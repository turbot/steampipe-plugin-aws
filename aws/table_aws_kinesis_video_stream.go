package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo"
	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo/types"

	kinesisv1 "github.com/aws/aws-sdk-go/service/kinesis"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKinesisVideoStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesis_video_stream",
		Description: "AWS Kinesis Video Stream",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("stream_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getKinesisVideoStream,
		},
		List: &plugin.ListConfig{
			Hydrate: listKinesisVideoStreams,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(kinesisv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "stream_name",
				Description: "The name of the stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stream_arn",
				Description: "The Amazon Resource Name (ARN) of the stream.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamARN"),
			},
			{
				Name:        "status",
				Description: "The status of the stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The version of the stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of the AWS Key Management Service (AWS KMS) key that Kinesis Video Streams uses to encrypt data on the stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "A time stamp that indicates when the stream was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "data_retention_in_hours",
				Description: "How long the stream retains data, in hours.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "device_name",
				Description: "The name of the device that is associated with the stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "media_type",
				Description: "The MediaType of the stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listKinesisVideoStreamTags,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("StreamARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listKinesisVideoStreams(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKinesisVideoStreams")

	// Create session
	svc, err := KinesisVideoClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_video_stream.listKinesisVideoStreams", "connection_error", err)
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

	input := &kinesisvideo.ListStreamsInput{
		MaxResults: aws.Int32(maxLimit),
	}
	paginator := kinesisvideo.NewListStreamsPaginator(svc, input, func(o *kinesisvideo.ListStreamsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_kinesis_video_stream.listKinesisVideoStreams", "api_error", err)
			return nil, err
		}
		for _, stream := range output.StreamInfoList {
			d.StreamListItem(ctx, stream)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getKinesisVideoStream(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var streamName string
	if h.Item != nil {
		streamName = *h.Item.(types.StreamInfo).StreamName
	} else {
		quals := d.EqualsQuals
		streamName = quals["stream_name"].GetStringValue()
	}

	// get service
	svc, err := KinesisVideoClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_video_stream.getKinesisVideoStream", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &kinesisvideo.DescribeStreamInput{
		StreamName: aws.String(streamName),
	}

	// Get call
	data, err := svc.DescribeStream(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_video_stream.getKinesisVideoStream", "api_error", err)
		return nil, err
	}
	return *data.StreamInfo, nil
}

// API call for fetching tags
func listKinesisVideoStreamTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(types.StreamInfo)

	// Create Session
	svc, err := KinesisVideoClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_video_stream.listKinesisVideoStreamTags", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &kinesisvideo.ListTagsForStreamInput{
		StreamName: data.StreamName,
	}

	// Get call
	op, err := svc.ListTagsForStream(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_video_stream.listKinesisVideoStreamTags", "api_error", err)
		return nil, err
	}
	return op, nil
}
