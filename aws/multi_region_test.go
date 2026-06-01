package aws

import (
	"testing"
)

func TestAwsLastResortRegionFromRegionWildcard(t *testing.T) {
	testCases := []struct {
		name           string
		regionWildcard string
		expected       string
	}{
		// Commercial regions resolve to us-east-1
		{"commercial us", "us-west-2", "us-east-1"},
		{"commercial us wildcard", "us-*", "us-east-1"},
		{"commercial eu", "eu-west-1", "us-east-1"},
		{"commercial ap", "ap-southeast-6", "us-east-1"},
		{"commercial af", "af-south-1", "us-east-1"},
		{"commercial ca", "ca-central-1", "us-east-1"},
		{"commercial me", "me-central-1", "us-east-1"},
		{"commercial sa", "sa-east-1", "us-east-1"},
		// Israel (Tel Aviv) and Mexico (Central) — added to the commercial prefix list
		{"commercial il", "il-central-1", "us-east-1"},
		{"commercial il wildcard", "il-*", "us-east-1"},
		{"commercial mx", "mx-central-1", "us-east-1"},
		{"commercial mx wildcard", "mx-*", "us-east-1"},
		// Obscure partitions take precedence over commercial prefixes
		{"gov", "us-gov-west-1", "us-gov-west-1"},
		{"gov wildcard", "us-gov-*", "us-gov-west-1"},
		{"china", "cn-north-1", "cn-northwest-1"},
		{"china wildcard", "cn-*", "cn-northwest-1"},
		{"isob", "us-isob-east-1", "us-isob-east-1"},
		{"iso", "us-iso-east-1", "us-iso-east-1"},
		// Unknown partition resolves to empty string
		{"unknown", "crap", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := awsLastResortRegionFromRegionWildcard(tc.regionWildcard)
			if got != tc.expected {
				t.Errorf("awsLastResortRegionFromRegionWildcard(%q) = %q, want %q", tc.regionWildcard, got, tc.expected)
			}
		})
	}
}
