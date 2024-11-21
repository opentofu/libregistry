// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package ociclient

import "testing"

func TestParseWWWAuthenticate(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput []OCIRawAuthScheme
		error          bool
	}

	testCases := []testCase{
		{
			name:           "empty",
			input:          "",
			expectedOutput: nil,
			error:          false,
		},
		{
			name:  "only-type",
			input: `Basic`,
			expectedOutput: []OCIRawAuthScheme{
				{
					Type: "Basic",
				},
			},
			error: false,
		},
		{
			name:  "start-equals",
			input: `=value`,
			error: true,
		},
		{
			name:  "type-with-comma",
			input: `"Basic"`,
			error: true,
		},
		{
			name:  "start-key",
			input: `key=value`,
			error: true,
		},
		{
			name:  "single-parameter-no-quotes",
			input: `Basic foo=bar`,
			expectedOutput: []OCIRawAuthScheme{
				{
					Type:   "Basic",
					Params: map[string]string{"foo": "bar"},
				},
			},
			error: false,
		},
		{
			name:  "key-with-comma",
			input: `Basic "foo"="bar""`,
			error: true,
		},
		{
			name:  "single-parameter-quotes",
			input: `Basic foo="bar"`,
			expectedOutput: []OCIRawAuthScheme{
				{
					Type:   "Basic",
					Params: map[string]string{"foo": "bar"},
				},
			},
			error: false,
		},
		{
			name:  "single-parameter-quotes-escape-1",
			input: `Basic foo="b\"a\"r"`,
			expectedOutput: []OCIRawAuthScheme{
				{
					Type:   "Basic",
					Params: map[string]string{"foo": "b\"a\"r"},
				},
			},
			error: false,
		},
		{
			name:  "single-parameter-quotes-inside",
			input: `Basic foo=b"ar"`,
			expectedOutput: []OCIRawAuthScheme{
				{
					Type:   "Basic",
					Params: map[string]string{"foo": "b\"ar\""},
				},
			},
			error: false,
		},
		{
			name:  "single-parameter-quotes-inside-escaped",
			input: `Basic foo=b\"ar"`,
			expectedOutput: []OCIRawAuthScheme{
				{
					Type:   "Basic",
					Params: map[string]string{"foo": "b\\\"ar\""},
				},
			},
			error: false,
		},
		{
			name:  "double-scheme",
			input: `Basic foo=bar, Digest baz=foo`,
			expectedOutput: []OCIRawAuthScheme{
				{
					Type:   "Basic",
					Params: map[string]string{"foo": "bar"},
				},
				{
					Type:   "Digest",
					Params: map[string]string{"baz": "foo"},
				},
			},
			error: false,
		},
		{
			name:  "double-scheme-no-comma",
			input: `Basic Digest baz=foo`,
			error: true,
		},
		{
			name:  "extra-comma",
			input: `Basic foo=bar, , Digest baz=foo`,
			expectedOutput: []OCIRawAuthScheme{
				{
					Type:   "Basic",
					Params: map[string]string{"foo": "bar"},
				},
				{
					Type:   "Digest",
					Params: map[string]string{"baz": "foo"},
				},
			},
			error: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseWWWAuthenticate(tc.input)
			if tc.error {
				if err == nil {
					t.Fatalf("Expected error not encountered")
				}
				return
			} else if !tc.error && err != nil {
				t.Fatalf("Unexpected error encountered: %v", err)
			}

			if len(result) != len(tc.expectedOutput) {
				t.Fatalf("Incorrect number of results: %d (expected: %d)", len(result), len(tc.expectedOutput))
			}

			for i, item := range result {
				if item.Type != tc.expectedOutput[i].Type {
					t.Fatalf("Output type mismatch for item %d: %s (expected: %s)", i, item.Type, tc.expectedOutput[i].Type)
				}
				if len(item.Params) != len(tc.expectedOutput[i].Params) {
					t.Fatalf("Incorrect number of parameters for item %d: %d (expected: %d)", i, len(item.Params), len(tc.expectedOutput[i].Params))
				}
				for expectedKey, expectedValue := range tc.expectedOutput[i].Params {
					value, ok := item.Params[expectedKey]
					if !ok {
						t.Fatalf("Key %s not found on item %d (expected value: %s)", expectedKey, i, expectedValue)
					}
					if expectedValue != value {
						t.Fatalf("Incorrect value %s for key %s on item %d (expected value: %s)", value, expectedKey, i, expectedValue)
					}
				}
			}
		})
	}
}
