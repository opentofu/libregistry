package provider_key_verifier

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (pkv providerKeyVerifier) downloadFile(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request (%w)", err)
	}

	response, err := pkv.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from %s: %w", url, err)
	}

	return contents, nil
}
