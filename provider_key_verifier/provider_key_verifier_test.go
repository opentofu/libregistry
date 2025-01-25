package provider_key_verifier

import (
	"net/http"
	"testing"
)

type ClientMock struct {
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestProviderConfig(t *testing.T) {
	mockClient := &ClientMock{}

	pkv, err := New(mockClient, nil, WithVersionsToCheck(5))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	if pkv.(*providerKeyVerifier).versionsToCheck != 5 {
		t.Fatalf("Incorrect number of versions to check: %v, expecting %v.", pkv.(*providerKeyVerifier).versionsToCheck, 10)
	}
}
