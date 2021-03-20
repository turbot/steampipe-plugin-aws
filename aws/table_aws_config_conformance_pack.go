package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsConfigConformancePack(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_conformance_pack",
		Description: "AWS Config Conformance Pack",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"NoSuchConformancePackException"}),
			Hydrate:           getConfigConformancePack,
		},
		List: &plugin.ListConfig{
			Hydrate: listConfigConformancePacks,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "Amazon Resource Name (ARN) of the conformance pack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConformancePackArn"),
			},
			{
				Name:        "name",
				Description: "Name of the conformance pack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConformancePackName"),
			},
			{
				Name:        "conformance_pack_id",
				Description: "ID of the conformance pack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConformancePackId"),
			},
			{
				Name:        "input_parameters",
				Description: "A list of ConformancePackInputParameter objects.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConformancePackInputParameters"),
			},

			{
				Name:        "created_by",
				Description: "AWS service that created the conformance pack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CreatedBy"),
			},
			{
				Name:        "delivery_s3_bucket",
				Description: "Amazon S3 bucket where AWS Config stores conformance pack templates.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DeliveryS3Bucket"),
			},
			{
				Name:        "delivery_s3_key_prefix",
				Description: "The prefix for the Amazon S3 delivery bucket",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DeliveryS3KeyPrefix"),
			},
			{
				Name:        "last_update_requested_time",
				Description: "Last update to the conformance pack.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdateRequestedTime"),
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
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listConfigConformancePacks", "AWS_REGION", region)

	// Create session
	svc, err := ConfigService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	op, err := svc.DescribeConformancePacks(
		&configservice.DescribeConformancePacksInput{})
	if err != nil {
		return nil, err
	}
	if op.ConformancePackDetails != nil {
		for _, ConformancePackDetails := range op.ConformancePackDetails {
			d.StreamListItem(ctx, ConformancePackDetails)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getConfigConformancePack(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getConfigConformancePack")
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("matrixRegionmatrixRegion", "matrixRegion", matrixRegion)

	// Create Session
	svc, err := ConfigService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	
	params := &configservice.DescribeConformancePacksInput{
		ConformancePackNames: []*string{aws.String(name)},
	}

	op, err := svc.DescribeConformancePacks(params)
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

//// TRANSFORM FUNCTIONS


