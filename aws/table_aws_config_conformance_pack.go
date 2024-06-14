package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"

	configservicev1 "github.com/aws/aws-sdk-go/service/configservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsConfigConformancePack(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_conformance_pack",
		Description: "AWS Config Conformance Pack",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchConformancePackException"}),
			},
			Hydrate: getConfigConformancePack,
			Tags:    map[string]string{"service": "config", "action": "DescribeConformancePacks"},
		},
		List: &plugin.ListConfig{
			Hydrate: listConfigConformancePacks,
			Tags:    map[string]string{"service": "config", "action": "DescribeConformancePacks"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(configservicev1.EndpointsID),
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
			{
				Name:        "template_ssm_document_details",
				Description: "An object that contains the name or Amazon Resource Name (ARN) of the Amazon Web Services Systems Manager document (SSM document) and the version of the SSM document that is used to create a conformance pack.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TemplateSSMDocumentDetails"),
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
	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		input.ConformancePackNames = []string{equalQuals["name"].GetStringValue()}
	}

	paginator := configservice.NewDescribeConformancePacksPaginator(svc, input, func(o *configservice.DescribeConformancePacksPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_config_conformance_pack.listConfigConformancePacks", "api_error", err)
			return nil, err
		}
		for _, conformancePack := range output.ConformancePackDetails {
			d.StreamListItem(ctx, conformancePack)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getConfigConformancePack(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	quals := d.EqualsQuals
	name := quals["name"].GetStringValue()

	// Create Session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_conformance_pack.getConfigConformancePack", "get_client_error", err)
		return nil, err
	}

	params := &configservice.DescribeConformancePacksInput{
		ConformancePackNames: []string{name},
	}

	op, err := svc.DescribeConformancePacks(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_conformance_pack.getConfigConformancePack", "api_error", err)
		return nil, err
	}

	if len(op.ConformancePackDetails) > 0 {
		return op.ConformancePackDetails[0], nil
	}

	return nil, nil
}
