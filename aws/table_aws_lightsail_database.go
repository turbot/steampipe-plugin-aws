package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"

	lightsailv1 "github.com/aws/aws-sdk-go/service/lightsail"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLightsailDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lightsail_database",
		Description: "AWS Lightsail Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getLightsailDatabase,
			Tags:       map[string]string{"service": "lightsail", "action": "GetRelationalDatabase"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidResourceName", "DoesNotExist"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listLightsailDatabases,
			Tags:    map[string]string{"service": "lightsail", "action": "GetRelationalDatabases"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lightsailv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arn"),
			},
			{
				Name:        "backup_retention_enabled",
				Description: "A Boolean value indicating whether automated backup retention is enabled for the database.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("BackupRetentionEnabled"),
			},
			{
				Name:        "ca_certificate_identifier",
				Description: "The certificate associated with the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CaCertificateIdentifier"),
			},
			{
				Name:        "created_at",
				Description: "The timestamp when the database was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreatedAt"),
			},
			{
				Name:        "engine",
				Description: "The database engine (e.g., MySQL).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Engine"),
			},
			{
				Name:        "engine_version",
				Description: "The database engine version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EngineVersion"),
			},
			{
				Name:        "latest_restorable_time",
				Description: "The latest point in time to which the database can be restored.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LatestRestorableTime"),
			},
			{
				Name:        "master_database_name",
				Description: "The name of the master database created when the Lightsail database resource is created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MasterDatabaseName"),
			},
			{
				Name:        "master_endpoint",
				Description: "The master endpoint of the database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MasterEndpoint"),
			},
			{
				Name:        "master_username",
				Description: "The master user name of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MasterUsername"),
			},
			{
				Name:        "preferred_backup_window",
				Description: "The daily time range during which automated backups are created for the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PreferredBackupWindow"),
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "The weekly time range during which system maintenance can occur on the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PreferredMaintenanceWindow"),
			},
			{
				Name:        "publicly_accessible",
				Description: "Specifies the accessibility options for the database.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PubliclyAccessible"),
			},
			{
				Name:        "relational_database_blueprint_id",
				Description: "The blueprint ID for the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RelationalDatabaseBlueprintId"),
			},
			{
				Name:        "relational_database_bundle_id",
				Description: "The bundle ID for the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RelationalDatabaseBundleId"),
			},
			{
				Name:        "secondary_availability_zone",
				Description: "Describes the secondary Availability Zone of a high availability database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecondaryAvailabilityZone"),
			},
			{
				Name:        "state",
				Description: "The current state of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State"),
			},
			{
				Name:        "support_code",
				Description: "The support code for the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SupportCode"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the database.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(getLightsailDatabaseTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listLightsailDatabases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := LightsailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_database.listLightsailDatabases", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &lightsail.GetRelationalDatabasesInput{}

	// List call
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		resp, err := svc.GetRelationalDatabases(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lightsail_database.listLightsailDatabases", "query_error", err)
			return nil, err
		}

		for _, db := range resp.RelationalDatabases {
			d.StreamListItem(ctx, db)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPageToken != nil {
			input.PageToken = resp.NextPageToken
		} else {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLightsailDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := LightsailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_database.getLightsailDatabase", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var name string
	if h.Item != nil {
		name = *h.Item.(types.RelationalDatabase).Name
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Empty check
	if name == "" {
		return nil, nil
	}

	params := &lightsail.GetRelationalDatabaseInput{
		RelationalDatabaseName: aws.String(name),
	}

	op, err := svc.GetRelationalDatabase(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_database.getLightsailDatabase", "query_error", err)
		return nil, err
	}

	return op.RelationalDatabase, nil
}

func getLightsailDatabaseTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.Value.([]types.Tag)
	if tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
