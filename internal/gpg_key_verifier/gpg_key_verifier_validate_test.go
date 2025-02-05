// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpg_key_verifier_test

import (
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
)

func generateTestData(plainMessage []byte) (string, []byte, error) {
	// Generate a crypto key
	armoredKey, err := helper.GenerateKey("", "test@opentofu.org", nil, "rsa", 1024)
	if err != nil {
		return "", nil, err
	}

	key, err := crypto.NewKeyFromArmored(armoredKey)
	if err != nil {
		return "", nil, err
	}

	unlockedKeyObj, err := key.Unlock(nil)
	if err != nil {
		return "", nil, err
	}

	signingKeyRing, err := crypto.NewKeyRing(unlockedKeyObj)
	if err != nil {
		return "", nil, err
	}

	dataToSign := crypto.NewPlainMessage(plainMessage)

	signature, err := signingKeyRing.SignDetached(dataToSign)
	if err != nil {
		return "", nil, err
	}

	publicKey, err := unlockedKeyObj.GetArmoredPublicKey()
	if err != nil {
		return "", nil, err
	}

	return publicKey, signature.GetBinary(), nil
}

func TestValidSignature(t *testing.T) {
	data := []byte("test\n")
	testKey, signature, err := generateTestData(data)
	if err != nil {
		t.Fatalf("Failed to generate testData (%v)", err)
	}

	gpgKeyVerifier, err := gpg_key_verifier.New(testKey)
	if err != nil {
		t.Fatalf("Failed to build gpgKeyVerifier (%v)", err)
	}

	err = gpgKeyVerifier.ValidateSignature(data, signature)
	if err != nil {
		t.Fatalf("Could not validate the signature (%v)", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	data := []byte("test_invalid\n")
	testKey, _, err := generateTestData(data)
	if err != nil {
		t.Fatalf("Failed to generate testData (%v)", err)
	}

	signature := []byte("invalid_signature")

	gpgKeyVerifier, err := gpg_key_verifier.New(testKey)
	if err != nil {
		t.Fatalf("Failed to build gpgKeyVerifier (%v)", err)
	}

	err = gpgKeyVerifier.ValidateSignature(data, signature)
	if err == nil {
		t.Fatalf("Err should be non-nil (%v)", err)
	}
}
