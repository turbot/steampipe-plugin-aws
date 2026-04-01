package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsMSKTopic(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_msk_topic",
		Description: "AWS MSK Topic",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"cluster_arn", "topic_name"}),
			Hydrate:    getMSKTopic,
			Tags:       map[string]string{"service": "kafka", "action": "DescribeTopic"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate:    listMSKTopics,
			Tags:       map[string]string{"service": "kafka", "action": "ListTopics"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "cluster_arn", Require: plugin.Required},
				{Name: "topic_name", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: describeMSKTopic,
				Tags: map[string]string{"service": "kafka", "action": "DescribeTopic"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_KAFKA_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cluster_arn",
				Description: "The Amazon Resource Name (ARN) of the MSK cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "topic_name",
				Description: "The name of the Kafka topic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "topic_arn",
				Description: "The Amazon Resource Name (ARN) of the topic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "partition_count",
				Description: "The number of partitions for the topic.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "replication_factor",
				Description: "The replication factor for the topic.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "out_of_sync_replica_count",
				Description: "The number of out-of-sync replicas for the topic.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "status",
				Description: "The status of the topic.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeMSKTopic,
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "configs",
				Description: "Topic configurations encoded as a Base64 string.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeMSKTopic,
				Transform:   transform.FromField("Configs"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TopicName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TopicArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type mskTopicRow struct {
	ClusterArn           *string
	TopicName            *string
	TopicArn             *string
	PartitionCount       *int32
	ReplicationFactor    *int32
	OutOfSyncReplicaCount *int32
}

//// LIST FUNCTION

func listMSKTopics(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	clusterArn := d.EqualsQuals["cluster_arn"].GetStringValue()
	if clusterArn == "" {
		return nil, nil
	}

	// The cluster_arn encodes the owning account (arn:aws:kafka:REGION:ACCOUNT_ID:cluster/...).
	// Steampipe fans this query out to all configured connections, but only the connection
	// whose account matches the cluster's account can call kafka-cluster:Connect on it.
	// Skip early for non-owning connections to avoid cross-account 403 errors.
	arnParts := strings.Split(clusterArn, ":")
	if len(arnParts) >= 5 {
		clusterAccountID := arnParts[4]
		callerIdentityData, err := getCallerIdentity(ctx, d, &plugin.HydrateData{})
		if err == nil && callerIdentityData != nil {
			callerIdentity := callerIdentityData.(*sts.GetCallerIdentityOutput)
			if callerIdentity.Account != nil && *callerIdentity.Account != clusterAccountID {
				logger.Debug("aws_msk_topic.listMSKTopics", "skipping connection - account mismatch",
					"connection_account", *callerIdentity.Account, "cluster_account", clusterAccountID)
				return nil, nil
			}
		}
	}

	svc, err := KafkaClient(ctx, d)
	if err != nil {
		logger.Error("aws_msk_topic.listMSKTopics", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 20 {
				maxLimit = 20
			} else {
				maxLimit = limit
			}
		}
	}

	input := kafka.ListTopicsInput{
		ClusterArn: aws.String(clusterArn),
		MaxResults: aws.Int32(maxLimit),
	}

	if d.EqualsQuals["topic_name"] != nil {
		input.TopicNameFilter = aws.String(d.EqualsQuals["topic_name"].GetStringValue())
	}

	paginator := kafka.NewListTopicsPaginator(svc, &input, func(o *kafka.ListTopicsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			logger.Error("aws_msk_topic.listMSKTopics", "api_error", err)
			return nil, err
		}

		for _, topic := range output.Topics {
			d.StreamListItem(ctx, mskTopicRow{
				ClusterArn:            aws.String(clusterArn),
				TopicName:             topic.TopicName,
				TopicArn:              topic.TopicArn,
				PartitionCount:        topic.PartitionCount,
				ReplicationFactor:     topic.ReplicationFactor,
				OutOfSyncReplicaCount: topic.OutOfSyncReplicaCount,
			})

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getMSKTopic(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	clusterArn := d.EqualsQuals["cluster_arn"].GetStringValue()
	topicName := d.EqualsQuals["topic_name"].GetStringValue()
	if clusterArn == "" || topicName == "" {
		return nil, nil
	}

	svc, err := KafkaClient(ctx, d)
	if err != nil {
		logger.Error("aws_msk_topic.getMSKTopic", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	output, err := svc.DescribeTopic(ctx, &kafka.DescribeTopicInput{
		ClusterArn: aws.String(clusterArn),
		TopicName:  aws.String(topicName),
	})
	if err != nil {
		logger.Error("aws_msk_topic.getMSKTopic", "api_error", err)
		return nil, err
	}

	return mskTopicRow{
		ClusterArn:        aws.String(clusterArn),
		TopicName:         output.TopicName,
		TopicArn:          output.TopicArn,
		PartitionCount:    output.PartitionCount,
		ReplicationFactor: output.ReplicationFactor,
	}, nil
}

func describeMSKTopic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	row := h.Item.(mskTopicRow)

	if row.ClusterArn == nil || row.TopicName == nil {
		return nil, nil
	}

	svc, err := KafkaClient(ctx, d)
	if err != nil {
		logger.Error("aws_msk_topic.describeMSKTopic", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	output, err := svc.DescribeTopic(ctx, &kafka.DescribeTopicInput{
		ClusterArn: row.ClusterArn,
		TopicName:  row.TopicName,
	})
	if err != nil {
		logger.Error("aws_msk_topic.describeMSKTopic", "api_error", err)
		return nil, err
	}

	return output, nil
}
