package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"

	wellarchitectedv1 "github.com/aws/aws-sdk-go/service/wellarchitected"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedShareInvitation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_share_invitation",
		Description: "AWS Well-Architected Share Invitation",
		List: &plugin.ListConfig{
			Hydrate: listWellArchitectedShareInvitations,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_name", Require: plugin.Optional},
				{Name: "lens_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "lens_arn",
				Description: "The ARN for the lens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lens_name",
				Description: "The full name of the lens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "permission_type",
				Description: "Permission granted on a workload share.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "share_invitation_id",
				Description: "The ID assigned to the share invitation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "share_resource_type",
				Description: "The resource type of the share invitation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "shared_by",
				Description: "An Amazon Web Services account ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "shared_with",
				Description: "The Amazon Web Services account ID, IAM role, organization ID, or organizational unit (OU) ID with which the workload is shared.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workload_name",
				Description: "The name of the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workload_id",
				Description: "The ID assigned to the workload.",
				Type:        proto.ColumnType_STRING,
			},
			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ShareInvitationId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedShareInvitationAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listWellArchitectedShareInvitations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create client
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_share_invitation.listWellArchitectedShareInvitations", "client_error", err)
		return nil, err
	}

	// Unsupported region, return no data
	if svc == nil {
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &wellarchitected.ListShareInvitationsInput{
		MaxResults: maxLimit,
	}

	if d.EqualsQualString("workload_name") != "" {
		if d.EqualsQualString("lens_name") != "" {
			return nil, nil
		}
		input.WorkloadNamePrefix = aws.String(d.EqualsQualString("workload_name"))
		input.ShareResourceType = types.ShareResourceTypeWorkload
	}
	if d.EqualsQualString("lens_name") != "" {
		if d.EqualsQualString("workload_name") != "" {
			return nil, nil
		}
		input.LensNamePrefix = aws.String(d.EqualsQualString("lens_name"))
		input.ShareResourceType = types.ShareResourceTypeLens
	}

	paginator := wellarchitected.NewListShareInvitationsPaginator(svc, input, func(o *wellarchitected.ListShareInvitationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_share_invitation.listWellArchitectedShareInvitations", "api_error", err)
			return nil, err
		}

		for _, item := range output.ShareInvitationSummaries {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getWellArchitectedShareInvitationAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	invitation := h.Item.(types.ShareInvitationSummary)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":wellarchitected:" + region + ":" + commonColumnData.AccountId + ":share-invitation/" + *invitation.ShareInvitationId}

	return akas, nil
}
