// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
)

func (r registryDataAPI) ListProviderNamespaceAliases(ctx context.Context) (map[string]string, error) {
	return map[string]string{
		"opentofu": "hashicorp",
	}, nil
}
