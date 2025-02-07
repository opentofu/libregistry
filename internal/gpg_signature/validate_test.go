// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpg_signature

import (
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// generateTestKey returns a PGP key
func generateTestKey(t *testing.T) *crypto.Key {
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

	return unlockedKey
}

// generateTestData receives a plain message and returns a public key and a signature
func generateTestData(t *testing.T, plainMessage string) (string, string, error) {
	// Generate a crypto key
	key := generateTestKey(t)

	signingKeyRing, err := crypto.NewKeyRing(key)
	if err != nil {
		return "", "", err
	}

	dataToSign := crypto.NewPlainMessageFromString(plainMessage)

	signature, err := signingKeyRing.SignDetached(dataToSign)
	if err != nil {
		return "", "", err
	}

	publicKey, err := key.GetArmoredPublicKey()
	if err != nil {
		return "", "", err
	}

	armoredSignature, err := signature.GetArmored()
	if err != nil {
		return "", "", err
	}

	return publicKey, armoredSignature, nil
}

func TestValidSignature(t *testing.T) {
	data := "test\n"
	testKey, signature, err := generateTestData(t, data)
	if err != nil {
		t.Fatalf("Failed to generate testData (%v)", err)
	}

	gpgKeyVerifier, err := New(testKey)
	if err != nil {
		t.Fatalf("Failed to build gpgKeyVerifier (%v)", err)
	}

	err = gpgKeyVerifier.Validate(signature, data)
	if err != nil {
		t.Fatalf("Could not validate the signature (%v)", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	data := "test_invalid\n"
	testKey, _, err := generateTestData(t, data)
	if err != nil {
		t.Fatalf("Failed to generate testData (%v)", err)
	}

	signature := "invalid_signature"

	gpgKeyVerifier, err := New(testKey)
	if err != nil {
		t.Fatalf("Failed to build gpgKeyVerifier (%v)", err)
	}

	err = gpgKeyVerifier.Validate(signature, data)
	if err == nil {
		t.Fatalf("Err should be non-nil (%v)", err)
	}
}
