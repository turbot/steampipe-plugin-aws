package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect/types"
	"github.com/hashicorp/go-hclog"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/context_key"
)

func TestListConnectInstanceAttributesInstanceIdFilter(t *testing.T) {
	testCases := []struct {
		name         string
		qualInstance string
		instanceId   string
		expectSkip   bool
	}{
		{
			name:         "matching instance_id is listed",
			qualInstance: "instance-a",
			instanceId:   "instance-a",
			expectSkip:   false,
		},
		{
			name:         "non-matching instance_id is skipped",
			qualInstance: "instance-a",
			instanceId:   "instance-b",
			expectSkip:   true,
		},
	}

	ctx := context.WithValue(context.Background(), context_key.Logger, hclog.NewNullLogger())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := &plugin.QueryData{
				EqualsQuals: map[string]*proto.QualValue{
					"instance_id": {Value: &proto.QualValue_StringValue{StringValue: tc.qualInstance}},
				},
			}
			h := &plugin.HydrateData{
				Item: types.InstanceSummary{Id: aws.String(tc.instanceId)},
			}

			// The query data has no region, so creating the Connect client fails.
			// An error therefore means the instance passed the instance_id filter.
			_, err := listConnectInstanceAttributes(ctx, d, h)
			skipped := err == nil

			if skipped != tc.expectSkip {
				t.Errorf("expected skip=%v, got skip=%v (err: %v)", tc.expectSkip, skipped, err)
			}
		})
	}
}
