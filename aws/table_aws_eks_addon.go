package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEksAddon(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_addon",
		Description: "AWS EKS Addon",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"addon_name", "cluster_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ResourceNotFoundException", "InvalidParameterException", "InvalidParameter"}),
			},
			Hydrate: getEksAddon,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEksClusters,
			Hydrate:       listEksAddons,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "addon_name",
				Description: "The name of the add-on.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the add-on.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksAddon,
				Transform:   transform.FromField("AddonArn"),
			},
			{
				Name:        "cluster_name",
				Description: "The name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "addon_version",
				Description: "The version of the add-on.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksAddon,
			},
			{
				Name:        "status",
				Description: "The status of the add-on.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksAddon,
			},
			{
				Name:        "created_at",
				Description: "The date and time that the add-on was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksAddon,
			},
			{
				Name:        "modified_at",
				Description: "The date and time that the add-on was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksAddon,
			},
			{
				Name:        "service_account_role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role that is bound to the Kubernetes service account used by the add-on.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksAddon,
			},
			{
				Name:        "health_issues",
				Description: "An object that represents the add-on's health issues.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksAddon,
				Transform:   transform.FromField("Health.Issues"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AddonName"),
			},
			{
				Name:        "tags",
				Description: "The metadata that you apply to the cluster to assist with categorization and organization.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksAddon,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksAddon,
				Transform:   transform.FromField("AddonArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEksAddons(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get cluster details
	clusterName := *h.Item.(types.Cluster).Name

	// Create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_addon.listEksAddons", "get_client_error", err)
		return nil, err
	}

	input := &eks.ListAddonsInput{
		ClusterName: &clusterName,
		MaxResults:  aws.Int32(100),
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxResults {
			if limit < 1 {
				input.MaxResults = aws.Int32(1)
			} else {
				input.MaxResults = aws.Int32(limit)
			}
		}
	}

	paginator := eks.NewListAddonsPaginator(svc, input, func(o *eks.ListAddonsPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_addon.listEksAddons", "api_error", err)
			return nil, err
		}

		for _, item := range output.Addons {
			d.StreamListItem(ctx, types.Addon{
				AddonName:   aws.String(item),
				ClusterName: &clusterName,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEksAddon(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var clusterName, addonName string
	if h.Item != nil {
		clusterName = *h.Item.(types.Addon).ClusterName
		addonName = *h.Item.(types.Addon).AddonName
	} else {
		clusterName = d.KeyColumnQuals["cluster_name"].GetStringValue()
		addonName = d.KeyColumnQuals["addon_name"].GetStringValue()
	}

	// create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_addon.getEksAddon", "get_client_error", err)
		return nil, err
	}

	params := &eks.DescribeAddonInput{
		AddonName:   &addonName,
		ClusterName: &clusterName,
	}

	op, err := svc.DescribeAddon(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_addon.getEksAddon", "api_error", err)
		return nil, err
	}

	return op.Addon, nil
}
