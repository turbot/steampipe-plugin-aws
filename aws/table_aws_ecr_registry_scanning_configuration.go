package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"

	ecrv1 "github.com/aws/aws-sdk-go/service/ecr"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcrRegistryScanningConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecr_registry_scanning_configuration",
		Description: "AWS ECR Registry Scanning Configuration",
		List: &plugin.ListConfig{
			Hydrate: getEcrRegistryScanningConfiguration,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ecrv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "registry_id",
				Description: "The ID of the registry.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scanning_configuration",
				Description: "The scanning configuration for the registry.",
				Type:        proto.ColumnType_JSON,
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryId"),
			},
		}),
	}
}

func getEcrRegistryScanningConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_registry_scanning_configuration.getEcrRegistryScanningConfiguration", "connection_error", err)
		return nil, err
	}

	params := &ecr.GetRegistryScanningConfigurationInput{}

	op, err := svc.GetRegistryScanningConfiguration(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_registry_scanning_configuration.getEcrRegistryScanningConfiguration", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, op)

	return nil, nil
}
