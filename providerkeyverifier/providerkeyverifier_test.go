// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package providerkeyverifier

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/types/provider"
	"github.com/stretchr/testify/require"
)

func generateKey(t testing.TB) *crypto.Key {
	t.Helper()
	armoredKey, err := helper.GenerateKey("opentofu", "test@opentofu.org", nil, "rsa", 1024)
	if err != nil {
		t.Fatalf("Error when generating the armored string: %v", err)
	}

	key, err := crypto.NewKeyFromArmored(armoredKey)
	if err != nil {
		t.Fatalf("Error when creating a new key from armored string: %v", err)
	}

	unlockedKey, err := key.Unlock(nil)
	if err != nil {
		t.Fatalf("Error when unlocking the key: %v", err)
	}

	return unlockedKey
}

// getPubKey returns a PGP public key
func getPubKey(t testing.TB, key *crypto.Key) string {
	t.Helper()
	pubKey, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("Failed to get the armored public key: %v", err)
	}

	return pubKey
}

// generate Signature and data
func generateSignedData(t testing.TB, key *crypto.Key, msg []byte) ([]byte, []byte) {
	t.Helper()
	var plainMsg = crypto.NewPlainMessage(msg)

	signingKeyRing, err := crypto.NewKeyRing(key)
	if err != nil {
		t.Fatalf("Failed to create a new key ring: %v", err)
	}

	pgpSignature, err := signingKeyRing.SignDetached(plainMsg)
	if err != nil {
		t.Fatalf("Failed to sign detached: %v", err)
	}

	return pgpSignature.GetBinary(), plainMsg.GetBinary()
}

// newTestServer is used to mock the HTTP requests and return the data
// `/SHASumsURL/` and `/SHASumsSignatureURL/` are used to mimic the opentofu's registry API
func newTestServer(t testing.TB, key *crypto.Key, expected []byte) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	sig, data := generateSignedData(t, key, []byte("message"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(expected)
		if err != nil {
			t.Fatalf("Couldn't write to testing response of /: %v", err)
		}
	})

	mux.HandleFunc("/SHASumsURL/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(data)
		if err != nil {
			t.Fatalf("Couldn't write to testing response of /SHASumsURL/: %v", err)
		}
	})
	mux.HandleFunc("/SHASumsSignatureURL/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(sig)
		if err != nil {
			t.Fatalf("Couldn't write to testing response of /SHASumsSignatureURL/: %v", err)
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

func setupProviderCall(t testing.TB, shaSumsURL string, shaSumsSignatureURL string) ProviderKeyVerifier {
	t.Helper()

	key := generateKey(t)
	pubKey := getPubKey(t, key)
	srv := newTestServer(t, key, []byte("test"))
	metadataAPI := &mockMetadata{
		shaSumsURL:          srv.URL + shaSumsURL,
		shaSumsSignatureURL: srv.URL + shaSumsSignatureURL,
	}

	httpClient := srv.Client()

	pkv, err := New(pubKey, metadataAPI, WithHTTPClient(httpClient))
	require.NoError(t, err)

	return pkv
}
