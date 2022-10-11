package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	pb "github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKinesisStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesis_stream",
		Description: "AWS Kinesis Stream",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("stream_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException", "InvalidParameter"}),
			},
			Hydrate: describeStream,
		},
		List: &plugin.ListConfig{
			Hydrate: listStreams,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "stream_name",
				Description: "The name of the stream being described.",
				Type:        pb.ColumnType_STRING,
			},
			{
				Name:        "stream_arn",
				Description: "The Amazon Resource Name (ARN) for the stream being described.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.StreamARN"),
			},
			{
				Name:        "stream_status",
				Description: "The current status of the stream being described.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.StreamStatus"),
			},
			{
				Name:        "stream_creation_timestamp",
				Description: "The approximate time that the stream was created.",
				Type:        pb.ColumnType_TIMESTAMP,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.StreamCreationTimestamp"),
			},
			{
				Name:        "encryption_type",
				Description: "The server-side encryption type used on the stream.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.EncryptionType"),
			},
			{
				Name:        "key_id",
				Description: "The GUID for the customer-managed AWS KMS key to use for encryption.",
				Type:        pb.ColumnType_STRING,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.KeyId"),
			},
			{
				Name:        "retention_period_hours",
				Description: "The current retention period, in hours.",
				Type:        pb.ColumnType_INT,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.RetentionPeriodHours"),
			},
			{
				Name:        "consumer_count",
				Description: "The number of enhanced fan-out consumers registered with the stream.",
				Type:        pb.ColumnType_INT,
				Hydrate:     describeStreamSummary,
				Transform:   transform.FromField("StreamDescriptionSummary.ConsumerCount"),
			},
			{
				Name:        "open_shard_count",
				Description: "The number of open shards in the stream.",
				Type:        pb.ColumnType_INT,
				Hydrate:     describeStreamSummary,
				Transform:   transform.FromField("StreamDescriptionSummary.OpenShardCount"),
			},
			{
				Name:        "has_more_shards",
				Description: "If set to true, more shards in the stream are available to describe.",
				Type:        pb.ColumnType_BOOL,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.HasMoreShards"),
			},
			{
				Name:        "shards",
				Description: "The shards that comprise the stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.Shards"),
			},
			{
				Name:        "enhanced_monitoring",
				Description: "Represents the current enhanced monitoring settings of the stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.EnhancedMonitoring"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the stream.",
				Type:        pb.ColumnType_JSON,
				Hydrate:     getAwsKinesisStreamTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        pb.ColumnType_STRING,
				Transform:   transform.FromField("StreamDescription.StreamName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        pb.ColumnType_JSON,
				Hydrate:     getAwsKinesisStreamTags,
				Transform:   transform.FromField("Tags").Transform(kinesisTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        pb.ColumnType_JSON,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.StreamARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listStreams(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := KinesisClient(ctx, d)
	if err != nil {
		return nil, err
	}
	pagesLeft := true
	maxLimit := int32(100)
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
	input := &kinesis.ListStreamsInput{
		Limit: aws.Int32(maxLimit),
	}
	for pagesLeft {
		result, err := svc.ListStreams(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, streams := range result.StreamNames {
			d.StreamListItem(ctx, streams)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if *result.HasMoreStreams {
			pagesLeft = true
			input.ExclusiveStartStreamName = &result.StreamNames[len(result.StreamNames)-1]
		} else {
			pagesLeft = false
		}
	}
	if err != nil {
		plugin.Logger(ctx).Error("listStreams", "ListStreams_error", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeStream(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("describeStream")

	var streamName string
	if h.Item != nil {
		streamName = h.Item.(string)
	} else {
		quals := d.KeyColumnQuals
		streamName = quals["stream_name"].GetStringValue()
	}

	// get service
	svc, err := KinesisClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &kinesis.DescribeStreamInput{
		StreamName: aws.String(streamName),
	}

	// Get call
	data, err := svc.DescribeStream(ctx, params)
	if err != nil {
		logger.Debug("describeStream__", "ERROR", err)
		return nil, err
	}
	return data, nil
}

// API call for Stream Summary
func describeStreamSummary(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("describeStreamSummary")

	streamName := getKinesisStreamName(h.Item)

	// get service
	svc, err := KinesisClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &kinesis.DescribeStreamSummaryInput{
		StreamName: aws.String(streamName),
	}

	// Get call
	data, err := svc.DescribeStreamSummary(ctx, params)
	if err != nil {
		logger.Debug("describeStreamSummary__", "ERROR", err)
		return nil, err
	}
	return data, nil
}

// API call for fetching tag list
func getAwsKinesisStreamTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsKinesisStreamTags")

	streamName := getKinesisStreamName(h.Item)

	// Create Session
	svc, err := KinesisClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &kinesis.ListTagsForStreamInput{
		StreamName: &streamName,
	}

	// Get call
	op, err := svc.ListTagsForStream(ctx, params)
	if err != nil {
		logger.Debug("getAwsKinesisStreamTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func kinesisTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("kinesisTagListToTurbotTags")
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

func getKinesisStreamName(data any) string  {
	switch item := data.(type){
	case *kinesis.DescribeStreamOutput:
		return *item.StreamDescription.StreamName
	case string:
		return item
	}
	return ""
}