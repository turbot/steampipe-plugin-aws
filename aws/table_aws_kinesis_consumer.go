package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKinesisConsumer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesis_consumer",
		Description: "AWS Kinesis Consumer",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("consumer_arn"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getAwsKinesisConsumer,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStreams,
			Hydrate:       listKinesisConsumers,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "consumer_name",
				Description: "The name of the consumer is something you choose when you register the consumer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "consumer_arn",
				Description: "When you register a consumer, Kinesis Data Streams generates an ARN for it.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConsumerARN"),
			},
			{
				Name:        "stream_arn",
				Description: "The ARN of the stream with which you registered the consumer.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKinesisConsumer,
				Transform:   transform.FromField("StreamARN"),
			},
			{
				Name:        "consumer_status",
				Description: "A consumer can't read data while in the CREATING or DELETING states.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "consumer_creation_timestamp",
				Description: "Timestamp when consumer was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConsumerName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConsumerARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listKinesisConsumers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
  var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listKinesisConsumers", "AWS_REGION", region)
	streamData := *h.Item.(*kinesis.DescribeStreamOutput)

	c, err := getCommonColumns(ctx, d, h)

	commonColumnData := c.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":kinesis:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":stream" + "/" + *streamData.StreamDescription.StreamName

	plugin.Logger(ctx).Trace("StreamArn", "arn", arn)

	// Create session
	svc, err := KinesisService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	err = svc.ListStreamConsumersPages(
		&kinesis.ListStreamConsumersInput{StreamARN: &arn},
		func(page *kinesis.ListStreamConsumersOutput, isLast bool) bool {
			for _, consumerData := range page.Consumers {
				d.StreamLeafListItem(ctx, consumerData)
			}
			return !isLast
		},
	)

	return nil, err
}

func getAwsKinesisConsumer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsKinesisConsumer")
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var arn string
	if h.Item != nil {
		i := h.Item.(*kinesis.Consumer)
		arn = *i.ConsumerARN
	} else {
		arn = d.KeyColumnQuals["consumer_arn"].GetStringValue()
	}

	// Create Session
	svc, err := KinesisService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &kinesis.DescribeStreamConsumerInput{
		ConsumerARN: &arn,
	}

	// Get call
	data, err := svc.DescribeStreamConsumer(params)
	if err != nil {
		logger.Debug("getAwsKinesisConsumer", "ERROR", err)
		return nil, err
	}

	return data.ConsumerDescription, nil
}
