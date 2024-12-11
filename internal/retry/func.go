// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package retry

import (
	"context"
	"fmt"
	"github.com/opentofu/libregistry/logger"
	"time"
)

// Func retries a function until it exhausts the maximum tries.
func Func(
	ctx context.Context,
	description string,
	what func() error,
	isRetryable func(err error) bool,
	maxTries int,
	waitTime time.Duration,
	logger logger.Logger,
) error {
	_, err := Func2(ctx, description, func() (any, error) {
		return nil, what()
	}, isRetryable, maxTries, waitTime, logger)
	return err
}

// Func2 retries a function until it exhausts the maximum tries. This variant returns a value.
func Func2[T any](
	ctx context.Context,
	description string,
	what func() (T, error),
	isRetryable func(err error) bool,
	maxTries int,
	waitTime time.Duration,
	logger logger.Logger,
) (T, error) {
	tries := 0
	for {
		logger.Trace(ctx, "Attempting to %s (try %d of %d)...", description, tries+1, maxTries)
		val, err := what()
		if err == nil {
			return val, nil
		}
		if !isRetryable(err) {
			logger.Trace(ctx, "Non-retryable error encountered while attempting to %s, aborting (%v)", description, err)
			return val, fmt.Errorf("non-retryable error encountered while attempting to %s, aborting (%w)", description, err)
		}
		tries++
		if tries >= maxTries {
			logger.Trace(ctx, "Max tries exhausted while attempting to %s, aborting (last error was: %v)", description, err)
			return val, fmt.Errorf("max tries exhausted while attempting to %s, aborting (last error was: %w)", description, err)
		}
		logger.Trace(ctx, "Failed to %s, retrying in %s (%w)", description, waitTime, err)
		select {
		case <-ctx.Done():
			logger.Trace(ctx, "Timeout while attempting to %s (last error was: %w)", description, err)
			return val, fmt.Errorf("timeout while attempting to %s (last error was: %w)", description, err)
		case <-time.After(waitTime):
		}
	}
}
