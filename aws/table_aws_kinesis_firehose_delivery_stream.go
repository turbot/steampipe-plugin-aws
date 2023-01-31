package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"

	kinesisv1 "github.com/aws/aws-sdk-go/service/kinesis"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
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
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "InvalidParameter", "ResourceNotFoundException"}),
			},
			Hydrate: describeFirehoseDeliveryStream,
		},
		List: &plugin.ListConfig{
			Hydrate: listFirehoseDeliveryStreams,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "delivery_stream_type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(kinesisv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "delivery_stream_name",
				Description: "The name of the delivery stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the delivery stream.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamARN"),
			},
			{
				Name:        "delivery_stream_status",
				Description: "The server-side encryption type used on the stream.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "delivery_stream_type",
				Description: "The delivery stream type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "version_id",
				Description: "The version id of the stream. Each time the destination is updated for a delivery stream, the version ID is changed, and the current version ID is required when updating the destination",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "create_timestamp",
				Description: "The date and time that the delivery stream was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "has_more_destinations",
				Description: "Indicates whether there are more destinations available to list.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "last_update_timestamp",
				Description: "The date and time that the delivery stream was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "delivery_stream_encryption_configuration",
				Description: "Indicates the server-side encryption (SSE) status for the delivery stream.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "destinations",
				Description: "The destinations for the stream.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "failure_description",
				Description: "Provides details in case one of the following operations fails due to an error related to KMS: CreateDeliveryStream, DeleteDeliveryStream, StartDeliveryStreamEncryption,StopDeliveryStreamEncryption.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "source",
				Description: "If the DeliveryStreamType parameter is KinesisStreamAsSource, a SourceDescription object describing the source Kinesis data stream.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the delivery stream.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listFirehoseDeliveryStreamTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listFirehoseDeliveryStreamTags,
				Transform:   transform.FromField("Tags").Transform(kinesisFirehoseTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listFirehoseDeliveryStreams(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := FirehoseClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_firehose_delivery_stream.listFirehoseDeliveryStreams", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxLimit := int32(1000)
	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	// // List call
	param := &firehose.ListDeliveryStreamsInput{
		Limit: aws.Int32(maxLimit),
	}
	//paginator function not availablable
	equalQuals := d.KeyColumnQuals
	if equalQuals["delivery_stream_type"] != nil {
		param.DeliveryStreamType = types.DeliveryStreamType(equalQuals["delivery_stream_type"].GetStringValue())
	}
	for {
		response, err := svc.ListDeliveryStreams(ctx, param)
		if err != nil {
			plugin.Logger(ctx).Error("aws_kinesis_firehose_delivery_stream.listFirehoseDeliveryStreams", "api_error", err)
			return nil, err
		}
		for _, stream := range response.DeliveryStreamNames {
			d.StreamListItem(ctx, types.DeliveryStreamDescription{
				DeliveryStreamName: aws.String(stream),
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
	var streamName string
	if h.Item != nil {
		streamName = *h.Item.(types.DeliveryStreamDescription).DeliveryStreamName
	} else {
		quals := d.KeyColumnQuals
		streamName = quals["delivery_stream_name"].GetStringValue()
	}

	// check if streamName is empty
	if streamName == "" {
		return nil, nil
	}

	// get service
	svc, err := FirehoseClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_firehose_delivery_stream.describeFirehoseDeliveryStream", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &firehose.DescribeDeliveryStreamInput{
		DeliveryStreamName: aws.String(streamName),
	}

	// Get call
	data, err := svc.DescribeDeliveryStream(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_firehose_delivery_stream.describeFirehoseDeliveryStream", "api_error", err)
		return nil, err
	}
	return *data.DeliveryStreamDescription, nil
}

// API call for fetching tag list
func listFirehoseDeliveryStreamTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	streamName := *h.Item.(types.DeliveryStreamDescription).DeliveryStreamName

	// Create Session
	svc, err := FirehoseClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_firehose_delivery_stream.listFirehoseDeliveryStreamTags", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &firehose.ListTagsForDeliveryStreamInput{
		DeliveryStreamName: &streamName,
	}

	// Get call
	op, err := svc.ListTagsForDeliveryStream(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_firehose_delivery_stream.listFirehoseDeliveryStreamTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func kinesisFirehoseTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

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
