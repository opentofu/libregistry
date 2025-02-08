// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

func generateKey(t *testing.T) *crypto.Key {
	armoredKey, err := helper.GenerateKey("opentofu", "test@opentofu.org", nil, "rsa", 1024)
	if err != nil {
		t.Error(err)
	}

	key, err := crypto.NewKeyFromArmored(armoredKey)
	if err != nil {
		t.Error(err)
	}

	unlockedKey, err := key.Unlock(nil)
	if err != nil {
		t.Error(err)
	}

	return unlockedKey
}

// generateTestPubKey returns a PGP public key
func getPubKey(t *testing.T, key *crypto.Key) string {
	pubKey, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Error(err)
	}

	return pubKey
}

// generate Signature and data
func generateSignedData(t *testing.T, key *crypto.Key, msg []byte) ([]byte, []byte) {
	var plainMsg = crypto.NewPlainMessage(msg)

	signingKeyRing, err := crypto.NewKeyRing(key)
	if err != nil {
		t.Error("failed to create a new key ring", err)
	}

	pgpSignature, err := signingKeyRing.SignDetached(plainMsg)
	if err != nil {
		t.Error("failed to sign detached", err)
	}

	return pgpSignature.GetBinary(), plainMsg.GetBinary()
}

// generateTestServer used to mock the HTTP requests and return what's expected
func generateTestServer(t *testing.T, key *crypto.Key, expected []byte) *httptest.Server {
	mux := http.NewServeMux()

	sig, data := generateSignedData(t, key, []byte("message"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(expected)
		if err != nil {
			t.Errorf("Couldn't write to testing response of /: %v", err)
		}
	})

	mux.HandleFunc("/SHASumsURL/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(data)
		if err != nil {
			t.Errorf("Couldn't write to testing response of /SHASumsURL/: %v", err)
		}
	})
	mux.HandleFunc("/SHASumsSignatureURL/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(sig)
		if err != nil {
			t.Errorf("Couldn't write to testing response of /SHASumsSignatureURL/: %v", err)
		}
	})

	srv := httptest.NewServer(mux)

	t.Cleanup(func() {
		srv.Close()
	})
	return srv
}

type mockMetadata struct {
	metadata.API
	shaSumsURL          string
	shaSumsSignatureURL string
}

func (m mockMetadata) GetProvider(ctx context.Context, addr provider.Addr, resolveAliases bool) (provider.Metadata, error) {
	return provider.Metadata{
		Versions: provider.VersionList{
			provider.Version{
				Version:             "0.2.0",
				SHASumsURL:          m.shaSumsURL,
				SHASumsSignatureURL: m.shaSumsSignatureURL,
			},
		},
	}, nil
}

func setupProviderCall(t *testing.T, shaSumsURL string, shaSumsSignatureURL string) ProviderKey {
	key := generateKey(t)
	pubKey := getPubKey(t, key)
	srv := generateTestServer(t, key, []byte("test"))
	metadataAPI := &mockMetadata{
		shaSumsURL:          srv.URL + shaSumsURL,
		shaSumsSignatureURL: srv.URL + shaSumsSignatureURL,
	}

	httpClient := srv.Client()

	pkv, err := New(pubKey, metadataAPI, WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	return pkv
}

func TestProviderConfig(t *testing.T) {
	key := generateKey(t)
	pubKey := getPubKey(t, key)

	pkv, err := New(pubKey, nil, WithNumVersionsToCheck(5))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	if pkv.(*providerKey).config.NumVersionsToCheck != 5 {
		t.Fatalf("Incorrect number of versions to check: %v, expecting %v.", pkv.(*providerKey).config.NumVersionsToCheck, 10)
	}
}
