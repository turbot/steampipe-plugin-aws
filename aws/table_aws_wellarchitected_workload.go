package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/wellarchitected"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedWorload(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_workload",
		Description: "AWS Well Architected Workload",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("workload_id"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getWellArchitectedWorload,
		},
		List: &plugin.ListConfig{
			Hydrate: listWellArchitectedWorloads,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "workload_name",
				Description: "The name of the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workload_arn",
				Description: "The ARN for the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workload_id",
				Description: "The ID assigned to the workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "architectural_design",
				Description: "The URL of the architectural design for the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "description",
				Description: "The description for the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "environment",
				Description: "The environment for the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "improvement_status",
				Description: "The improvement status for a workload.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "industry",
				Description: "The industry for the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "industry_type",
				Description: "The industry type for the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "is_review_owner_update_acknowledged",
				Description: "Flag indicating whether the workload owner has acknowledged that the review owner field is required.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "notes",
				Description: "The notes associated with the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "owner",
				Description: "An AWS account ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "review_owner",
				Description: "The review owner of the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "review_restriction_date",
				Description: "The date and time recorded.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "share_invitation_id",
				Description: "The ID assigned to the share invitation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The date and time recorded.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "account_ids",
				Description: "The list of AWS account IDs associated with the workload.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "aws_regions",
				Description: "The list of AWS Regions associated with the workload, for example, us-east-2, or ca-central-1.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "lenses",
				Description: "The list of lenses associated with the workload. Each lens is identified by its LensSummary$LensAlias.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "non_aws_regions",
				Description: "The list of non-AWS Regions associated with the workload.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "pillar_priorities",
				Description: "The priorities of the pillars, which are used to order items in the improvement plan. ",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "risk_counts",
				Description: "A map from risk names to the count of how questions have that rating.",
				Type:        proto.ColumnType_JSON,
			},
			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkloadName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedWorload,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkloadArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listWellArchitectedWorloads(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listWellArchitectedWorloads", "AWS_REGION", region)

	// Create session
	svc, err := WellArchitectedService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	err = svc.ListWorkloadsPages(
		&wellarchitected.ListWorkloadsInput{},
		func(page *wellarchitected.ListWorkloadsOutput, lastPage bool) bool {
			for _, Workload := range page.WorkloadSummaries {
				d.StreamListItem(ctx, Workload)
			}
			return !lastPage
		},
	)

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWellArchitectedWorload(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getWellArchitectedWorload")

	var id string
	if h.Item != nil {
		id = workloadID(h.Item)
	} else {
		quals := d.KeyColumnQuals
		id = quals["workload_id"].GetStringValue()
	}

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := WellArchitectedService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &wellarchitected.GetWorkloadInput{
		WorkloadId: aws.String(id),
	}

	op, err := svc.GetWorkload(params)
	if err != nil {
		logger.Debug("getWellArchitectedWorload", "ERROR", err)
		return nil, err
	}

	return op.Workload, nil
}

//// TRANSFORM FUNCTIONS

func workloadID(item interface{}) string {
	switch item.(type) {
	case *wellarchitected.WorkloadSummary:
		return *item.(*wellarchitected.WorkloadSummary).WorkloadId
	case *wellarchitected.Workload:
		return *item.(*wellarchitected.Workload).WorkloadId
	}
	return ""
}
