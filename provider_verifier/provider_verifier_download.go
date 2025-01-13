package provider_verifier

import (
	"context"
	"fmt"
	"io"
)

func (kv keyVerification) DownloadFile(ctx context.Context, url string) ([]byte, error) {
	response, err := kv.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
