package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/firehose"
	pb "github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKinesisFirehoseDeliveryStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesis_firehose_delivery_stream",
		Description: "AWS Kinesis Firehose Delivery Stream",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("delivery_stream_name"),
			ShouldIgnoreError: isNotFoundError([]string{}),
			Hydrate:           describeFirehoseDeliveryStream,
		},
		List: &plugin.ListConfig{
			Hydrate: listFirehoseDeliveryStreams,
		},
		GetMatrixItem: BuildRegionList,
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
				Transform:   transform.FromField("DeliveryStreamDescription.DeliveryStreamARN"),
			},
			{
				Name:        "delivery_stream_status",
				Description: "The server-side encryption type used on the stream.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.DeliveryStreamStatus"),
			},
			{
				Name:        "delivery_stream_type",
				Description: "The delivery stream type.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.DeliveryStreamType"),
			},
			{
				Name:        "version_id",
				Description: "The version id of the stream. Each time the destination is updated for a delivery stream, the version ID is changed, and the current version ID is required when updating the destination",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.VersionId"),
			},
			{
				Name:        "create_timestamp",
				Description: "The date and time that the delivery stream was created.",
				Type:        pb.ColumnType_TIMESTAMP,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.CreateTimestamp"),
			},
			{
				Name:        "has_more_destinations",
				Description: "Indicates whether there are more destinations available to list.",
				Type:        pb.ColumnType_BOOL,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.HasMoreDestinations"),
			},
			{
				Name:        "last_update_timestamp",
				Description: "The date and time that the delivery stream was last updated.",
				Type:        pb.ColumnType_TIMESTAMP,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.LastUpdateTimestamp"),
			},
			{
				Name:        "delivery_stream_encryption_configuration",
				Description: "Indicates the server-side encryption (SSE) status for the delivery stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.DeliveryStreamEncryptionConfiguration"),
			},
			{
				Name:        "destinations",
				Description: "The destinations for the stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.Destinations"),
			},
			{
				Name:        "failure_description",
				Description: "Provides details in case one of the following operations fails due to an error related to KMS: CreateDeliveryStream, DeleteDeliveryStream, StartDeliveryStreamEncryption,StopDeliveryStreamEncryption.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.FailureDescription"),
			},
			{
				Name:        "source",
				Description: "If the DeliveryStreamType parameter is KinesisStreamAsSource, a SourceDescription object describing the source Kinesis data stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeFirehoseDeliveryStream,
				Transform:   transform.FromField("DeliveryStreamDescription.Source"),
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
				Transform:   transform.FromField("DeliveryStreamDescription.DeliveryStreamName"),
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
				Transform:   transform.FromField("DeliveryStreamDescription.DeliveryStreamARN").Transform(arnToAkas),
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
	param := &firehose.ListDeliveryStreamsInput{}
	for {
		response, err := svc.ListDeliveryStreams(param)
		if err != nil {
			return nil, err
		}
		for _, stream := range response.DeliveryStreamNames {
			d.StreamListItem(ctx, &firehose.DeliveryStreamDescription{
				DeliveryStreamName: stream,
			})
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
	return data, nil
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
