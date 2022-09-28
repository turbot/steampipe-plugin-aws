package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsConfigConformancePack(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_conformance_pack",
		Description: "AWS Config Conformance Pack",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"NoSuchConformancePackException"}),
			},
			Hydrate: getConfigConformancePack,
		},
		List: &plugin.ListConfig{
			Hydrate: listConfigConformancePacks,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the conformance pack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConformancePackName"),
			},
			{
				Name:        "arn",
				Description: "Amazon Resource Name (ARN) of the conformance pack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConformancePackArn"),
			},
			{
				Name:        "conformance_pack_id",
				Description: "ID of the conformance pack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_by",
				Description: "AWS service that created the conformance pack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delivery_s3_bucket",
				Description: "Amazon S3 bucket where AWS Config stores conformance pack templates.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delivery_s3_key_prefix",
				Description: "The prefix for the Amazon S3 delivery bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_update_requested_time",
				Description: "Last update to the conformance pack.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "input_parameters",
				Description: "A list of ConformancePackInputParameter objects.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConformancePackInputParameters"),
			},

			// Standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConformancePackArn").Transform(arnToAkas),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConformancePackName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listConfigConformancePacks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_conformance_pack.listConfigConformancePacks", "get_client_error", err)
		return nil, err
	}

	input := &configservice.DescribeConformancePacksInput{
		Limit: int32(20),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(input.Limit) {
			input.Limit = int32(*limit)
		}
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.ConformancePackNames = []string{equalQuals["name"].GetStringValue()}
	}

	paginator := configservice.NewDescribeConformancePacksPaginator(svc, input, func(o *configservice.DescribeConformancePacksPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sns_topic.listAwsSnsTopics", "api_error", err)
			return nil, err
		}
		for _, conformancePack := range output.ConformancePackDetails {
			d.StreamListItem(ctx, conformancePack)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getConfigConformancePack(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	
	logger := plugin.Logger(ctx)
	logger.Trace("getConfigConformancePack")
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()

	// Create Session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_conformance_pack.getConfigConformancePack", "get_client_error", err)
		return nil, err
	}

	params := &configservice.DescribeConformancePacksInput{
		ConformancePackNames: []string{*aws.String(name)},
	}

	op, err := svc.DescribeConformancePacks(ctx, params)
	if err != nil {
		logger.Debug("getConfigConformancePack", "ERROR", err)
		return nil, err
	}

	if op != nil {
		logger.Debug("getConfigConformancePack", "SUCCESS", op)
		return op.ConformancePackDetails[0], nil
	}

	return nil, nil
}
