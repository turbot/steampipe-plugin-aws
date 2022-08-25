package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsCloudwatchLogResourcePolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_resource_policy",
		Description: "AWS CloudWatch Log Resource Policy",
		List: &plugin.ListConfig{
			Hydrate: listCloudwatchLogResourcePolicies,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "policy_name",
				Description: "The name of the resource policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_time",
				Description: "Timestamp showing when this policy was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdatedTime").Transform(convertTimestamp),
			},
			{
				Name:        "policy",
				Description: "The details of the policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyDocument"),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy document in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyDocument").Transform(unescape).Transform(policyToCanonical),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PolicyName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudwatchLogResourcePolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listCloudwatchLogResourcePolicies")

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		logger.Error("listCloudwatchLogResourcePolicies", "error_CloudWatchLogsService", err)
		return nil, err
	}

	// Set MaxItems to the maximum number allowed
	input := cloudwatchlogs.DescribeResourcePoliciesInput{
		Limit: aws.Int64(50),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			input.Limit = limit
		}
	}

	for {
		resp, err := svc.DescribeResourcePolicies(&input)
		if err != nil {
			logger.Error("listCloudwatchLogResourcePolicies", "error_DescribeResourcePolicies", err)
			return nil, err
		}

		// Stream results
		for _, policy := range resp.ResourcePolicies {
			d.StreamListItem(ctx, policy)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextToken == nil {
			break
		}

		input.NextToken = resp.NextToken
	}

	return nil, nil
}
