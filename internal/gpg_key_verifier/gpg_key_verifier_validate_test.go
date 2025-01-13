// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package gpg_key_verifier_test

import (
	"encoding/base64"
	"testing"

	"github.com/opentofu/libregistry/internal/gpg_key_verifier"
)

const validPublicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mI0EZ4UfnAEEALNp2JW3iLceciejZu0u5ep3OJsYWg8UFep8cfzsn92xPuc2xiLZ
CioNpRzoJ4SZIxtbr7iGXrX3InMOdVPLTWrhYyJ5R2jE/b9n3g2vGgjBl0IWFGWQ
/NvABdQTjAPZv3/JNKy6UPYjkYlk/BNCP6mXpassBk26KxeOAXqJWNHjABEBAAG0
EXRlc3RAb3BlbnRvZnUub3JniNEEEwEIADsWIQSiDE0NNHdLnuzViwI0kurTLyqK
4wUCZ4UfnAIbAwULCQgHAgIiAgYVCgkICwIEFgIDAQIeBwIXgAAKCRA0kurTLyqK
45gGA/4nRQxhQ9hT1r3+wgPPVOa9Dx15eUHV55i9ASEBE3SNpteTAL+rJFgz4Sxn
6ydYwl4Mog0dWnbhzjuCQUAZVs3Gdt0jcNoycRT0CJhn9w2dZg8QA3ZVL5pdhWrU
cnz1VHGl3je06lF4ZHRwRrqLhM/J0SQv+cvkNst7eS7W7xiZVriNBGeFH5wBBACt
B+O88KYGQO0ZNyxc+ZkkD0dHqVBtKaPDLjixVaRVLEzDRPJCJlsF/XDtW+01zgrK
FzBYJHAvmuc92L3E+6ciqTvCX0lsK9/KjhSqcGxFAQN/mMgORoGeXspZc/uurDm3
X5rHIQsZV1X6XFRUPvz1R7MMu5i//jtS9L5FmNio+wARAQABiLYEGAEIACAWIQSi
DE0NNHdLnuzViwI0kurTLyqK4wUCZ4UfnAIbDAAKCRA0kurTLyqK4/v7A/9p9sdb
xBBYIsdAvCAojpHlZEUBiFWDKzD3/UEjAiFrlI6IjMAD/B6wl+dVtnPuXX6i62zZ
eg25uYGiUXUsE63xL1U9Pjmq+NXAt4L3/1xjjgGh/nmVJBnGUApwDzq+D6hOhYTG
79CondpeqED7JwlC6AJm+5zuHWqSxMMyCY48HQ==
=ExfT
-----END PGP PUBLIC KEY BLOCK-----`

func getSignatureDecoded() ([]byte, bool) {
	data, err := base64.StdEncoding.DecodeString(`iLMEAAEIAB0WIQSiDE0NNHdLnuzViwI0kurTLyqK4wUCZ4UfsgAKCRA0kurTLyqK4zHmA/4zGQe/JwphSmF6AreWo8RLoMLFqHcSM5UucIDxDo1Q07nx1uOKO4YS4ecNSANCktXqYWcSuZhLO2ujuV1VBAnl9U/VEhLsFYzVH1gQFYiJ0Jkep6oLrhifNsgTHBIJtCD2WWllatqSXT1Q9u3/CdBhlbBHAXiXoRo8bJanUphojA==`)

	if err != nil {
		return []byte{}, true
	}
	return data, false
}

func TestValidSignature(t *testing.T) {
	testKey := []byte(validPublicKey)

	gpgKeyVerifier, err := gpg_key_verifier.New(testKey)
	if err != nil {
		t.Fatalf("Failed to build gpgKeyVerifier (%v)", err)
	}

	data := []byte("test\n")
	signature, _ := getSignatureDecoded()
	err = gpgKeyVerifier.ValidateSignature(data, signature)
	if err != nil {
		t.Fatalf("Could not validate the signature (%v)", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	testKey := []byte(validPublicKey)

	gpgKeyVerifier, err := gpg_key_verifier.New(testKey)
	if err != nil {
		t.Fatalf("Failed to build gpgKeyVerifier (%v)", err)
	}

	data := []byte("test_invalid\n")
	signature, _ := getSignatureDecoded()

	err = gpgKeyVerifier.ValidateSignature(data, signature)
	if err == nil {
		t.Fatalf("Err should be non-nil (%v)", err)
	}
}
