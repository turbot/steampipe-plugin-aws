package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/stretchr/testify/assert"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func TestGetMessageField(t *testing.T) {
	tests := []struct {
		name     string
		item     interface{}
		expected []string
		wantErr  bool
	}{
		{
			name: "CloudWatch log event",
			item: types.FilteredLogEvent{
				Message: aws.String("2 123456789012 eni-12345678 10.0.0.1 10.0.0.2 80 443 6 10 100 1600000000 1600000001 ACCEPT OK"),
			},
			expected: []string{"2", "123456789012", "eni-12345678", "10.0.0.1", "10.0.0.2", "80", "443", "6", "10", "100", "1600000000", "1600000001", "ACCEPT", "OK"},
			wantErr:  false,
		},
		{
			name: "S3 flow log event",
			item: s3FlowLogEvent{
				FilteredLogEvent: types.FilteredLogEvent{
					Message: aws.String("2 123456789012 eni-12345678 10.0.0.1 10.0.0.2 80 443 6 10 100 1600000000 1600000001 ACCEPT OK"),
				},
			},
			expected: []string{"2", "123456789012", "eni-12345678", "10.0.0.1", "10.0.0.2", "80", "443", "6", "10", "100", "1600000000", "1600000001", "ACCEPT", "OK"},
			wantErr:  false,
		},
		{
			name:     "Unknown item type",
			item:     struct{}{},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getMessageField(context.Background(), nil, &plugin.HydrateData{Item: tt.item})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetField(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		idx      int
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "Valid field",
			fields:   []string{"2", "123456789012", "eni-12345678", "10.0.0.1"},
			idx:      2,
			expected: "eni-12345678",
			wantErr:  false,
		},
		{
			name:     "Dash field (treated as null)",
			fields:   []string{"2", "-", "eni-12345678"},
			idx:      1,
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getField(context.Background(), &transform.TransformData{
				Value: tt.fields,
				Param: tt.idx,
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTableAwsVpcFlowLogEventListKeyColumns(t *testing.T) {
	columns := tableAwsVpcFlowLogEventListKeyColumns()
	assert.NotNil(t, columns)

	// Verify we have the expected columns
	expectedColumns := map[string]bool{
		"log_source":      true,
		"log_group_name":  true,
		"bucket_name":     true,
		"s3_prefix":       true,
		"log_stream_name": true,
		"filter":          true,
		"region":          true,
		"timestamp":       true,
		"event_id":        true,
		"interface_id":    true,
		"src_addr":        true,
		"dst_addr":        true,
		"src_port":        true,
		"dst_port":        true,
		"action":          true,
		"log_status":      true,
	}

	for _, col := range columns {
		assert.True(t, expectedColumns[col.Name], "Column %s should be in the expected columns", col.Name)
	}
}
