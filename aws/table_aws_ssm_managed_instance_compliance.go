package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSSMManagedInstanceCompliance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_managed_instance_compliance",
		Description: "AWS SSM Managed Instance Compliance",
		List: &plugin.ListConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidResourceId", "ValidationException"}),
			},
			Hydrate: listSsmManagedInstanceCompliances,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "resource_id", Require: plugin.Required},
				{Name: "resource_type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
		MaxResults:  aws.Int64(50),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["resource_type"] != nil {
		if equalQuals["resource_type"].GetStringValue() != "" {
			params.ResourceTypes = []*string{aws.String(equalQuals["resource_type"].GetStringValue())}
		} else {
			params.ResourceTypes = getListValues(equalQuals["resource_type"].GetListValue())
		}
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			if *limit < 1 {
				params.MaxResults = aws.Int64(1)
			} else {
				params.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListComplianceItemsPages(
		params,
		func(page *ssm.ListComplianceItemsOutput, isLast bool) bool {
			for _, item := range page.ComplianceItems {
				d.StreamListItem(ctx, item)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
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
