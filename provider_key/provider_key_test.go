// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// generateTestPubKey returns a PGP public key
func generateTestPubKey(t *testing.T) string {
	armoredKey, err := helper.GenerateKey("", "test@opentofu.org", nil, "rsa", 1024)
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

	pubKey, err := unlockedKey.GetArmoredPublicKey()
	if err != nil {
		t.Error(err)
	}

	return pubKey
}

func generateTestServer(t *testing.T, expected string) *httptest.Server {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expected))
		}),
	)

	t.Cleanup(func() {
		srv.Close()
	})
	return srv
}

func TestProviderConfig(t *testing.T) {
	srv := generateTestServer(t, "test")
	httpClient := srv.Client()
	pubKey := generateTestPubKey(t)

	pkv, err := New(pubKey, nil, WithNumVersionsToCheck(5), WithHTTPClient(httpClient))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	if pkv.(*providerKey).config.NumVersionsToCheck != 5 {
		t.Fatalf("Incorrect number of versions to check: %v, expecting %v.", pkv.(*providerKey).config.NumVersionsToCheck, 10)
	}
}
