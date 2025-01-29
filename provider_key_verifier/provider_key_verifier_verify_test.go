package provider_key_verifier

import (
	"context"
	"testing"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
)

type mockMetadata struct {
	metadata.API
}

func (m mockMetadata) GetProvider(ctx context.Context, addr provider.Addr, test bool) (provider.Metadata, error) {
	providerMetadata := provider.Metadata{
		Versions: provider.VersionList{
			provider.Version{
				Version: "0.2.0",
			},
		},
	}
	return providerMetadata, nil
}

func TestProviderVerify(t *testing.T) {
	metadataAPI := &mockMetadata{}
	expectedData := []byte("test")
	httpClient := *generateTestClient(expectedData)

	ctx := context.Background()

	keyData, err := generateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	pkv, err := New(keyData, metadataAPI, WithHTTPClient(httpClient), WithCheckFn(func(pkv providerKeyVerifier, ctx context.Context, version provider.Version) error {
		return nil
	}))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	addr := provider.Addr{
		Name:      "test",
		Namespace: "opentofu",
	}

	data, err := pkv.VerifyProvider(ctx, addr)
	if err != nil {
		t.Fatalf("Failed to verify key: %v", err)
	}

	if data[0].Version != "0.2.0" {
		t.Fatalf("Wrong version was returned %s", data[0])
	}
}
