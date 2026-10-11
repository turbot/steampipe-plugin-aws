package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
)

func TestResourceTagMapFilters(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want [][2]string
	}{
		{"single tag", `[{"key": "Name", "value": "Dev"}]`, [][2]string{{"Name", "Dev"}}},
		{"multiple tags", `[{"key": "Name", "value": "Dev"}, {"key": "Team", "value": "Ops"}]`, [][2]string{{"Name", "Dev"}, {"Team", "Ops"}}},
		{"missing value is skipped", `[{"key": "Name"}, {"key": "Team", "value": "Ops"}]`, [][2]string{{"Team", "Ops"}}},
		{"missing key is skipped", `[{"value": "Dev"}]`, nil},
		{"invalid JSON", `not json`, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resourceTagMapFilters(tt.val)
			if len(got) != len(tt.want) {
				t.Fatalf("resourceTagMapFilters(%q) returned %d filters, want %d", tt.val, len(got), len(tt.want))
			}
			for i, f := range got {
				if f.Comparison != types.MapComparisonEquals {
					t.Errorf("filter %d: comparison %q, want %q", i, f.Comparison, types.MapComparisonEquals)
				}
				if aws.ToString(f.Key) != tt.want[i][0] || aws.ToString(f.Value) != tt.want[i][1] {
					t.Errorf("filter %d: got %s=%s, want %s=%s", i, aws.ToString(f.Key), aws.ToString(f.Value), tt.want[i][0], tt.want[i][1])
				}
			}
		})
	}
}
