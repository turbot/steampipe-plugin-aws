package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueCatalogDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_catalog_database",
		Description: "AWS Glue Catalog Database",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"EntityNotFoundException", "InvalidParameter"}),
			Hydrate:           getGlueCatalogDatabase,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueCatalogDatabases,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the database. For Hive compatibility, this is folded to lowercase when it is stored.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "catalog_id",
				Description: "The ID of the Data Catalog in which the database resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time at which the metadata database was created in the catalog.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A description of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location_uri",
				Description: "The location of the database (for example, an HDFS path).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy_create_time",
				Description: "The date and time at which the policy was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getGlueCatalogDatabasePolicy,
				Transform:   transform.FromField("CreateTime"),
			},
			{
				Name:        "policy_hash",
				Description: "Contains the hash value associated with this policy.",
				Hydrate:     getGlueCatalogDatabasePolicy,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy_update_time",
				Description: "The date and time at which the policy was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getGlueCatalogDatabasePolicy,
				Transform:   transform.FromField("UpdateTime"),
			},
			{
				Name:        "create_table_default_permissions",
				Description: "Creates a set of default permissions on the table for principals.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "parameters",
				Description: "These key-value pairs define parameters and properties of the database.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy",
				Description: "Contains the requested policy document",
				Hydrate:     getGlueCatalogDatabasePolicy,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyInJson"),
			},
			{
				Name:        "policy_std",
				Description: "Contains the contents of the requested policy document in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueCatalogDatabasePolicy,
				Transform:   transform.FromField("PolicyInJson").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "target_database",
				Description: "A DatabaseIdentifier structure that describes a target database for resource linking.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueCatalogDatabaseAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueCatalogDatabases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GlueService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listGlueCatalogDatabases", "error_GlueService", err)
		return nil, err
	}

	// List call
	err = svc.GetDatabasesPages(
		&glue.GetDatabasesInput{},
		func(page *glue.GetDatabasesOutput, isLast bool) bool {
			for _, database := range page.DatabaseList {
				d.StreamListItem(ctx, database)

			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listGlueCatalogDatabases", "error_GetDatabasesPages", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueCatalogDatabase(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGlueCatalogDatabase")

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := GlueService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getGlueCatalogDatabase", "error_GlueService", err)
		return nil, err
	}

	// Build the params
	params := &glue.GetDatabaseInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.GetDatabase(params)
	if err != nil {
		plugin.Logger(ctx).Error("getGlueCatalogDatabase", "error_GetDatabase", err)
		return nil, err
	}

	return data.Database, nil
}

func getGlueCatalogDatabaseAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGlueCatalogDatabaseAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*glue.Database)

	// Get common columns
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getGlueCatalogDatabaseAkas", "error_getCommonColumnsCached", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":database/" + *data.Name

	return []string{aka}, nil
}

func getGlueCatalogDatabasePolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGlueCatalogDatabasePolicy")

	// Create Session
	svc, err := GlueService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getGlueCatalogDatabasePolicy", "error_GlueService", err)
		return nil, err
	}

	// Get call
	data, err := svc.GetResourcePolicy(&glue.GetResourcePolicyInput{})
	if err != nil {
		plugin.Logger(ctx).Error("getGlueCatalogDatabasePolicy", "error_GetResourcePolicy", err)
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "EntityNotFoundException" {
				return nil, nil
			}
		}
		return nil, err
	}

	return data, nil
}
