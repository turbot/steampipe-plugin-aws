package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/shield"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsShieldProtection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_protection",
		Description: "Aws Shield Protection",
		Get:         &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate: getShieldProtection,
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate: listShieldProtections,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the protection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier (ID) of the protection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the protection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProtectionArn"),
			},
			{
				Name:        "resource_arn",
				Description: "The ARN (Amazon Resource Name) of the AWS resource that is protected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_check_ids",
				Description: "The unique identifier (ID) for the Route 53 health check that's associated with the protection.",
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
				Transform:   transform.FromField("ProtectionArn").Transform(transform.EnsureStringArray),
			},

		}),
	}
}

//// LIST FUNCTION

func listShieldProtections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listShieldProtections", "AWS_REGION", region)

	// Create Session
	svc, err := ShieldService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListProtectionsPages(
		&shield.ListProtectionsInput{},
		func(page *shield.ListProtectionsOutput, isLast bool) bool {
			for _, protection := range page.Protections {
				d.StreamListItem(ctx, protection)
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listShieldProtections", "query_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getShieldProtection(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	protectionId := d.KeyColumnQuals["id"].GetStringValue()

	// Create service
	svc, err := ShieldService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &shield.DescribeProtectionInput{
		ProtectionId: &protectionId,
	}

	op, err := svc.DescribeProtection(params)
	if err != nil {
		return nil, err
	}

	if op.Protection != nil {
		return op.Protection, nil
	}

	return nil, nil
}
