package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"

	configservicev1 "github.com/aws/aws-sdk-go/service/configservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsConfigDeliveryChannel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_delivery_channel",
		Description: "AWS Config Delivery Channel",
		List: &plugin.ListConfig{
			Hydrate: listConfigDeliveryChannels,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},		
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchDeliveryChannelException"}),
			},
			Tags:    map[string]string{"service": "config", "action": "DescribeDeliveryChannels"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(configservicev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the delivery channel.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "s3_bucket_name",
				Description: "The name of the Amazon S3 bucket to which AWS Config delivers configuration snapshots and configuration history files.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("S3BucketName"),
			},
			{
				Name:        "s3_key_prefix",
				Description: "The prefix for the specified Amazon S3 bucket.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("S3KeyPrefix"),
			},
			{
				Name:        "s3_kms_key_arn",
				Description: "The Amazon Resource Name (ARN) of the KMS key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("S3KmsKeyArn"),
			},
			{
				Name:        "sns_topic_arn",
				Description: "The Amazon Resource Name (ARN) of the Amazon SNS topic to which AWS Config sends notifications about configuration changes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnsTopicARN"),
			},
			{
				Name:        "delivery_frequency",
				Description: "The frequency with which the AWS Config delivers configuration snapshots to the Amazon S3 bucket.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConfigSnapshotDeliveryProperties.DeliveryFrequency"),
			},
			{
				Name:        "status",
				Description: "The current status of the delivery channel.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConfigDeliveryChannelStatus,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsDeliveryChannelAkas,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listConfigDeliveryChannels(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_delivery_channel.listConfigDeliveryChannels", "get_client_error", err)
		return nil, err
	}

	input := &configservice.DescribeDeliveryChannelsInput{}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		input.DeliveryChannelNames = []string{equalQuals["name"].GetStringValue()}
	}

	op, err := svc.DescribeDeliveryChannels(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_delivery_channel.listConfigDeliveryChannels", "api_error", err)
		return nil, err
	}

	if op.DeliveryChannels != nil {
		for _, deliveryChannel := range op.DeliveryChannels {
			d.StreamListItem(ctx, deliveryChannel)

			// Context can be cancelled due to manual cancellation or limit being reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getConfigDeliveryChannelStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	deliveryChannel := h.Item.(types.DeliveryChannel)

	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_delivery_channel.getConfigDeliveryChannelStatus", "get_client_error", err)
		return nil, err
	}

	params := &configservice.DescribeDeliveryChannelStatusInput{
		DeliveryChannelNames: []string{*deliveryChannel.Name},
	}

	status, err := svc.DescribeDeliveryChannelStatus(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_delivery_channel.getConfigDeliveryChannelStatus", "api_error", err)
		return nil, err
	}

	if len(status.DeliveryChannelsStatus) < 1 {
		return nil, nil
	}

	return status.DeliveryChannelsStatus[0], nil
}

func getAwsDeliveryChannelAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	deliveryChannel := h.Item.(types.DeliveryChannel)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_delivery_channel.getAwsDeliveryChannelAkas", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":config:" + region + ":" + commonColumnData.AccountId + ":delivery-channel" + "/" + *deliveryChannel.Name

	return arn, nil
}
