package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/synthetics"
	"github.com/aws/aws-sdk-go-v2/service/synthetics/types"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableAwsSyntheticsCanary(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_synthetics_canary",
		Description: "AWS CloudWatch Synthetics Canary",
		Get: &plugin.GetConfig{
			// Canaries with the same name can exist in different regions for an account.
			KeyColumns: plugin.AllColumns([]string{"name", "region"}),
			Hydrate:    getSyntheticsCanary,
			Tags:       map[string]string{"service": "synthetics", "action": "GetCanary"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
			Hydrate: listSyntheticsCanaries,
			Tags:    map[string]string{"service": "synthetics", "action": "DescribeCanaries"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SYNTHETICS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the canary.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "artifact_config",
				Description: "The configuration for canary artifacts.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "artifact_s3_location",
				Description: "The S3 location where canary artifacts are stored.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ArtifactS3Location"),
			},
			{
				Name:        "browser_configs",
				Description: "The configuration for browser types used for a canary run.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "code",
				Description: "The Lambda handler and code location of the canary.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dry_run_config",
				Description: "The dry run configurations of the canary.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "engine_arn",
				Description: "The engine ARN of the canary.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_configs",
				Description: "The list of engine configurations of the canary.",
				Type:        proto.ColumnType_JSON,
			},

			{
				Name:        "execution_role_arn",
				Description: "The ARN of the execution role used to run the canary.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "failure_retention_period_in_days",
				Description: "The number of days that data for failed runs is retained.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "id",
				Description: "The unique ID of the canary.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioned_resource_cleanup",
				Description: "The resource cleanup configuration of the canary.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "run_config",
				Description: "The run configuration of the canary.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "runtime_version",
				Description: "The runtime version to use for the canary.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schedule",
				Description: "The schedule for canary runs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "status",
				Description: "The status of the canary.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "success_retention_period_in_days",
				Description: "The number of days that data for successful runs is retained.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tags",
				Description: "The list of tags assigned to the canary.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "timeline",
				Description: "The timeline information of the canary.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "visual_reference",
				Description: "The visual reference used as the baseline for screenshot comparisons.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "visual_references",
				Description: "The list of visual references used as the baseline for screenshot comparisons.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_config",
				Description: "The VPC configuration of the canary.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSyntheticsCanaryArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSyntheticsCanaries(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := SyntheticsClient(ctx, d, d.EqualsQualString(matrixKeyRegion))
	if err != nil {
		plugin.Logger(ctx).Error("aws_synthetics_canary.listSyntheticsCanaries", "client_error", err)
		return nil, err
	}

	paginator := synthetics.NewDescribeCanariesPaginator(svc, &synthetics.DescribeCanariesInput{}, func(o *synthetics.DescribeCanariesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// Apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_synthetics_canary.listSyntheticsCanaries", "api_error", err)
			return nil, err
		}

		for _, canaryDetail := range output.Canaries {
			d.StreamListItem(ctx, canaryDetail)

			// Context may be cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSyntheticsCanary(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	canaryName := d.EqualsQuals["name"].GetStringValue()

	// Get client
	svc, err := SyntheticsClient(ctx, d, d.EqualsQualString(matrixKeyRegion))
	if err != nil {
		plugin.Logger(ctx).Error("aws_synthetics_canary.getSyntheticsCanary", "client_error", err)
		return nil, err
	}

	output, err := svc.GetCanary(ctx, &synthetics.GetCanaryInput{
		Name: aws.String(canaryName),
	})
	if err != nil {
		plugin.Logger(ctx).Error("aws_synthetics_canary.getSyntheticsCanary", "api_error", err)
		return nil, err
	}

	if output.Canary == nil {
		return nil, nil
	}

	return *output.Canary, nil
}

func getSyntheticsCanaryArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	canary := h.Item.(types.Canary)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_synthetics_canary.getSyntheticsCanaryArn", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// arn format - https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazoncloudwatchsynthetics.html
	arn := "arn:" + commonColumnData.Partition + ":synthetics:" + region + ":" + commonColumnData.AccountId + ":canary:" + *canary.Name

	return arn, nil
}
