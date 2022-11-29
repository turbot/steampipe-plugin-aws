package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
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
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidResourceId", "ValidationException"}),
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
	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_managed_instance_compliance.listSsmManagedInstanceCompliances", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	instanceId := d.KeyColumnQuals["resource_id"].GetStringValue()

	// Build the params
	maxItems := int32(50)
	input := &ssm.ListComplianceItemsInput{
		ResourceIds: []string{instanceId},
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["resource_type"] != nil {
		input.ResourceTypes = []string{equalQuals["resource_type"].GetStringValue()}
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 5 {
				maxItems = int32(5)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := ssm.NewListComplianceItemsPaginator(svc, input, func(o *ssm.ListComplianceItemsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_managed_instance_compliance.listSsmManagedInstanceCompliances", "api_error", err)
			return nil, err
		}

		for _, item := range output.ComplianceItems {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSSMManagedInstanceComplianceAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(types.ComplianceItem)
	region := d.KeyColumnQualString(matrixKeyRegion)

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_managed_instance_compliance.getSSMInstanceComplianceAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":managed-instance/" + *data.ResourceId + "/compliance-item/" + *data.Id + ":" + *data.ComplianceType}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func ssmManagedInstanceComplianceTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.ComplianceItem)

	title := *data.Id
	if len(*data.Title) > 0 {
		title = *data.Title
	}

	return title, nil
}
