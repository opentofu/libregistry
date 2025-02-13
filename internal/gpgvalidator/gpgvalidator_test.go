// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpgvalidator

import (
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
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
