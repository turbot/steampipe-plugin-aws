package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/shield"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsShieldDRTAccess(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_drt_access",
		Description: "AWS Shield DRT Access",
		List: &plugin.ListConfig{
			Hydrate: listAwsShieldDRTAccess,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags:    map[string]string{"service": "shield", "action": "DescribeDRTAccess"},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the role the SRT used to access your AWS account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_bucket_list",
				Description: "The list of Amazon S3 buckets accessed by the SRT.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LogBucketList").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listAwsShieldDRTAccess(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_drt_access.listAwsShieldDRTAccess", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	data, err := svc.DescribeDRTAccess(ctx, &shield.DescribeDRTAccessInput{})

	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_protection.getAwsProtection", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, data)

	return nil, nil
}
