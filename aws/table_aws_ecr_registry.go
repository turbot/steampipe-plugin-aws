package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcrRegistry(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecr_registry",
		Description: "AWS ECR Private Registry",
		List: &plugin.ListConfig{
			Hydrate: listEcrRegistry,
			Tags:    map[string]string{"service": "ecr", "action": "DescribeRegistry"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getEcrRegistryPolicy,
				Tags: map[string]string{"service": "ecr", "action": "GetRegistryPolicy"},
				IgnoreConfig: &plugin.IgnoreConfig{
					ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"RegistryPolicyNotFoundException"}),
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_API_ECR_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "registry_id",
				Description: "The ID of the registry.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the registry.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEcrRegistryArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "replication_configuration",
				Description: "The replication configuration for the registry.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy",
				Description: "The JSON text of the registry permissions policy.",
				Hydrate:     getEcrRegistryPolicy,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyText"),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Hydrate:     getEcrRegistryPolicy,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PolicyText").Transform(unescape).Transform(policyToCanonical),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcrRegistryArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEcrRegistry(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_registry.listEcrRegistry", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	op, err := svc.DescribeRegistry(ctx, &ecr.DescribeRegistryInput{})
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_registry.listEcrRegistry", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, op)

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEcrRegistryPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ECRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_registry.getEcrRegistryPolicy", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	op, err := svc.GetRegistryPolicy(ctx, &ecr.GetRegistryPolicyInput{})
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_registry.getEcrRegistryPolicy", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getEcrRegistryArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	registry := h.Item.(*ecr.DescribeRegistryOutput)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ecr_registry.getEcrRegistryArn", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	if registry == nil || registry.RegistryId == nil {
		return nil, nil
	}

	arn := "arn:" + commonColumnData.Partition + ":ecr:" + region + ":" + *registry.RegistryId + ":registry/" + *registry.RegistryId
	return arn, nil
}
