package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisvideo"
	pb "github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKinesisVideoStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesis_video_stream",
		Description: "AWS Kinesis Video Stream",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("stream_name"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getKinesisVideoStream,
		},
		List: &plugin.ListConfig{
			Hydrate: listKinesisVideoStreams,
		},
		GetMatrixItem: BuildRegionList,
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
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listKinesisVideoStreams", "AWS_REGION", region)

	// Create session
	svc, err := KinesisVideoService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListStreamsPages(
		&kinesisvideo.ListStreamsInput{},
		func(page *kinesisvideo.ListStreamsOutput, isLast bool) bool {
			for _, stream := range page.StreamInfoList {
				d.StreamListItem(ctx, stream)
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

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	var streamName string
	if h.Item != nil {
		streamName = *h.Item.(*kinesisvideo.StreamInfo).StreamName
	} else {
		quals := d.KeyColumnQuals
		streamName = quals["stream_name"].GetStringValue()
	}

	// get service
	svc, err := KinesisVideoService(ctx, d, region)
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

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	data := h.Item.(*kinesisvideo.StreamInfo)

	// Create Session
	svc, err := KinesisVideoService(ctx, d, region)
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