package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisvideo"
	pb "github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKinesisVideoStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesis_video_stream",
		Description: "AWS Kinesis Video Stream",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("stream_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getKinesisVideoStream,
		},
		List: &plugin.ListConfig{
			Hydrate: listKinesisVideoStreams,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "stream_name",
				Description: "The name of the stream.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "stream_arn",
				Description: "The Amazon Resource Name (ARN) of the stream.",
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("StreamARN"),
			},
			{
				Name:        "status",
				Description: "The status of the stream.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The version of the stream.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of the AWS Key Management Service (AWS KMS) key that Kinesis Video Streams uses to encrypt data on the stream.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "A time stamp that indicates when the stream was created.",
				Type:        pb.ColumnType_TIMESTAMP,
			},
			{
				Name:        "data_retention_in_hours",
				Description: "How long the stream retains data, in hours.",
				Type:        pb.ColumnType_INT,
			},
			{
				Name:        "device_name",
				Description: "The name of the device that is associated with the stream.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "media_type",
				Description: "The MediaType of the stream.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        pb.ColumnType_JSON,
				Hydrate:     listKinesisVideoStreamTags,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("StreamName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        pb.ColumnType_JSON,
				Transform:   transform.FromField("StreamARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listKinesisVideoStreams(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKinesisVideoStreams")

	// Create session
	svc, err := KinesisVideoService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &kinesisvideo.ListStreamsInput{
		MaxResults: aws.Int64(10000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListStreamsPages(
		input,
		func(page *kinesisvideo.ListStreamsOutput, isLast bool) bool {
			for _, stream := range page.StreamInfoList {
				d.StreamListItem(ctx, stream)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getKinesisVideoStream(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getKinesisVideoStream")

	var streamName string
	if h.Item != nil {
		streamName = *h.Item.(*kinesisvideo.StreamInfo).StreamName
	} else {
		quals := d.KeyColumnQuals
		streamName = quals["stream_name"].GetStringValue()
	}

	// get service
	svc, err := KinesisVideoService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &kinesisvideo.DescribeStreamInput{
		StreamName: aws.String(streamName),
	}

	// Get call
	data, err := svc.DescribeStream(params)
	if err != nil {
		logger.Debug("describeStream__", "ERROR", err)
		return nil, err
	}
	return data.StreamInfo, nil
}

// API call for fetching tags
func listKinesisVideoStreamTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listKinesisVideoStreamTags")

	data := h.Item.(*kinesisvideo.StreamInfo)

	// Create Session
	svc, err := KinesisVideoService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &kinesisvideo.ListTagsForStreamInput{
		StreamName: data.StreamName,
	}

	// Get call
	op, err := svc.ListTagsForStream(params)
	if err != nil {
		logger.Debug("listKinesisVideoStreamTags", "ERROR", err)
		return nil, err
	}
	return op, nil
}
