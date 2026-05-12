package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEksPodIdentityAssociation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_pod_identity_association",
		Description: "AWS EKS Pod Identity Association",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"cluster_name", "association_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameterException", "InvalidParameter"}),
			},
			Hydrate: getEksPodIdentityAssociation,
			Tags:    map[string]string{"service": "eks", "action": "DescribePodIdentityAssociation"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEKSClusters,
			Hydrate:       listEKSPodIdentityAssociations,
			Tags:          map[string]string{"service": "eks", "action": "ListPodIdentityAssociations"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cluster_name", Require: plugin.Optional},
				{Name: "namespace", Require: plugin.Optional},
				{Name: "service_account", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getEksPodIdentityAssociation,
				Tags: map[string]string{"service": "eks", "action": "DescribePodIdentityAssociation"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EKS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cluster_name",
				Description: "The name of the cluster that the association is in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_id",
				Description: "The ID of the association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_arn",
				Description: "The Amazon Resource Name (ARN) of the association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "The name of the Kubernetes namespace inside the cluster that the association is in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_account",
				Description: "The name of the Kubernetes service account that the association uses.",
				Type:        proto.ColumnType_STRING,
			},
			// Fields that require DescribePodIdentityAssociation
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role associated with the service account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksPodIdentityAssociation,
			},
			{
				Name:        "created_at",
				Description: "The timestamp that the association was created at.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksPodIdentityAssociation,
			},
			{
				Name:        "modified_at",
				Description: "The most recent timestamp that the association was modified at.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksPodIdentityAssociation,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssociationId"),
			},
			{
				Name:        "tags",
				Description: "The metadata that you apply to the association to assist with categorization and organization.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksPodIdentityAssociation,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssociationArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

// PodIdentityAssociationInfo holds the list-level summary fields
type PodIdentityAssociationInfo struct {
	ClusterName    *string
	AssociationId  *string
	AssociationArn *string
	Namespace      *string
	ServiceAccount *string
}

//// LIST FUNCTION

func listEKSPodIdentityAssociations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	clusterName := *h.Item.(types.Cluster).Name

	// Apply optional cluster_name filter — skip this cluster if it doesn't match
	filterClusterName := d.EqualsQuals["cluster_name"].GetStringValue()
	if filterClusterName != "" && filterClusterName != clusterName {
		return nil, nil
	}

	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_pod_identity_association.listEKSPodIdentityAssociations", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &eks.ListPodIdentityAssociationsInput{
		ClusterName: &clusterName,
		MaxResults:  aws.Int32(100),
	}

	// Apply optional namespace filter
	if ns := d.EqualsQuals["namespace"].GetStringValue(); ns != "" {
		input.Namespace = aws.String(ns)
	}

	// Apply optional service_account filter
	if sa := d.EqualsQuals["service_account"].GetStringValue(); sa != "" {
		input.ServiceAccount = aws.String(sa)
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

	paginator := eks.NewListPodIdentityAssociationsPaginator(svc, input, func(o *eks.ListPodIdentityAssociationsPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_pod_identity_association.listEKSPodIdentityAssociations", "api_error", err)
			return nil, err
		}

		for _, item := range output.Associations {
			d.StreamListItem(ctx, &PodIdentityAssociationInfo{
				ClusterName:    item.ClusterName,
				AssociationId:  item.AssociationId,
				AssociationArn: item.AssociationArn,
				Namespace:      item.Namespace,
				ServiceAccount: item.ServiceAccount,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEksPodIdentityAssociation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var clusterName, associationId string

	if h.Item != nil {
		info := h.Item.(*PodIdentityAssociationInfo)
		clusterName = aws.ToString(info.ClusterName)
		associationId = aws.ToString(info.AssociationId)
	} else {
		clusterName = d.EqualsQuals["cluster_name"].GetStringValue()
		associationId = d.EqualsQuals["association_id"].GetStringValue()
	}

	// check for empty parameters
	if clusterName == "" || associationId == "" {
		return nil, nil
	}

	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_pod_identity_association.getEksPodIdentityAssociation", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &eks.DescribePodIdentityAssociationInput{
		ClusterName:   aws.String(clusterName),
		AssociationId: aws.String(associationId),
	}

	output, err := svc.DescribePodIdentityAssociation(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_pod_identity_association.getEksPodIdentityAssociation", "api_error", err)
		return nil, err
	}

	if output.Association == nil {
		return nil, nil
	}

	return output.Association, nil
}
