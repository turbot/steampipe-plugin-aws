package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityHub(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_hub",
		Description: "AWS Security hub",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("hub_arn"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidVolume.NotFound", "InvalidParameterValue"}),
			Hydrate:           describeHub,
		},
		List: &plugin.ListConfig{
			Hydrate: listHub,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "hub_arn",
				Description: "The ARN of the Hub resource that was retrieved.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_enable_controls",
				Description: "Whether to automatically enable new controls when they are added to standards that are enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "subscribed_at",
				Description: "The date and time when Security Hub was enabled in the account.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			/// Standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HubArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listHub(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listEBSVolume", "AWS_REGION", region)

	// Create session
	svc, err := SecurityHubService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	resp, err := svc.DescribeHub(&securityhub.DescribeHubInput{})

	d.StreamListItem(ctx, resp)

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeHub(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("describeHub")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	hubArn := d.KeyColumnQuals["hub_arn"].GetStringValue()

	// get service
	svc, err := SecurityHubService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &securityhub.DescribeHubInput{
		HubArn: &hubArn,
	}

	// Get call
	op, err := svc.DescribeHub(params)
	if err != nil {
		plugin.Logger(ctx).Debug("describeHub", "ERROR", err)
		return nil, err
	}
	return op, nil
}
