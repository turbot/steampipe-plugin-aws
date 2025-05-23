package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAwsGlueDataCatalogEncryptionSettings(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_data_catalog_encryption_settings",
		Description: "AWS Glue Data Catalog Encryption Settings",
		List: &plugin.ListConfig{
			Hydrate: listGlueDataCatalogEncryptionSettings,
			Tags:    map[string]string{"service": "glue", "action": "GetDataCatalogEncryptionSettings"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_GLUE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "encryption_at_rest",
				Description: "A list of public keys to be used by the DataCatalogEncryptionSettingss for authentication.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "connection_password_encryption",
				Description: "A list of security group identifiers used in this DataCatalogEncryptionSettings.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueDataCatalogEncryptionSettings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_data_catalog_encryption_settings.listGlueDataCatalogEncryptionSettings", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &glue.GetDataCatalogEncryptionSettingsInput{}

	// List call
	result, err := svc.GetDataCatalogEncryptionSettings(ctx, input)

	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_data_catalog_encryption_settings.listGlueDataCatalogEncryptionSettings", "api_error", err)
		return nil, err
	}
	d.StreamListItem(ctx, result.DataCatalogEncryptionSettings)

	return nil, nil
}
