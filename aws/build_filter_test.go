package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func TestBuildFilter(t *testing.T) {
	tests := []struct {
		name          string
		equalQuals    plugin.KeyColumnEqualsQualMap
		expectedItems []string
	}{
		{
			name:          "empty quals",
			equalQuals:    plugin.KeyColumnEqualsQualMap{},
			expectedItems: []string{},
		},
		{
			name: "string quals",
			equalQuals: plugin.KeyColumnEqualsQualMap{
				"action":       &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: "ACCEPT"}},
				"log_status":   &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: "OK"}},
				"interface_id": &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: "eni-12345678"}},
			},
			expectedItems: []string{"ACCEPT", "OK", "eni-12345678"},
		},
		{
			name: "network address quals",
			equalQuals: plugin.KeyColumnEqualsQualMap{
				"src_addr": &proto.QualValue{Value: &proto.QualValue_InetValue{InetValue: &proto.Inet{Addr: "10.0.0.1"}}},
				"dst_addr": &proto.QualValue{Value: &proto.QualValue_InetValue{InetValue: &proto.Inet{Addr: "10.0.0.2"}}},
			},
			expectedItems: []string{"10.0.0.1", "10.0.0.2"},
		},
		{
			name: "port quals",
			equalQuals: plugin.KeyColumnEqualsQualMap{
				"src_port": &proto.QualValue{Value: &proto.QualValue_Int64Value{Int64Value: 80}},
				"dst_port": &proto.QualValue{Value: &proto.QualValue_Int64Value{Int64Value: 443}},
			},
			expectedItems: []string{"80", "443"},
		},
		{
			name: "mixed quals",
			equalQuals: plugin.KeyColumnEqualsQualMap{
				"action":       &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: "ACCEPT"}},
				"src_addr":     &proto.QualValue{Value: &proto.QualValue_InetValue{InetValue: &proto.Inet{Addr: "10.0.0.1"}}},
				"dst_port":     &proto.QualValue{Value: &proto.QualValue_Int64Value{Int64Value: 443}},
				"interface_id": &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: "eni-12345678"}},
			},
			expectedItems: []string{"ACCEPT", "10.0.0.1", "443", "eni-12345678"},
		},
		{
			name: "irrelevant quals",
			equalQuals: plugin.KeyColumnEqualsQualMap{
				"action":         &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: "ACCEPT"}},
				"unrelated_qual": &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: "should-not-appear"}},
			},
			expectedItems: []string{"ACCEPT"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFilter(tt.equalQuals)

			// Check that the result has the expected number of items
			assert.Equal(t, len(tt.expectedItems), len(result))

			// Check that all expected items are in the result
			for _, expectedItem := range tt.expectedItems {
				found := false
				for _, resultItem := range result {
					if resultItem == expectedItem {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected item %s not found in result", expectedItem)
			}
		})
	}
}
