package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/firehose"
	pb "github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKinesisFirehoseDeliveryStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesis_firehose_delivery_stream",
		Description: "AWS Kinesis Firehose Delivery Stream",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("delivery_stream_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException", "InvalidParameter", "ResourceNotFoundException"}),
			},
			Hydrate: describeFirehoseDeliveryStream,
		},
		List: &plugin.ListConfig{
			Hydrate: listFirehoseDeliveryStreams,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "delivery_stream_type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "delivery_stream_name",
				Description: "The name of the delivery stream.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the delivery stream.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamARN"),
			},
			{
				Name:        "delivery_stream_status",
				Description: "The server-side encryption type used on the stream.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "delivery_stream_type",
				Description: "The delivery stream type.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "version_id",
				Description: "The version id of the stream. Each time the destination is updated for a delivery stream, the version ID is changed, and the current version ID is required when updating the destination",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "create_timestamp",
				Description: "The date and time that the delivery stream was created.",
				Type:        pb.ColumnType_TIMESTAMP,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "has_more_destinations",
				Description: "Indicates whether there are more destinations available to list.",
				Type:        pb.ColumnType_BOOL,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "last_update_timestamp",
				Description: "The date and time that the delivery stream was last updated.",
				Type:        pb.ColumnType_TIMESTAMP,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "delivery_stream_encryption_configuration",
				Description: "Indicates the server-side encryption (SSE) status for the delivery stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "destinations",
				Description: "The destinations for the stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "failure_description",
				Description: "Provides details in case one of the following operations fails due to an error related to KMS: CreateDeliveryStream, DeleteDeliveryStream, StartDeliveryStreamEncryption,StopDeliveryStreamEncryption.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "source",
				Description: "If the DeliveryStreamType parameter is KinesisStreamAsSource, a SourceDescription object describing the source Kinesis data stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the delivery stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     listFirehoseDeliveryStreamTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        pb.ColumnType_JSON,
				Hydrate:     listFirehoseDeliveryStreamTags,
				Transform:   transform.FromField("Tags").Transform(kinesisFirehoseTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listFirehoseDeliveryStreams(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := FirehoseService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	param := &firehose.ListDeliveryStreamsInput{
		Limit: aws.Int64(10000),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["delivery_stream_type"] != nil {
		param.DeliveryStreamType = aws.String(equalQuals["delivery_stream_type"].GetStringValue())
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *param.Limit {
			if *limit < 1 {
				param.Limit = aws.Int64(1)
			} else {
				param.Limit = limit
			}
		}
	}

	for {
		response, err := svc.ListDeliveryStreams(param)
		if err != nil {
			return nil, err
		}
		for _, stream := range response.DeliveryStreamNames {
			d.StreamListItem(ctx, &firehose.DeliveryStreamDescription{
				DeliveryStreamName: stream,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				break
			}
		}
		if !*response.HasMoreDeliveryStreams {
			break
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeFirehoseDeliveryStream(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("describeFirehoseDeliveryStream")

	var streamName string
	if h.Item != nil {
		streamName = *h.Item.(*firehose.DeliveryStreamDescription).DeliveryStreamName
	} else {
		quals := d.KeyColumnQuals
		streamName = quals["delivery_stream_name"].GetStringValue()
	}

	// check if streamName is empty
	if streamName == "" {
		return nil, nil
	}

	// get service
	svc, err := FirehoseService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &firehose.DescribeDeliveryStreamInput{
		DeliveryStreamName: aws.String(streamName),
	}

	// Get call
	data, err := svc.DescribeDeliveryStream(params)
	if err != nil {
		logger.Debug("describeDeliveryStream__", "ERROR", err)
		return nil, err
	}
	return data.DeliveryStreamDescription, nil
}

// API call for fetching tag list
func listFirehoseDeliveryStreamTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listFirehoseDeliveryStreamTags")

	streamName := *h.Item.(*firehose.DeliveryStreamDescription).DeliveryStreamName

	// Create Session
	svc, err := FirehoseService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &firehose.ListTagsForDeliveryStreamInput{
		DeliveryStreamName: &streamName,
	}

	// Get call
	op, err := svc.ListTagsForDeliveryStream(params)
	if err != nil {
		logger.Debug("listFirehoseDeliveryStreamTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func kinesisFirehoseTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("kinesisFirehoseTagListToTurbotTags")
	tagList := d.Value.([]*firehose.Tag)

	if tagList == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
