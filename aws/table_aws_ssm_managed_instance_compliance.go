package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSSMManagedInstanceCompliance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_managed_instance_compliance",
		Description: "AWS SSM Managed Instance Compliance",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("resource_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidResourceId", "ValidationException"}),
			Hydrate:           listSsmManagedInstanceCompliances,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "An ID for the compliance item.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "A title for the compliance item.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Title"),
			},
			{
				Name:        "resource_id",
				Description: "An ID for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the compliance item.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compliance_type",
				Description: "The compliance type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "severity",
				Description: "The severity of the compliance status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "details",
				Description: "A key-value combination details for the compliance item.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "execution_summary",
				Description: "A summary for the compliance item.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(ssmManagedInstanceComplianceTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSSMManagedInstanceComplianceAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listSsmManagedInstanceCompliances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSsmManagedInstanceCompliances")

	// Create session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	instanceId := d.KeyColumnQuals["resource_id"].GetStringValue()

	// Build the params
	params := &ssm.ListComplianceItemsInput{
		ResourceIds: []*string{aws.String(instanceId)},
	}

	// List call
	err = svc.ListComplianceItemsPages(
		params,
		func(page *ssm.ListComplianceItemsOutput, isLast bool) bool {
			for _, item := range page.ComplianceItems {
				d.StreamListItem(ctx, item)
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Trace("listSsmManagedInstanceCompliances", "ListComplianceItemsPages_error", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSSMManagedInstanceComplianceAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSSMInstanceComplianceAkas")
	data := h.Item.(*ssm.ComplianceItem)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":managed-instance/" + *data.ResourceId + "/compliance-item/" + *data.Id + ":" + *data.ComplianceType}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func ssmManagedInstanceComplianceTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ssm.ComplianceItem)

	title := *data.Id
	if len(*data.Title) > 0 {
		title = *data.Title
	}

	return title, nil
}
