package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func TestAwsKmsKeyLastUsageTimestamp(t *testing.T) {
	timestamp := time.Date(2026, time.August, 25, 12, 0, 0, 0, time.UTC)
	testCases := []struct {
		name     string
		response interface{}
		expected *time.Time
	}{
		{
			name: "populated timestamp",
			response: &kms.GetKeyLastUsageOutput{
				KeyLastUsage: &types.KeyLastUsageData{Timestamp: &timestamp},
			},
			expected: &timestamp,
		},
		{
			name: "usage without timestamp",
			response: &kms.GetKeyLastUsageOutput{
				KeyLastUsage: &types.KeyLastUsageData{},
			},
			expected: nil,
		},
		{
			name:     "empty usage",
			response: &kms.GetKeyLastUsageOutput{},
			expected: nil,
		},
		{
			name:     "nil usage",
			response: (*kms.GetKeyLastUsageOutput)(nil),
			expected: nil,
		},
	}

	var column *transform.ColumnTransforms
	for _, candidate := range tableAwsKmsKey(context.Background()).Columns {
		if candidate.Name == "last_key_usage_timestamp" {
			column = candidate.Transform
			break
		}
	}
	if column == nil {
		t.Fatal("last_key_usage_timestamp column transform not found")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := column.Execute(context.Background(), &transform.TransformData{HydrateItem: tc.response})
			if err != nil {
				t.Fatalf("transform returned an error: %v", err)
			}
			gotTimestamp, ok := got.(*time.Time)
			if tc.expected == nil {
				if got != nil && (!ok || gotTimestamp != nil) {
					t.Errorf("transform returned %v, want nil", got)
				}
				return
			}
			if !ok || gotTimestamp == nil || !gotTimestamp.Equal(*tc.expected) {
				t.Errorf("transform returned %v, want %v", got, tc.expected)
			}
		})
	}
}
