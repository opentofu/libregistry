// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
)

func (r registryDataAPI) ListProviderNamespaceAliases(ctx context.Context) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (r registryDataAPI) PutProviderNamespaceAlias(ctx context.Context, from string, to string) error {
	//TODO implement me
	panic("implement me")
}

func (r registryDataAPI) DeleteProviderNamespaceAlias(ctx context.Context, from string) error {
	//TODO implement me
	panic("implement me")
}
