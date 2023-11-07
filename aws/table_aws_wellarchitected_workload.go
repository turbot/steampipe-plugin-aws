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

func tableAwsWellArchitectedWorkload(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_workload",
		Description: "AWS Well-Architected Workload",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("workload_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getWellArchitectedWorkload,
			Tags:    map[string]string{"service": "wellarchitected", "action": "GetWorkload"},
		},
		List: &plugin.ListConfig{
			Hydrate: listWellArchitectedWorkloads,
			Tags:    map[string]string{"service": "wellarchitected", "action": "ListWorkloads"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_name", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getWellArchitectedWorkload,
				Tags: map[string]string{"service": "wellarchitected", "action": "GetWorkload"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
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
				Hydrate:     getWellArchitectedWorkload,
			},
			{
				Name:        "description",
				Description: "The description for the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorkload,
			},
			{
				Name:        "environment",
				Description: "The environment for the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorkload,
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
				Hydrate:     getWellArchitectedWorkload,
			},
			{
				Name:        "industry_type",
				Description: "The industry type for the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorkload,
			},
			{
				Name:        "is_review_owner_update_acknowledged",
				Description: "Flag indicating whether the workload owner has acknowledged that the review owner field is required.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getWellArchitectedWorkload,
			},
			{
				Name:        "notes",
				Description: "The notes associated with the workload.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWellArchitectedWorkload,
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
				Hydrate:     getWellArchitectedWorkload,
			},
			{
				Name:        "review_restriction_date",
				Description: "The date and time recorded.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getWellArchitectedWorkload,
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
				Hydrate:     getWellArchitectedWorkload,
			},
			{
				Name:        "aws_regions",
				Description: "The list of AWS Regions associated with the workload, for example, us-east-2, or ca-central-1.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedWorkload,
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
				Hydrate:     getWellArchitectedWorkload,
			},
			{
				Name:        "pillar_priorities",
				Description: "The priorities of the pillars, which are used to order items in the improvement plan. ",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWellArchitectedWorkload,
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
				Hydrate:     getWellArchitectedWorkload,
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

func listWellArchitectedWorkloads(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_workload.listWellArchitectedWorkloads", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &wellarchitected.ListWorkloadsInput{
		MaxResults: maxLimit,
	}

	equalQuals := d.EqualsQuals
	if equalQuals["workload_name"] != nil {
		if equalQuals["workload_name"].GetStringValue() != "" {
			input.WorkloadNamePrefix = aws.String(equalQuals["workload_name"].GetStringValue())
		}
	}

	paginator := wellarchitected.NewListWorkloadsPaginator(svc, input, func(o *wellarchitected.ListWorkloadsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_wellarchitected_workload.listWellArchitectedWorkloads", "api_error", err)
			return nil, err
		}

		for _, items := range output.WorkloadSummaries {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWellArchitectedWorkload(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		id = workloadID(h.Item)
	} else {
		quals := d.EqualsQuals
		id = quals["workload_id"].GetStringValue()
	}

	// Create Session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_workload.getWellArchitectedWorkload", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &wellarchitected.GetWorkloadInput{
		WorkloadId: aws.String(id),
	}

	op, err := svc.GetWorkload(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_workload.getWellArchitectedWorkload", "api_error", err)
		return nil, err
	}

	return op.Workload, nil
}

//// TRANSFORM FUNCTIONS

func workloadID(item interface{}) string {
	switch item := item.(type) {
	case types.WorkloadSummary:
		return *item.WorkloadId
	case *types.Workload:
		return *item.WorkloadId
	}
	return ""
}
