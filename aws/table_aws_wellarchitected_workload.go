package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/wellarchitected"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedWorkload(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_workload",
		Description: "AWS Well-Architected Workload",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("workload_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getWellArchitectedWorkload,
		},
		List: &plugin.ListConfig{
			Hydrate: listWellArchitectedWorkloads,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "workload_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	logger := plugin.Logger(ctx)
	logger.Trace("listWellArchitectedWorkloads")

	// Create session
	svc, err := WellArchitectedService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &wellarchitected.ListWorkloadsInput{
		MaxResults: aws.Int64(50),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["workload_name"] != nil {
		if equalQuals["workload_name"].GetStringValue() != "" {
			input.WorkloadNamePrefix = aws.String(equalQuals["workload_name"].GetStringValue())
		}
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	err = svc.ListWorkloadsPages(
		input,
		func(page *wellarchitected.ListWorkloadsOutput, lastPage bool) bool {
			for _, Workload := range page.WorkloadSummaries {
				d.StreamListItem(ctx, Workload)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	if err != nil {
		var awsErr awserr.Error
		// AWS Well-Architected Tool is not supported in all regions. For unsupported regions the API throws an error, e.g.,
		// Post "https://wellarchitected.ap-northeast-3.amazonaws.com/workloadsSummaries": dial tcp: lookup wellarchitected.ap-northeast-3.amazonaws.com: no such host
		if errors.As(err, &awsErr) {
			if awsErr.OrigErr() != nil {
				if strings.Contains(awsErr.OrigErr().Error(), "no such host") {
					return nil, nil
				}
			}
		}
		logger.Error("listWellArchitectedWorkloads", "list", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWellArchitectedWorkload(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getWellArchitectedWorkload")

	var id string
	if h.Item != nil {
		id = workloadID(h.Item)
	} else {
		quals := d.KeyColumnQuals
		id = quals["workload_id"].GetStringValue()
	}

	// Create Session
	svc, err := WellArchitectedService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &wellarchitected.GetWorkloadInput{
		WorkloadId: aws.String(id),
	}

	op, err := svc.GetWorkload(params)
	if err != nil {
		var awsErr awserr.Error
		// AWS Well-Architected Tool is not supported in all regions. For unsupported regions the API throws an error, e.g.,
		// Post "https://wellarchitected.ap-northeast-3.amazonaws.com/workloadsSummaries": dial tcp: lookup wellarchitected.ap-northeast-3.amazonaws.com: no such host
		if errors.As(err, &awsErr) {
			if awsErr.OrigErr() != nil {
				if strings.Contains(awsErr.OrigErr().Error(), "no such host") {
					return nil, nil
				}
			}
		}
		logger.Error("getWellArchitectedWorkload", "get", err)
		return nil, err
	}

	return op.Workload, nil
}

//// TRANSFORM FUNCTIONS

func workloadID(item interface{}) string {
	switch item := item.(type) {
	case *wellarchitected.WorkloadSummary:
		return *item.WorkloadId
	case *wellarchitected.Workload:
		return *item.WorkloadId
	}
	return ""
}
