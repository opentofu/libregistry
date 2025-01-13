package provider_verifier

import (
	"context"
	"fmt"
	"io"
)

func (kv keyVerification) downloadFile(ctx context.Context, url string) ([]byte, error) {
	response, err := kv.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
