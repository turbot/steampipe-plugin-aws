package aws

import (
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/quals"
)

func TestBuildTagFilter(t *testing.T) {
	testCases := []struct {
		name           string
		qualJSON       string
		expectedCount  int
		expectedKeys   []string
		expectedValues map[string][]string
		expectError    bool
	}{
		{
			name:          "empty quals",
			qualJSON:      "",
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:           "single key without values",
			qualJSON:       `[{"key":"Environment"}]`,
			expectedCount:  1,
			expectedKeys:   []string{"Environment"},
			expectedValues: map[string][]string{"Environment": nil},
			expectError:    false,
		},
		{
			name:           "single key with single value",
			qualJSON:       `[{"key":"Environment","values":["prod"]}]`,
			expectedCount:  1,
			expectedKeys:   []string{"Environment"},
			expectedValues: map[string][]string{"Environment": {"prod"}},
			expectError:    false,
		},
		{
			name:           "single key with multiple values",
			qualJSON:       `[{"key":"Environment","values":["prod","dev","staging"]}]`,
			expectedCount:  1,
			expectedKeys:   []string{"Environment"},
			expectedValues: map[string][]string{"Environment": {"prod", "dev", "staging"}},
			expectError:    false,
		},
		{
			name:          "multiple keys (AND logic)",
			qualJSON:      `[{"key":"Environment","values":["prod"]},{"key":"Team","values":["backend"]}]`,
			expectedCount: 2,
			expectedKeys:  []string{"Environment", "Team"},
			expectedValues: map[string][]string{
				"Environment": {"prod"},
				"Team":        {"backend"},
			},
			expectError: false,
		},
		{
			name:          "mixed keys with and without values",
			qualJSON:      `[{"key":"CostCenter"},{"key":"Environment","values":["prod","dev"]}]`,
			expectedCount: 2,
			expectedKeys:  []string{"CostCenter", "Environment"},
			expectedValues: map[string][]string{
				"CostCenter":  nil,
				"Environment": {"prod", "dev"},
			},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			qualJSON:    `[{"key":"Environment"`,
			expectError: true,
		},
		{
			name:        "invalid values type - string instead of array",
			qualJSON:    `[{"key":"Environment","values":"prod"}]`,
			expectError: true,
		},
		{
			name:        "invalid values type - number instead of array",
			qualJSON:    `[{"key":"Environment","values":123}]`,
			expectError: true,
		},
		{
			name:        "invalid values type - mixed types in array",
			qualJSON:    `[{"key":"Environment","values":["prod",1]}]`,
			expectError: true,
		},
		{
			name:        "empty array (should error)",
			qualJSON:    `[]`,
			expectError: true,
		},
		{
			name:        "empty key (should error)",
			qualJSON:    `[{"key":"","values":["prod"]}]`,
			expectError: true,
		},
		{
			name:        "missing key field (should error)",
			qualJSON:    `[{"values":["prod"]}]`,
			expectError: true,
		},
		{
			name:           "empty values array (should not set Values)",
			qualJSON:       `[{"key":"Environment","values":[]}]`,
			expectedCount:  1,
			expectedKeys:   []string{"Environment"},
			expectedValues: map[string][]string{"Environment": nil},
			expectError:    false,
		},
		{
			name:           "values with empty strings (should be kept, matches API behavior)",
			qualJSON:       `[{"key":"Environment","values":["prod","","dev"]}]`,
			expectedCount:  1,
			expectedKeys:   []string{"Environment"},
			expectedValues: map[string][]string{"Environment": {"prod", "", "dev"}},
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := &plugin.QueryData{
				Quals: map[string]*plugin.KeyColumnQuals{},
			}

			if tc.qualJSON != "" {
				d.Quals["tag_filters"] = &plugin.KeyColumnQuals{
					Name: "tag_filters",
					Quals: quals.QualSlice{
						&quals.Qual{
							Column:   "tag_filters",
							Operator: "=",
							Value: &proto.QualValue{
								Value: &proto.QualValue_JsonbValue{
									JsonbValue: tc.qualJSON,
								},
							},
						},
					},
				}
			}

			result, err := buildTagFilter(d)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tc.expectedCount {
				t.Errorf("expected %d filters, got %d", tc.expectedCount, len(result))
				return
			}

			keysSeen := make(map[string]bool)
			for _, filter := range result {
				if filter.Key == nil {
					t.Errorf("filter key should not be nil")
					continue
				}

				key := *filter.Key
				keysSeen[key] = true

				expectedFound := false
				for _, expectedKey := range tc.expectedKeys {
					if key == expectedKey {
						expectedFound = true
						break
					}
				}

				if !expectedFound {
					t.Errorf("unexpected key: %s", key)
				}

				if tc.expectedValues != nil {
					expectedVals := tc.expectedValues[key]
					if expectedVals == nil && filter.Values != nil {
						t.Errorf("key %s: expected no values, got %v", key, filter.Values)
					} else if expectedVals != nil {
						if len(filter.Values) != len(expectedVals) {
							t.Errorf("key %s: expected %d values, got %d", key, len(expectedVals), len(filter.Values))
						} else {
							for i, expectedVal := range expectedVals {
								if filter.Values[i] != expectedVal {
									t.Errorf("key %s: expected value[%d] = %s, got %s", key, i, expectedVal, filter.Values[i])
								}
							}
						}
					}
				}
			}

			for _, expectedKey := range tc.expectedKeys {
				if !keysSeen[expectedKey] {
					t.Errorf("expected key %s not found in results", expectedKey)
				}
			}
		})
	}
}
