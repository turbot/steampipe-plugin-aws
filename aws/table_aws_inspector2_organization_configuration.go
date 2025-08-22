package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInspector2OrganizationConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector2_organization_configuration",
		Description: "AWS Inspector2 Organization Configuration",
		List: &plugin.ListConfig{
			Hydrate: listInspector2OrganizationConfiguration,
			Tags:    map[string]string{"service": "inspector2", "action": "DescribeOrganizationConfiguration"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_INSPECTOR2_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "ec2_auto_enable",
				Description: "Represents whether Amazon EC2 scans are automatically enabled for new members of your Amazon Inspector organization.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AutoEnable.Ec2"),
			},
			{
				Name:        "ecr_auto_enable",
				Description: "Represents whether Amazon ECR scans are automatically enabled for new members of your Amazon Inspector organization.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AutoEnable.Ecr"),
			},
			{
				Name:        "lambda_auto_enable",
				Description: "Represents whether Amazon Web Services Lambda standard scans are automatically enabled for new members of your Amazon Inspector organization.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AutoEnable.Lambda"),
			},
			{
				Name:        "lambda_code_auto_enable",
				Description: "Represents whether Lambda code scans are automatically enabled for new members of your Amazon Inspector organization.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AutoEnable.LambdaCode"),
			},
			{
				Name:        "max_account_limit_reached",
				Description: "Represents whether your organization has reached the maximum Amazon Web Services account limit for Amazon Inspector.",
				Type:        proto.ColumnType_BOOL,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getInspector2OrganizationConfigurationTitle),
			},
		}),
	}
}

//// LIST FUNCTION

func listInspector2OrganizationConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := Inspector2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector2_organization_configuration.listInspector2OrganizationConfiguration", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &inspector2.DescribeOrganizationConfigurationInput{}

	// Get organization configuration
	data, err := svc.DescribeOrganizationConfiguration(ctx, params)
	if err != nil {
		// For the regions where we have not enable it we will receive the rror: aws: operation error Inspector2: DescribeOrganizationConfiguration, https response error StatusCode: 403, RequestID: 8bb92eba-7e8d-4186-a83a-640e98b5621f, AccessDeniedException: Invoking account does not have access to describe the organization configuration.
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("Invoking account does not have access to describe the organization configuration")) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_inspector2_organization_configuration.listInspector2OrganizationConfiguration", "api_error", err)
		return nil, err
	}

	// Stream the single organization configuration
	d.StreamListItem(ctx, data)

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getInspector2OrganizationConfigurationTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	region := d.MatrixItem[matrixKeyRegion]

	title := region.(string) + " Inspector2 Organization Configuration"
	return title, nil
}
