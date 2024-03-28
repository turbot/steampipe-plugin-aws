package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	rdsv1 "github.com/aws/aws-sdk-go/service/rds"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBEngineVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_engine_version",
		Description: "AWS RDS DB Engine Version",
		List: &plugin.ListConfig{
			Hydrate: listRDSDBEngineVersions,
			Tags:    map[string]string{"service": "rds", "action": "DescribeDBEngineVersions"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "engine", Require: plugin.Optional},
				{Name: "engine_version", Require: plugin.Optional},
				{Name: "db_parameter_group_family", Require: plugin.Optional},
				{Name: "list_supported_character_sets", Require: plugin.Optional},
				{Name: "list_supported_timezones", Require: plugin.Optional},
				{Name: "default_only", Require: plugin.Optional},
				{Name: "engine_mode", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "engine",
				Description: "The name of the database engine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "engine_version",
				Description: "The version number of the database engine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the custom engine version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBEngineVersionArn"),
			},
			{
				Name:        "status",
				Description: "The status of the DB engine version, either available or deprecated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The creation time of the DB engine version.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "custom_db_engine_version_manifest",
				Description: "JSON string that lists the installation files and parameters that RDS Custom uses to create a custom engine version (CEV).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CustomDBEngineVersionManifest"),
			},
			{
				Name:        "list_supported_character_sets",
				Description: "A value that indicates whether to list the supported character sets for each engine version.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromQual("list_supported_character_sets"),
				Default:     false,
			},
			{
				Name:        "engine_mode",
				Description: "Accepts DB engine modes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("engine_mode"),
			},
			{
				Name:        "list_supported_timezones",
				Description: "A value that indicates whether to list the supported time zones for each engine version.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromQual("list_supported_timezones"),
				Default:     false,
			},
			{
				Name:        "default_only",
				Description: "A value that indicates whether only the default version of the specified engine or engine and major version combination is returned.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromQual("list_supported_timezones"),
				Default:     false,
			},
			{
				Name:        "db_engine_description",
				Description: "The description of the database engine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBEngineDescription"),
			},
			{
				Name:        "db_engine_media_type",
				Description: "A value that indicates the source media provider of the AMI based on the usage operation. Applicable for RDS Custom for SQL Server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBEngineMediaType"),
			},
			{
				Name:        "db_engine_version_description",
				Description: "The description of the database engine version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBEngineVersionDescription"),
			},
			{
				Name:        "db_parameter_group_family",
				Description: "The name of the DB parameter group family for the database engine.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBParameterGroupFamily"),
			},
			{
				Name:        "database_installation_files_s3_bucket_name",
				Description: "The name of the Amazon S3 bucket that contains your database installation files.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseInstallationFilesS3BucketName"),
			},
			{
				Name:        "database_installation_files_s3_prefix",
				Description: "The Amazon S3 directory that contains the database installation files. If not specified, then no prefix is assumed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseInstallationFilesS3Prefix"),
			},
			{
				Name:        "kms_key_id",
				Description: "The Amazon Web Services KMS key identifier for an encrypted CEV. This parameter is required for RDS Custom, but optional for Amazon RDS.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KMSKeyId"),
			},
			{
				Name:        "major_engine_version",
				Description: "The major engine version of the CEV.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "supports_babelfish",
				Description: "A value that indicates whether the engine version supports Babelfish for Aurora PostgreSQL.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "supports_certificate_rotation_without_restart",
				Description: "A value that indicates whether the engine version supports rotating the server certificate without rebooting the DB instance.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "supports_global_databases",
				Description: "A value that indicates whether you can use Aurora global databases with a specific DB engine version.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "supports_log_exports_to_cloudwatch_logs",
				Description: "A value that indicates whether the engine version supports exporting the log types specified by ExportableLogTypes to CloudWatch Logs.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "supports_parallel_query",
				Description: "A value that indicates whether you can use Aurora parallel query with a specific DB engine version.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "supports_read_replica",
				Description: "Indicates whether the database engine version supports read replicas.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "exportable_log_types",
				Description: "The types of logs that the database engine has available for export to CloudWatch Logs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "image",
				Description: "The EC2 image.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "supported_feature_names",
				Description: "A list of features supported by the DB engine.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "supported_nchar_character_sets",
				Description: "A list of the character sets supported by the Oracle DB engine for the NcharCharacterSetName parameter of the CreateDBInstance operation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SupportedNcharCharacterSets"),
			},
			{
				Name:        "supported_timezones",
				Description: "A list of the time zones supported by this engine for the Timezone parameter of the CreateDBInstance action.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "valid_upgrade_target",
				Description: "A list of engine versions that this database engine version can be upgraded to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tag_list",
				Description: "A list of tags.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EngineVersion"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TagList").Transform(getRDSDBEngineVersionTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBEngineVersionArn").Transform(getRDSDBEngineVersionAka),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBEngineVersions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_engine_version.listRDSDBEngineVersions", "connection_error", err)
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

	input := &rds.DescribeDBEngineVersionsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	if d.EqualsQualString("db_parameter_group_family") != "" {
		input.DBParameterGroupFamily = aws.String(d.EqualsQualString("db_parameter_group_family"))
	}
	if d.EqualsQualString("engine") != "" {
		input.Engine = aws.String(d.EqualsQualString("engine"))
	}
	if d.EqualsQualString("engine_version") != "" {
		input.EngineVersion = aws.String(d.EqualsQualString("engine_version"))
	}
	if d.EqualsQuals["list_supported_character_sets"] != nil {
		input.ListSupportedCharacterSets = aws.Bool(d.EqualsQuals["list_supported_character_sets"].GetBoolValue())
	}
	if d.EqualsQuals["list_supported_timezones"] != nil {
		input.ListSupportedTimezones = aws.Bool(d.EqualsQuals["list_supported_timezones"].GetBoolValue())
	}
	if d.EqualsQuals["default_only"] != nil {
		input.DefaultOnly = aws.Bool(d.EqualsQuals["default_only"].GetBoolValue())
	}

	// Additional input filter
	filters := buildEngineVersionInputFilter(d.EqualsQuals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := rds.NewDescribeDBEngineVersionsPaginator(svc, input, func(o *rds.DescribeDBEngineVersionsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_engine_version.listRDSDBEngineVersions", "api_error", err)
			return nil, err
		}

		for _, item := range output.DBEngineVersions {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// TRANSFORM FUNCTIONS

func getRDSDBEngineVersionTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	engineVersion := d.HydrateItem.(types.DBEngineVersion)

	if engineVersion.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range engineVersion.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}

func getRDSDBEngineVersionAka(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	engineVersion := d.HydrateItem.(types.DBEngineVersion)

	if engineVersion.DBEngineVersionArn == nil {
		return []string{}, nil
	} else {
		return transform.EnsureStringArray(ctx, d)
	}
}

//// UTILITY FUNCTION

func buildEngineVersionInputFilter(equalQuals plugin.KeyColumnEqualsQualMap) []types.Filter {
	filters := []types.Filter{}

	filterQuals := map[string]string{
		"status":      "status",
		"engine_mode": "engine-mode",
	}

	for qual, filterKey := range filterQuals {
		if equalQuals[qual] != nil {
			filter := types.Filter{}
			filter.Name = aws.String(filterKey)
			value := equalQuals[qual].GetStringValue()
			filter.Values = []string{value}
			filters = append(filters, filter)
		}
	}

	return filters
}
