package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"

	kinesisEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKinesisStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kinesis_stream",
		Description: "AWS Kinesis Stream",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("stream_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameter"}),
			},
			Hydrate: describeStream,
			Tags:    map[string]string{"service": "kinesis", "action": "DescribeStream"},
		},
		List: &plugin.ListConfig{
			Hydrate: listStreams,
			Tags:    map[string]string{"service": "kinesis", "action": "ListStreams"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(kinesisEndpoint.AWS_KINESIS_SERVICE_ID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: describeStream,
				Tags: map[string]string{"service": "kinesis", "action": "DescribeStream"},
			},
			{
				Func: describeStreamSummary,
				Tags: map[string]string{"service": "kinesis", "action": "DescribeStreamSummary"},
			},
			{
				Func: getAwsKinesisStreamTags,
				Tags: map[string]string{"service": "kinesis", "action": "ListTagsForStream"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "stream_name",
				Description: "The name of the stream being described.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamDescription.StreamName"),
			},
			{
				Name:        "stream_arn",
				Description: "The Amazon Resource Name (ARN) for the stream being described.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.StreamARN"),
			},
			{
				Name:        "stream_status",
				Description: "The current status of the stream being described.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.StreamStatus"),
			},
			{
				Name:        "stream_creation_timestamp",
				Description: "The approximate time that the stream was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.StreamCreationTimestamp"),
			},
			{
				Name:        "encryption_type",
				Description: "The server-side encryption type used on the stream.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.EncryptionType"),
			},
			{
				Name:        "key_id",
				Description: "The GUID for the customer-managed AWS KMS key to use for encryption.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.KeyId"),
			},
			{
				Name:        "retention_period_hours",
				Description: "The current retention period, in hours.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.RetentionPeriodHours"),
			},
			{
				Name:        "consumer_count",
				Description: "The number of enhanced fan-out consumers registered with the stream.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeStreamSummary,
				Transform:   transform.FromField("StreamDescriptionSummary.ConsumerCount"),
			},
			{
				Name:        "open_shard_count",
				Description: "The number of open shards in the stream.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeStreamSummary,
				Transform:   transform.FromField("StreamDescriptionSummary.OpenShardCount"),
			},
			{
				Name:        "has_more_shards",
				Description: "If set to true, more shards in the stream are available to describe.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.HasMoreShards"),
			},
			{
				Name:        "shards",
				Description: "The shards that comprise the stream.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.Shards"),
			},
			{
				Name:        "enhanced_monitoring",
				Description: "Represents the current enhanced monitoring settings of the stream.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.EnhancedMonitoring"),
			},
			{
				Name:        "stream_mode_details",
				Description: "Represents the current mode of the stream.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeStream,
				Transform:   transform.FromField("StreamDescription.StreamModeDetails"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the stream.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKinesisStreamTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamDescription.StreamName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKinesisStreamTags,
				Transform:   transform.FromField("Tags").Transform(kinesisTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
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
		plugin.Logger(ctx).Error("aws_kinesis_stream.listStreams", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
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

	// API doesn't support aws-sdk-go-v2 paginator as of date
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		result, err := svc.ListStreams(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_kinesis_stream.listStreams", "api_error", err)
			return nil, err
		}

		for _, streams := range result.StreamNames {
			d.StreamListItem(ctx, &kinesis.DescribeStreamOutput{
				StreamDescription: &types.StreamDescription{
					StreamName: aws.String(streams),
				},
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
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

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeStream(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var streamName string
	if h.Item != nil {
		streamName = *h.Item.(*kinesis.DescribeStreamOutput).StreamDescription.StreamName
	} else {
		quals := d.EqualsQuals
		streamName = quals["stream_name"].GetStringValue()
	}

	// get service
	svc, err := KinesisClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_stream.describeStream", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &kinesis.DescribeStreamInput{
		StreamName: aws.String(streamName),
	}

	// Get call
	data, err := svc.DescribeStream(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_stream.describeStream", "api_error", err)
		return nil, err
	}
	return data, nil
}

// API call for Stream Summary
func describeStreamSummary(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	streamName := *h.Item.(*kinesis.DescribeStreamOutput).StreamDescription.StreamName

	// get service
	svc, err := KinesisClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_stream.describeStreamSummary", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &kinesis.DescribeStreamSummaryInput{
		StreamName: aws.String(streamName),
	}

	// Get call
	data, err := svc.DescribeStreamSummary(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_stream.describeStreamSummary", "api_error", err)
		return nil, err
	}
	return data, nil
}

// API call for fetching tag list
func getAwsKinesisStreamTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	streamName := *h.Item.(*kinesis.DescribeStreamOutput).StreamDescription.StreamName

	// Create Session
	svc, err := KinesisClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_stream.getAwsKinesisStreamTags", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &kinesis.ListTagsForStreamInput{
		StreamName: &streamName,
	}

	// Get call
	op, err := svc.ListTagsForStream(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kinesis_stream.getAwsKinesisStreamTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func kinesisTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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
