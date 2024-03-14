package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"

	eksv1 "github.com/aws/aws-sdk-go/service/eks"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type InsightConfig struct {
	ClusterName *string
	types.Insight
}

//// TABLE DEFINITION

func tableAwsEksInsight(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_insight",
		Description: "AWS EKS Insight",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "cluster_name"}),
			Hydrate:    getEKSInsight,
			Tags:       map[string]string{"service": "eks", "action": "DescribeInsight"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameterException", "InvalidParameter"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEKSClusters,
			Hydrate:       listEKSInsights,
			Tags:          map[string]string{"service": "eks", "action": "ListInsights"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cluster_name", Require: plugin.Optional},
			},
		},

		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getEKSInsight,
				Tags: map[string]string{"service": "eks", "action": "DescribeInsight"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(eksv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name associated with an Amazon EKS insight.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "id",
				Description: "The ID of the insight.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "cluster_name",
				Description: "The name of the cluster that the Insight is for.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "category",
				Description: "The category of the insight.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "kubernetes_version",
				Description: "The Kubernetes minor version associated with an insight if applicable.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "last_refresh_time",
				Description: "The time Amazon EKS last successfully completed a refresh of this insight check on the cluster.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "last_transition_time",
				Description: "The time the status of the insight last changed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "description",
				Description: "The description of the insight which includes alert criteria, remediation recommendation, and additional resources (contains Markdown).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "insight_status",
				Description: "An object containing more detail on the status of the insight.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "recommendation",
				Description: "A summary of how to remediate the finding of this insight if applicable.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "additional_info",
				Description: "Links to sources that provide additional context on the insight.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "resources",
				Description: "The details about each resource listed in the insight check result.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "category_specific_summary",
				Description: "Summary information that relates to the category of the insight. Currently only returned with certain insights having category UPGRADE_READINESS .",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSInsight,
			},
			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
				Hydrate:     getEKSInsight,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Id").Transform(transform.EnsureStringArray),
				Hydrate:     getEKSInsight,
			},
		}),
	}
}

//// LIST FUNCTION

func listEKSInsights(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get cluster details
	cluster := d.EqualsQuals["cluster_name"].GetStringValue()
	clusterName := *h.Item.(types.Cluster).Name
	if cluster != "" {
		if cluster != clusterName {
			return nil, nil
		}
	}
	// Create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_insight.listEKSInsights", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(100)
	input := &eks.ListInsightsInput{
		ClusterName: &clusterName,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := eks.NewListInsightsPaginator(svc, input, func(o *eks.ListInsightsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_insight.listEKSInsights", "api_error", err)
			return nil, err
		}

		for _, insightSummary := range output.Insights {
			d.StreamListItem(ctx, &InsightConfig{&clusterName, types.Insight{
				Id: insightSummary.Id,
			}})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEKSInsight(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var clusterName, insightId string
	if h.Item != nil {
		clusterName = *h.Item.(*InsightConfig).ClusterName
		insightId = *h.Item.(*InsightConfig).Id
	} else {
		clusterName = d.EqualsQuals["cluster_name"].GetStringValue()
		insightId = d.EqualsQuals["id"].GetStringValue()
	}

	// check if clusterName or insightId is empty
	if clusterName == "" || insightId == "" {
		return nil, nil
	}

	// create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_insight.getEKSInsight", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &eks.DescribeInsightInput{
		ClusterName: &clusterName,
		Id:          &insightId,
	}

	op, err := svc.DescribeInsight(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_insight.getEKSInsight", "api_error", err)
		return nil, err
	}
	return op.Insight, nil
}
