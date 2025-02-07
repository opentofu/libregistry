// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

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

func generateTestClient(t *testing.T, expected string) *http.Client {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", expected)
		}),
	)

	t.Cleanup(func() {
		srv.Close()
	})
	return srv.Client()
}

func TestProviderConfig(t *testing.T) {
	httpClient := generateTestClient(t, "test")
	pubKey := generateTestPubKey(t)

	pkv, err := New(pubKey, nil, WithNumVersionsToCheck(5), WithHTTPClient(httpClient))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	if pkv.(*providerKey).config.NumVersionsToCheck != 5 {
		t.Fatalf("Incorrect number of versions to check: %v, expecting %v.", pkv.(*providerKey).config.NumVersionsToCheck, 10)
	}
}
