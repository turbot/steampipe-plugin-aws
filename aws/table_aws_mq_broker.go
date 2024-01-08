package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mq"
	"github.com/aws/aws-sdk-go-v2/service/mq/types"

	mqv1 "github.com/aws/aws-sdk-go/service/mq"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsMQBroker(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_mq_broker",
		Description: "AWS MQ Broker",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("broker_id"),
			Hydrate:    getMQBroker,
			Tags:       map[string]string{"service": "mq", "action": "DescribeBroker"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMQBrokers,
			Tags:    map[string]string{"service": "mq", "action": "ListBrokers"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getMQBroker,
				Tags: map[string]string{"service": "mq", "action": "DescribeBroker"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(mqv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "broker_name",
				Description: "The broker's name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "broker_id",
				Description: "The unique ID that Amazon MQ generates for the broker.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the broker.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BrokerArn"),
			},
			{
				Name:        "broker_state",
				Description: "The broker's status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deployment_mode",
				Description: "The broker's deployment mode.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created",
				Description: "The time when the broker was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "host_instance_type",
				Description: "The broker's instance type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "authentication_strategy",
				Description: "The authentication strategy used to secure the broker. The default is SIMPLE.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "data_replication_mode",
				Description: "Describes whether this broker is a part of a data replication pair.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "engine_type",
				Description: "The type of broker engine. Currently, Amazon MQ supports ACTIVEMQ and RABBITMQ.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "The broker engine's version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "pending_authentication_strategy",
				Description: "The authentication strategy that will be applied when the broker is rebooted. The default is SIMPLE.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "pending_data_replication_mode",
				Description: "Describes whether this broker will be a part of a data replication pair after reboot.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "pending_engine_version",
				Description: "The broker engine version to upgrade to.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "pending_host_instance_type",
				Description: "The broker's host instance type to upgrade to.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "publicly_accessible",
				Description: "Enables connections from applications outside of the VPC that hosts the broker's subnets.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "storage_type",
				Description: "The broker's storage type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "auto_minor_version_upgrade",
				Description: "Enables automatic upgrades to new minor versions for brokers, as new versions are released and supported by Amazon MQ.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getMQBroker,
			},

			// JSON columns
			{
				Name:        "actions_required",
				Description: "Actions required for a broker.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "broker_instances",
				Description: "A list of information about allocated brokers.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "configurations",
				Description: "The list of all revisions for the specified configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "data_replication_metadata",
				Description: "The replication details of the data replication-enabled broker. Only returned if dataReplicationMode is set to CRDR.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "encryption_options",
				Description: "Encryption options for the broker.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "ldap_server_metadata",
				Description: "The metadata of the LDAP server used to authenticate and authorize connections to the broker.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "logs",
				Description: "The list of information about logs currently enabled and pending to be deployed for the specified broker.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "pending_ldap_server_metadata",
				Description: "The metadata of the LDAP server that will be used to authenticate and authorize connections to the broker after it is rebooted.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "maintenance_window_start_time",
				Description: "The parameters that determine the WeeklyStartTime.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "pending_data_replication_metadata",
				Description: "The pending replication details of the data replication-enabled broker. Only returned if pendingDataReplicationMode is set to CRDR.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "pending_security_groups",
				Description: "The list of pending security groups to authorize connections to brokers.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "security_groups",
				Description: "The list of rules (1 minimum, 125 maximum) that authorize connections to brokers.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "subnet_ids",
				Description: "The list of groups that define which subnets and IP ranges the broker can use from different Availability Zones.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "users",
				Description: "The list of all broker usernames for the specified broker.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: "A list of tags attached to the broker.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMQBroker,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BrokerName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BrokerArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listMQBrokers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := MQClient(ctx, d)
	if err != nil {
		logger.Error("aws_mq_broker.listMQBrokers", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 5 {
				maxLimit = 5
			} else {
				maxLimit = limit
			}
		}
	}

	input := mq.ListBrokersInput{
		MaxResults: maxLimit,
	}

	paginator := mq.NewListBrokersPaginator(svc, &input, func(o *mq.ListBrokersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_mq_broker.listMQBrokers", "api_error", err)
			return nil, err
		}

		for _, broker := range output.BrokerSummaries {
			d.StreamListItem(ctx, broker)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getMQBroker(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	brokerId := ""
	if h.Item != nil {
		broker := h.Item.(types.BrokerSummary)
		brokerId = *broker.BrokerId
	} else {
		brokerId = d.EqualsQualString("broker_id")
	}
	if brokerId == "" {
		return nil, nil
	}

	// Create service
	svc, err := MQClient(ctx, d)
	if err != nil {
		logger.Error("aws_mq_broker.getMQBroker", "service_creation_error", err)
		return nil, err
	}

	// Unsupported region, return no data
	if svc == nil {
		return nil, nil
	}

	params := &mq.DescribeBrokerInput{
		BrokerId: aws.String(brokerId),
	}

	op, err := svc.DescribeBroker(ctx, params)
	if err != nil {
		logger.Error("aws_mq_broker.getMQBroker", "api_error", err)
		return nil, err
	}

	return op, nil
}