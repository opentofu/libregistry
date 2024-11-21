// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/opentofu/libregistry/logger"
	"golang.org/x/sync/errgroup"
)

type ociClient struct {
	tempDirectory string
	rawClient     RawOCIClient
	logger        logger.Logger
}

func (o ociClient) ListReferences(ctx context.Context, addr OCIAddr) ([]OCIReference, OCIWarnings, error) {
	if err := addr.Validate(); err != nil {
		return nil, nil, err
	}
	o.logger.Debug(ctx, "Listing references for OCI image %s...", addr)
	warnings, err := o.rawClient.Check(ctx, addr.Registry)
	if err != nil {
		return nil, warnings, err
	}
	response, newWarnings, err := o.rawClient.ContentDiscovery(ctx, addr)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}
	return response.Tags, warnings, err
}

func (o ociClient) ResolvePlatformImageDigest(ctx context.Context, addrRef OCIAddrWithReference, opts ...ClientPullOpt) (OCIDigest, OCIWarnings, error) {
	if err := addrRef.Validate(); err != nil {
		return "", nil, err
	}

	pullConfig := ClientPullConfig{}
	for _, opt := range opts {
		if err := opt(&pullConfig); err != nil {
			return "", nil, err
		}
	}
	if err := pullConfig.ApplyDefaultsAndValidate(); err != nil {
		return "", nil, err
	}

	o.logger.Debug(ctx, "Resolving image digest for %s...", addrRef)
	warnings, err := o.rawClient.Check(ctx, addrRef.Registry)
	if err != nil {
		return "", warnings, err
	}

	o.logger.Trace(ctx, "Getting main manifest for %s...", addrRef)
	mainManifest, newWarnings, err := o.rawClient.GetManifest(ctx, addrRef)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		o.logger.Trace(ctx, "Getting main manifest for %s failed. (%v)", addrRef, err)
		return "", warnings, err
	}
	indexManifest, ok := mainManifest.AsIndexManifest()
	if !ok {
		o.logger.Trace(ctx, "Main manifest for %s is not an image or index manifest, %s found instead.", addrRef, mainManifest.GetMediaType())
		return "", warnings, fmt.Errorf("main manifest for %s is not an image or index manifest (%s found instead)", addrRef, mainManifest.GetMediaType())
	}
	o.logger.Trace(ctx, "Found multi-arch index manifest for %s, searching for a platform image for %s / %s", addrRef, pullConfig.GOOS, pullConfig.GOARCH, addrRef)
	var platform *OCIRawDescriptor
	for _, platformManifest := range indexManifest.Manifests {
		platformManifest := platformManifest
		if platformManifest.Platform.OS == pullConfig.GOOS &&
			platformManifest.Platform.Architecture == pullConfig.GOARCH &&
			(platformManifest.MediaType == MediaTypeDockerImage || platformManifest.MediaType == MediaTypeOCIImage) {
			platform = &platformManifest.OCIRawDescriptor
			break
		}
	}
	if platform == nil {
		// TODO typed error
		o.logger.Trace(ctx, "No suitable image found for GOOS %s and GOARCH %s in %s", pullConfig.GOOS, pullConfig.GOARCH, addrRef)
		return "", warnings, fmt.Errorf("no suitable image found for GOOS %s and GOARCH %s in %s", pullConfig.GOOS, pullConfig.GOARCH, addrRef)
	}
	return platform.Digest, warnings, nil
}

func (o ociClient) PullImageWithImageDigest(ctx context.Context, addrRef OCIAddrWithDigest) (PulledOCIImage, OCIWarnings, error) {
	if err := addrRef.Validate(); err != nil {
		return nil, nil, err
	}

	o.logger.Debug(ctx, "Pulling image %s...", addrRef)
	warnings, err := o.rawClient.Check(ctx, addrRef.Registry)
	if err != nil {
		return nil, warnings, err
	}

	o.logger.Trace(ctx, "Getting main manifest for %s...", addrRef)
	mainManifest, newWarnings, err := o.rawClient.GetManifest(ctx, OCIAddrWithReference{
		OCIAddr:   addrRef.OCIAddr,
		Reference: OCIReference(addrRef.Digest),
	})
	warnings = append(warnings, newWarnings...)
	if err != nil {
		o.logger.Trace(ctx, "Getting main manifest for %s failed. (%v)", addrRef, err)
		return nil, warnings, err
	}
	imageManifest, ok := mainManifest.AsImageManifest()
	if !ok {
		return nil, warnings, fmt.Errorf("the specified digest does not point to an image manifest")
	}
	layers := imageManifest.Layers
	if err := addrRef.Validate(); err != nil {
		return nil, nil, err
	}
	o.logger.Trace(ctx, "Downloading %d layers for %s...", len(layers), addrRef)
	errGroup, errGroupCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(6)
	lock := &sync.Mutex{}
	result := &pulledOCIImage{
		layers:       nil,
		currentLayer: -1,
	}
	for _, layer := range layers {
		layer := layer
		// Note: we need to replace colons with underscores here because Windows doesn't support them in filenames.
		fn := path.Join(o.tempDirectory, strings.ReplaceAll(string(layer.Digest), ":", "_")+".tar.gz")
		result.layers = append(result.layers, pulledLayer{
			digest:   layer.Digest,
			tempFile: fn,
		})
		errGroup.Go(func() error {
			o.logger.Trace(ctx, "Downloading %s layer %s...", addrRef, layer)
			blob, newWarnings, err := o.rawClient.GetBlob(errGroupCtx, OCIAddrWithDigest{
				addrRef.OCIAddr,
				layer.Digest,
			})
			lock.Lock()
			warnings = append(warnings, newWarnings...)
			lock.Unlock()
			if err != nil {
				// TODO typed error
				o.logger.Trace(ctx, "Failed to pull %s layer %s. (%v)", addrRef, layer, err)
				return fmt.Errorf("failed to pull layer %s (%w)", layer.Digest, err)
			}
			fh, err := os.Create(fn)
			if err != nil {
				// TODO typed error
				o.logger.Trace(ctx, "Failed to create temporary file %s for %s layer %s. (%v)", fn, addrRef, layer, err)
				return fmt.Errorf("failed to create temporary file %s for layer %s (%w)", fn, layer.Digest, err)
			}
			if _, err := io.Copy(fh, blob); err != nil {
				// TODO typed error
				o.logger.Trace(ctx, "Failed to write temporary file %s for %s layer %s. (%v)", fn, addrRef, layer, err)
				return fmt.Errorf("failed to write temporary file %s for layer %s (%w)", fn, layer.Digest, err)
			}
			if err := fh.Close(); err != nil {
				// TODO typed error
				o.logger.Trace(ctx, "Failed to close temporary file %s for %s layer %s. (%v)", fn, addrRef, layer, err)
				return fmt.Errorf("failed to close temporary file %s for layer %s (%w)", fn, layer.Digest, err)
			}
			o.logger.Trace(ctx, "Downloading %s layer %s complete.", addrRef, layer)
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		o.logger.Trace(ctx, "Downloading layers for %s failed, cleaning up temporary files... (%v)", addrRef, err)
		lock.Lock()
		defer lock.Unlock()
		if err := result.Close(); err != nil {
			o.logger.Warn(ctx, "Could not clean up partial pull for %s (%w)", addrRef, err)
		}
		return nil, warnings, err
	}
	o.logger.Trace(ctx, "Downloading layers for %s complete.", addrRef)
	lock.Lock()
	defer lock.Unlock()
	return result, warnings, nil
}

func (o ociClient) PullImage(ctx context.Context, addrRef OCIAddrWithReference, opts ...ClientPullOpt) (PulledOCIImage, OCIWarnings, error) {
	digest, warnings, err := o.ResolvePlatformImageDigest(ctx, addrRef, opts...)
	if err != nil {
		return nil, warnings, err
	}
	image, newWarnings, err := o.PullImageWithImageDigest(ctx, OCIAddrWithDigest{
		OCIAddr: addrRef.OCIAddr,
		Digest:  digest,
	})
	return image, append(warnings, newWarnings...), err
}

// pulledOCIImage implements the PulledOCIImage interface to give access to downloaded layers.
// TODO when files are overwritten or removed, this is currently not taken into account.
type pulledOCIImage struct {
	layers       []pulledLayer
	currentLayer int
	fh           *os.File
	gzip         *gzip.Reader
	tar          *tar.Reader
	header       *tar.Header
}

func (p *pulledOCIImage) Next() (bool, error) {
	for {
		// If we have no tar file open, go to the next layer. If there are no more layers, we are done.
		if p.tar == nil {
			done, err := p.nextLayer()
			if err != nil {
				return false, err
			}
			if done {
				return false, nil
			}
		}
		// Here we try to set the tar pointer to the next file.
		if p.tar != nil {
			if err := p.nextFileInTar(); err != nil {
				return false, err
			}
		}
		// Here we check if nextFileInTar opened a new file or if the tar file has reached EOF.
		// If we reach EOF, loop around and try the next tar file.
		if p.header != nil {
			return true, nil
		}
	}
}

func (p *pulledOCIImage) nextFileInTar() error {
	var err error
	p.header = nil
	p.header, err = p.tar.Next()
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Rotate to the next file.
			p.tar = nil
			// No error handling needed as we are only reading
			if err := p.gzip.Close(); err != nil {
				return fmt.Errorf("unexpected error while closing GZIP stream %s (%w)", p.layers[p.currentLayer].tempFile, err)
			}
			p.gzip = nil
			if err := p.fh.Close(); err != nil {
				return fmt.Errorf("unexpected error while closing file %s (%w)", p.layers[p.currentLayer].tempFile, err)
			}
			p.fh = nil
		} else {
			return fmt.Errorf("unexpected error while reading tar file %s (%w)", p.layers[p.currentLayer].tempFile, err)
		}
	}
	return nil
}

func (p *pulledOCIImage) nextLayer() (bool, error) {
	var err error
	if p.currentLayer+1 >= len(p.layers) {
		// No more archives
		return true, nil
	}
	p.currentLayer++
	p.fh, err = os.Open(p.layers[p.currentLayer].tempFile)
	if err != nil {
		return false, fmt.Errorf("failed to open %s (%w)", p.layers[p.currentLayer].tempFile, err)
	}
	p.gzip, err = gzip.NewReader(p.fh)
	if err != nil {
		return false, fmt.Errorf("failed to open GZIP stream for %s (%w)", p.layers[p.currentLayer].tempFile, err)
	}
	p.tar = tar.NewReader(p.gzip)
	return false, nil
}

func (p *pulledOCIImage) FileInfo() fs.FileInfo {
	if p.header != nil {
		return p.header.FileInfo()
	}
	return nil
}

func (p *pulledOCIImage) Filename() string {
	if p.header != nil {
		return p.header.Name
	}
	return ""
}

func (p *pulledOCIImage) Read(data []byte) (n int, err error) {
	if p.tar != nil {
		return p.tar.Read(data)
	}
	return 0, fmt.Errorf("no valid layer open, please call Next() before calling read")
}

func (p *pulledOCIImage) Close() error {
	p.tar = nil
	if p.gzip != nil {
		if err := p.gzip.Close(); err != nil {
			return fmt.Errorf("unexpected error while closing %s (%w)", p.layers[p.currentLayer].tempFile, err)
		}
	}
	if p.fh != nil {
		if err := p.fh.Close(); err != nil {
			return fmt.Errorf("unexpected error while closing %s (%w)", p.layers[p.currentLayer].tempFile, err)
		}
	}
	if runtime.GOOS == "windows" {
		// Make sure all file handles are closed on Windows so no files are blocked from removal.
		runtime.GC()
	}
	errGroup := errgroup.Group{}
	for _, layer := range p.layers {
		layer := layer
		errGroup.Go(func() error {
			if err := os.RemoveAll(layer.tempFile); err != nil {
				return fmt.Errorf("failed to remove temporary file %s (%w)", layer.tempFile, err)
			}
			return nil
		})
	}
	p.currentLayer = math.MaxInt
	return errGroup.Wait()
}

type pulledLayer struct {
	digest   OCIDigest
	tempFile string
}

var _ OCIClient = &ociClient{}
