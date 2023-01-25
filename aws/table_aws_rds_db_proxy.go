package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	rdsv1 "github.com/aws/aws-sdk-go/service/rds"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBProxy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_proxy",
		Description: "AWS RDS DB Proxy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("db_proxy_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DBProxyNotFoundFault"}),
			},
			Hydrate: getRDSDBProxy,
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBProxies,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAction"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "db_proxy_name",
				Description: "The identifier for the proxy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBProxyName"),
			},
			{
				Name:        "db_proxy_arn",
				Description: "The Amazon Resource Name (ARN) for the proxy",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBProxyArn"),
			},
			{
				Name:        "created_date",
				Description: "The date and time when the proxy was first created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "The current status of this proxy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "debug_logging",
				Description: "Whether the proxy includes detailed information about SQL statements in its logs.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "endpoint",
				Description: "The endpoint that you can use to connect to the DB proxy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_family",
				Description: "The kinds of databases that the proxy can connect to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "idle_client_timeout",
				Description: "The number of seconds a connection to the proxy can have no activity before the proxy drops the client connection.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "require_tls",
				Description: "Indicates whether Transport Layer Security (TLS) encryption is required for connections to the proxy.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("RequireTLS"),
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) for the IAM role that the proxy uses to access Amazon Secrets Manager.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_date",
				Description: "The date and time when the proxy was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "vpc_id",
				Description: "Provides the VPC ID of the DB proxy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auth",
				Description: "One or more data structures specifying the authorization mechanism to connect to the associated RDS DB instance or Aurora DB cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_security_group_ids",
				Description: "Provides a list of VPC security groups that the proxy belongs to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_subnet_ids",
				Description: "The EC2 subnet IDs for the proxy.",
				Type:        proto.ColumnType_JSON,
			},
			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBProxyName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRDSDBProxyTags,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBProxyArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBProxies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := RDSDBProxyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_proxy.listRDSDBProxies", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	// Limiting the results
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

	input := &rds.DescribeDBProxiesInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	paginator := rds.NewDescribeDBProxiesPaginator(svc, input, func(o *rds.DescribeDBProxiesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_proxy.listRDSDBProxies", "api_error", err)
			return nil, err
		}

		for _, items := range output.DBProxies {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBProxy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	dbProxyName := d.KeyColumnQuals["db_proxy_name"].GetStringValue()

	if dbProxyName == "" {
		return nil, nil
	}

	// Create service
	svc, err := RDSDBProxyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_proxy.getRDSDBProxy", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	params := &rds.DescribeDBProxiesInput{
		DBProxyName: aws.String(dbProxyName),
	}

	op, err := svc.DescribeDBProxies(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_proxy.getRDSDBProxy", "api_error", err)
		return nil, err
	}

	if len(op.DBProxies) > 0 {
		return op.DBProxies[0], nil
	}
	return nil, nil
}

func getRDSDBProxyTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		logger.Error("aws_rds_db_proxy.getRDSDBProxyTags", "service_creation_error", err)
		return nil, err
	}

	rdsDbProxy := h.Item.(types.DBProxy)
	params := &rds.ListTagsForResourceInput{
		ResourceName: rdsDbProxy.DBProxyArn,
	}

	tags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_proxy.getRDSDBProxyTags", "api_error", err)
		return nil, err
	}

	return tags.TagList, nil
}
