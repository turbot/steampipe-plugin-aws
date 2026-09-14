package aws

import (
	"reflect"
	"testing"
)

func TestMatchRegionsToPatterns(t *testing.T) {
	validRegions := []string{
		"us-east-1", "us-east-2", "us-west-2",
		"eu-west-1", "eu-central-1",
		"ap-southeast-1",
		"me-south-1", "me-central-1",
		"us-gov-west-1",
	}

	testCases := []struct {
		name     string
		patterns []string
		expected []string
	}{
		// Positive patterns (existing behavior)
		{"wildcard all", []string{"*"}, []string{"us-east-1", "us-east-2", "us-west-2", "eu-west-1", "eu-central-1", "ap-southeast-1", "me-south-1", "me-central-1", "us-gov-west-1"}},
		{"exact region", []string{"us-east-1"}, []string{"us-east-1"}},
		{"prefix wildcard", []string{"us-*"}, []string{"us-east-1", "us-east-2", "us-west-2", "us-gov-west-1"}},
		{"question mark wildcard", []string{"us-east-?"}, []string{"us-east-1", "us-east-2"}},
		{"character class", []string{"me-[^s]*"}, []string{"me-central-1"}},
		{"multiple positive patterns", []string{"us-east-*", "eu-*"}, []string{"us-east-1", "us-east-2", "eu-west-1", "eu-central-1"}},
		{"overlapping patterns deduped", []string{"us-*", "us-east-1"}, []string{"us-east-1", "us-east-2", "us-west-2", "us-gov-west-1"}},
		{"no match", []string{"cn-*"}, nil},
		{"empty pattern list", []string{}, nil},
		{"malformed pattern matches nothing", []string{"us-[east-1"}, nil},

		// Exclusion patterns
		{"exclude exact", []string{"*", "!me-south-1"}, []string{"us-east-1", "us-east-2", "us-west-2", "eu-west-1", "eu-central-1", "ap-southeast-1", "me-central-1", "us-gov-west-1"}},
		{"exclude pattern", []string{"*", "!us-*"}, []string{"eu-west-1", "eu-central-1", "ap-southeast-1", "me-south-1", "me-central-1"}},
		{"exclude question mark", []string{"us-*", "!us-east-?"}, []string{"us-west-2", "us-gov-west-1"}},
		{"exclude gov from us", []string{"us-*", "!us-gov-*"}, []string{"us-east-1", "us-east-2", "us-west-2"}},
		{"exclude wins regardless of order", []string{"!me-south-1", "*"}, []string{"us-east-1", "us-east-2", "us-west-2", "eu-west-1", "eu-central-1", "ap-southeast-1", "me-central-1", "us-gov-west-1"}},
		{"exclude wins over explicit include", []string{"us-east-1", "!us-east-1"}, []string{}},
		{"exclude splits shared prefix", []string{"me-*", "!me-south-1"}, []string{"me-central-1"}},
		{"exclude with no match is a no-op", []string{"us-east-*", "!cn-north-1"}, []string{"us-east-1", "us-east-2"}},
		{"exclude outside include set is a no-op", []string{"eu-*", "!us-east-1"}, []string{"eu-west-1", "eu-central-1"}},
		{"duplicate exclusions", []string{"me-*", "!me-south-1", "!me-south-1"}, []string{"me-central-1"}},
		{"multiple exclusions", []string{"*", "!me-*", "!us-gov-*", "!ap-southeast-1"}, []string{"us-east-1", "us-east-2", "us-west-2", "eu-west-1", "eu-central-1"}},
		{"only exclusions matches nothing", []string{"!me-south-1"}, nil},
		{"only exclusion wildcard matches nothing", []string{"!*"}, nil},
		{"exclude everything", []string{"*", "!*"}, []string{}},
		{"malformed exclusion excludes nothing", []string{"us-east-*", "!us-[east-1"}, []string{"us-east-1", "us-east-2"}},
		{"bare bang matches nothing", []string{"us-east-1", "!"}, []string{"us-east-1"}},
		{"double bang is literal after prefix strip", []string{"*", "!!us-east-1"}, []string{"us-east-1", "us-east-2", "us-west-2", "eu-west-1", "eu-central-1", "ap-southeast-1", "me-south-1", "me-central-1", "us-gov-west-1"}},
		{"mid-string bang is literal, matches nothing", []string{"us-!east-1"}, nil},
		{"mid-string bang in exclusion is a no-op", []string{"us-east-*", "!us-!east-1"}, []string{"us-east-1", "us-east-2"}},
		{"trailing bang is literal, matches nothing", []string{"us-east-1!"}, nil},
		// Go's path.Match negates character classes with ^, not shell-style !.
		// Inside a class, ! is a literal member, so [!s] matches the letter s.
		{"shell-style class negation does not negate", []string{"me-[!s]*"}, []string{"me-south-1"}},
		{"go-style class negation", []string{"me-[^s]*"}, []string{"me-central-1"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchRegionsToPatterns(tc.patterns, validRegions)
			if len(got) == 0 && len(tc.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("matchRegionsToPatterns(%v) = %v, want %v", tc.patterns, got, tc.expected)
			}
		})
	}
}

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
