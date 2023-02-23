package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/emr"

	emrv1 "github.com/aws/aws-sdk-go/service/emr"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEmrBlockPublicAccessConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_emr_block_public_access_configuration",
		Description: "AWS EMR Block Public Access Configuration",
		List: &plugin.ListConfig{
			Hydrate: listBlockPublicAccessConfigurations,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(emrv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "created_by_arn",
				Description: "The Amazon Resource Name that created or last modified the configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BlockPublicAccessConfigurationMetadata.CreatedByArn"),
			},
			{
				Name:        "block_public_security_group_rules",
				Description: "Indicates whether Amazon EMR block public access is enabled (true) or disabled (false).",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("BlockPublicAccessConfiguration.BlockPublicSecurityGroupRules"),
			},
			{
				Name:        "creation_date",
				Description: "The date and time that the configuration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("BlockPublicAccessConfigurationMetadata.CreationDateTime"),
			},
			{
				Name:        "classification",
				Description: "The classification within a configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BlockPublicAccessConfiguration.Classification"),
			},
			{
				Name:        "configurations",
				Description: "A list of additional configurations to apply within a configuration object.",
				Transform:   transform.FromField("BlockPublicAccessConfiguration.Configurations"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "permitted_public_security_group_rule_ranges",
				Description: "Specifies ports and port ranges that are permitted to have security group rules that allow inbound traffic from all public sources.",
				Transform:   transform.FromField("BlockPublicAccessConfiguration.PermittedPublicSecurityGroupRuleRanges"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "properties",
				Description: "A set of properties specified within a configuration classification.",
				Transform:   transform.FromField("BlockPublicAccessConfiguration.Properties"),
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

func listBlockPublicAccessConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := EMRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_cluster_block_public_access.listBlockPublicAccessConfigurations", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &emr.GetBlockPublicAccessConfigurationInput{}

	op, err := svc.GetBlockPublicAccessConfiguration(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_cluster_block_public_access.listBlockPublicAccessConfigurations", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, op)

	return nil, nil
}
