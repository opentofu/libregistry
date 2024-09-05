package provider_test

import (
	"testing"

	"github.com/opentofu/libregistry/types"
	"github.com/opentofu/libregistry/types/provider"
)

func TestVersionNumber(t *testing.T) {
	type testCase struct {
		version                 provider.VersionNumber
		expectError             bool
		expectedMajor           int
		expectedMinor           int
		expectedPatch           int
		expectedStability       string
		expectedStabilityNumber int
	}

	for name, tc := range map[string]testCase{
		"major": {
			"v1",
			true,
			0, 0, 0, "", 0,
		},
		"minor": {
			"v1.2",
			true,
			0, 0, 0, "", 0,
		},
		"simple": {
			"v1.2.3",
			false,
			1, 2, 3, "", 0,
		},
		"unprefixed": {
			"1.2.3",
			false,
			1, 2, 3, "", 0,
		},
		"stability": {
			"1.2.3-alpha",
			true,
			0, 0, 0, "", 0,
		},
		"full": {
			"1.2.3-alpha1",
			false,
			1, 2, 3, "alpha", 1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			versionNumber := types.VersionNumber(tc.version)
			major, minor, patch, stability, stabilityNumber, err := versionNumber.Parse()
			if tc.expectError && err == nil {
				t.Fatalf("Expected error was not returned.")
			} else if !tc.expectError && err != nil {
				t.Fatalf("Unexpected error returned: %v", err)
			}
			if major != tc.expectedMajor {
				t.Fatalf("Incorrect major version: %d (expected: %d)", major, tc.expectedMajor)
			}
			if minor != tc.expectedMinor {
				t.Fatalf("Incorrect minor version: %d (expected: %d)", minor, tc.expectedMinor)
			}
			if patch != tc.expectedPatch {
				t.Fatalf("Incorrect patch version: %d (expected: %d)", patch, tc.expectedPatch)
			}
			if stability != tc.expectedStability {
				t.Fatalf("Incorrect stability: %s (expected: %s)", stability, tc.expectedStability)
			}
			if stabilityNumber != tc.expectedStabilityNumber {
				t.Fatalf("Incorrect patch version: %d (expected: %d)", stabilityNumber, tc.expectedStabilityNumber)
			}
		})
	}
}
