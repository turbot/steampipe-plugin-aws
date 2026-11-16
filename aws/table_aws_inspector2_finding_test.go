package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
)

func TestResourceTagMapFilters(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    [][2]string
		wantErr bool
	}{
		{"single tag", `[{"key": "Name", "value": "Dev"}]`, [][2]string{{"Name", "Dev"}}, false},
		{"multiple tags", `[{"key": "Name", "value": "Dev"}, {"key": "Team", "value": "Ops"}]`, [][2]string{{"Name", "Dev"}, {"Team", "Ops"}}, false},
		{"missing value", `[{"key": "Name"}, {"key": "Team", "value": "Ops"}]`, nil, true},
		{"missing key", `[{"value": "Dev"}]`, nil, true},
		{"invalid JSON", `not json`, nil, true},
		{"object instead of array", `{"key": "Name", "value": "Dev"}`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resourceTagMapFilters(tt.val)
			if (err != nil) != tt.wantErr {
				t.Fatalf("resourceTagMapFilters(%q) error = %v, wantErr %v", tt.val, err, tt.wantErr)
			}
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
