// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package provider_key_verifier

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

func generateTestClient(expected []byte) *http.Client {
	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", expected)
		}),
	)

	defer svr.Close()

	return svr.Client()
}

func generateKey() (string, error) {
	armoredKey, err := helper.GenerateKey("", "test@opentofu.org", nil, "rsa", 1024)
	if err != nil {
		return "", err
	}

	key, err := crypto.NewKeyFromArmored(armoredKey)
	if err != nil {
		return "", err
	}

	unlockedKeyObj, err := key.Unlock(nil)
	if err != nil {
		return "", err
	}

	pubKey, err := unlockedKeyObj.GetArmoredPublicKey()
	if err != nil {
		return "", err
	}

	return pubKey, nil
}

func TestProviderConfig(t *testing.T) {
	httpClient := *generateTestClient([]byte("test"))
	key, err := generateKey()
	if err != nil {
		t.Fatalf("couldn't create key: %v", err)
	}

	pkv, err := New(key, nil, WithVersionsToCheck(5), WithHTTPClient(httpClient))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	if pkv.(*providerKeyVerifier).versionsToCheck != 5 {
		t.Fatalf("Incorrect number of versions to check: %v, expecting %v.", pkv.(*providerKeyVerifier).versionsToCheck, 10)
	}
}

func TestProviderNoConfig(t *testing.T) {
	httpClient := *generateTestClient([]byte("test"))
	key, err := generateKey()
	if err != nil {
		t.Fatalf("couldn't create key: %v", err)
	}

	_, err = New(key, nil, WithHTTPClient(httpClient))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}
}
