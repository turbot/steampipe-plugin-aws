package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice/types"

	databasemigrationservicev1 "github.com/aws/aws-sdk-go/service/databasemigrationservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDmsEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dms_endpoint",
		Description: "AWS DMS Endpoint",
		List: &plugin.ListConfig{
			Hydrate: listDmsEndpoints,
			// The API returns an "InvalidParameterValueException" error when an attempt is made to filter results by ARN in a region where the resource is unavailable.
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundFault", "InvalidParameterValueException"}),
			},
			Tags:    map[string]string{"service": "dms", "action": "DescribeEndpoints"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "endpoint_identifier",
					Require: plugin.Optional,
				},
				{
					Name:    "arn",
					Require: plugin.Optional,
				},
				{
					Name:    "endpoint_type",
					Require: plugin.Optional,
				},
				{
					Name:    "engine_name",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getDmsEndpointTags,
				Tags: map[string]string{"service": "dms", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(databasemigrationservicev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "endpoint_identifier",
				Description: "The database endpoint identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) string that uniquely identifies the endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointArn"),
			},
			{
				Name:        "certificate_arn",
				Description: "The Amazon Resource Name (ARN) used for SSL connection to the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_name",
				Description: "The name of the database at the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint_type",
				Description: "The type of endpoint. Valid values are source and target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_display_name",
				Description: "The expanded name for the engine name. For example, if the EngineName parameter is 'aurora', this value would be 'Amazon Aurora MySQL'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_name",
				Description: "The database engine name. Valid values, depending on the EndpointType, include 'mysql', 'oracle', 'postgres', 'mariadb', 'aurora', 'aurora-postgresql', 'redshift', 's3', 'db2', 'db2-zos', 'azuredb', 'sybase', 'dynamodb', 'mongodb', 'kinesis', 'kafka', 'elasticsearch', 'documentdb', 'sqlserver', 'neptune', and 'babelfish'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "external_id",
				Description: "Value returned by a call to CreateEndpoint that can be used for cross-account validation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "external_table_definition",
				Description: "The external table definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "extra_connection_attributes",
				Description: "Additional connection attributes used to connect to the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "An KMS key identifier that is used to encrypt the connection parameters for the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server_name",
				Description: "The name of the server at the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_access_role_arn",
				Description: "The Amazon Resource Name (ARN) used by the service to access the IAM role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_mode",
				Description: "The SSL mode used to connect to the endpoint. The default value is none.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "username",
				Description: "The user name used to connect to the endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "The port value used to access the endpoint.",
				Type:        proto.ColumnType_INT,
			},
			// JSON columns
			{
				Name:        "dms_transfer_settings",
				Description: "The settings for the DMS Transfer type source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "doc_db_settings",
				Description: "Provides information that defines a DocumentDB endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "timestream_settings",
				Description: "The settings for the Amazon Timestream target endpoint. For more information, see the TimestreamSettings structure.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dynamo_db_settings",
				Description: "The settings for the DynamoDB target endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "elasticsearch_settings",
				Description: "The settings for the OpenSearch source endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "gcp_my_sql_settings",
				Description: "Settings in JSON format for the source GCP MySQL endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("GcpMySQLSettings"),
			},
			{
				Name:        "ibm_db2_settings",
				Description: "The settings for the IBM Db2 LUW source endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IBMDb2Settings"),
			},
			{
				Name:        "kafka_settings",
				Description: "The settings for the Apache Kafka target endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "kinesis_settings",
				Description: "The settings for the Amazon Kinesis target endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "microsoft_sql_server_settings",
				Description: "The settings for the Microsoft SQL Server source and target endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MicrosoftSQLServerSettings"),
			},
			{
				Name:        "mongo_db_settings",
				Description: "The settings for the MongoDB source endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "my_sql_settings",
				Description: "The settings for the MySQL source and target endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MySQLSettings"),
			},
			{
				Name:        "neptune_settings",
				Description: "The settings for the Amazon Neptune target endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "oracle_settings",
				Description: "The settings for the Oracle source and target endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "postgre_sql_settings",
				Description: "The settings for the PostgreSQL source and target endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PostgreSQLSettings"),
			},
			{
				Name:        "redis_settings",
				Description: "The settings for the Redis target endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "redshift_settings",
				Description: "Settings for the Amazon Redshift endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "s3_settings",
				Description: "The settings for the S3 target endpoint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("S3Settings"),
			},
			{
				Name:        "sybase_settings",
				Description: "The settings for the SAP ASE source and target endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags currently associated with the replication instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDmsEndpointTags,
				Transform:   transform.FromField("TagList"),
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EndpointIdentifier"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDmsEndpointTags,
				Transform:   transform.From(dmsEndpointTagListToTagsMap),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EndpointArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listDmsEndpoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DatabaseMigrationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_endpoint.listDmsEndpoints", "connection_error", err)
		return nil, err
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

	// Build the params
	input := &databasemigrationservice.DescribeEndpointsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	var filter []types.Filter

	// Additonal Filter
	if d.EqualsQualString("endpoint_identifier") != "" {
		paramFilter := types.Filter{
			Name:   aws.String("endpoint-id"),
			Values: []string{d.EqualsQualString("endpoint_identifier")},
		}
		filter = append(filter, paramFilter)
	}
	if d.EqualsQualString("arn") != "" {
		paramFilter := types.Filter{
			Name:   aws.String("endpoint-arn"),
			Values: []string{d.EqualsQualString("arn")},
		}
		filter = append(filter, paramFilter)
	}
	if d.EqualsQualString("endpoint_type") != "" {
		paramFilter := types.Filter{
			Name:   aws.String("endpoint-type"),
			Values: []string{d.EqualsQualString("endpoint_type")},
		}
		filter = append(filter, paramFilter)
	}
	if d.EqualsQualString("engine_name") != "" {
		paramFilter := types.Filter{
			Name:   aws.String("engine-name"),
			Values: []string{d.EqualsQualString("engine_name")},
		}
		filter = append(filter, paramFilter)
	}
	input.Filters = filter

	paginator := databasemigrationservice.NewDescribeEndpointsPaginator(svc, input, func(o *databasemigrationservice.DescribeEndpointsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dms_endpoint.listDmsEndpoints", "api_error", err)
			return nil, err
		}

		for _, items := range output.Endpoints {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDmsEndpointTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	endpointArn := h.Item.(types.Endpoint).EndpointArn

	// Create service
	svc, err := DatabaseMigrationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_endpoint.getDmsEndpointTags", "connection_error", err)
		return nil, err
	}

	params := &databasemigrationservice.ListTagsForResourceInput{
		ResourceArn: endpointArn,
	}

	endpointTags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dms_endpoint.getDmsEndpointTags", "api_error", err)
		return nil, err
	}

	return endpointTags, nil
}

//// TRANSFORM FUNCTIONS

func dmsEndpointTagListToTagsMap(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*databasemigrationservice.ListTagsForResourceOutput)

	// Mapping the resource tags inside turbotTags
	if data.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range data.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
