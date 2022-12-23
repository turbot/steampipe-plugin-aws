package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/securitylake"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityLakeDataLake(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securitylake_data_lake",
		Description: "AWS Security Lake Data Lake",
		List: &plugin.ListConfig{
			Hydrate: getSecurityLakeDataLake,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "encryption_key",
				Description: "The type of encryption key used by Security Lake to encrypt the lake configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_role_arn",
				Description: "Replication settings for the Amazon S3 buckets. This parameter uses the IAM role created by you that is managed by Security Lake, to ensure the replication setting is correct.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "s3_bucket_arn",
				Description: "Amazon Resource Names (ARNs) uniquely identify Amazon Web Services resources.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Retrieves the status of the configuration operation for an account in Amazon Security Lake.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replication_destination_regions",
				Description: "Replication enables automatic, asynchronous copying of objects across Amazon S3 buckets.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "retention_settings",
				Description: "Retention settings for the destination Amazon S3 buckets.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TagsMap"),
			},
		}),
	}
}

//// LIST FUNCTION

func getSecurityLakeDataLake(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Client
	svc, err := SecurityLakeClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securitylake_data_lake.listSecurityLakeDataLakes", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &securitylake.GetDatalakeInput{}

	resp, err := svc.GetDatalake(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securitylake_data_lake.listSecurityLakeDataLakes", "api_error", err)
		return nil, err
	}

	for _, v := range resp.Configurations {
		d.StreamLeafListItem(ctx, v)
	}

	return nil, nil
}
