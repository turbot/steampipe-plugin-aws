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

func tableAwsEksAccessPolicyAssociation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_access_policy_association",
		Description: "AWS EKS Access Policy Association",
		List: &plugin.ListConfig{
			ParentHydrate: listEKSClusters,
			Hydrate:       listEKSAccessPolicyAssociations,
			Tags:          map[string]string{"service": "eks", "action": "ListAssociatedAccessPolicies"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cluster_name", Require: plugin.Optional},
				{Name: "principal_arn", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EKS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cluster_name",
				Description: "The name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "principal_arn",
				Description: "The ARN of the IAM principal for the AccessEntry.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policy_arn",
				Description: "The ARN of the access policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associated_at",
				Description: "The date and time that the access policy was associated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "modified_at",
				Description: "The date and time that the access policy association was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "access_scope",
				Description: "The scope of the access policy association.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "access_scope_type",
				Description: "The type of the access scope (cluster or namespace).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssociatedAccessPolicy.AccessScope.Type"),
			},
			{
				Name:        "access_scope_namespaces",
				Description: "The Kubernetes namespaces included in the access scope.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssociatedAccessPolicy.AccessScope.Namespaces"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssociatedAccessPolicy.PolicyArn"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssociatedAccessPolicy.PolicyArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type AssociatedAccessPolicyInfo struct {
	ClusterName  *string
	PrincipalArn *string
	types.AssociatedAccessPolicy
}

//// LIST FUNCTION

func listEKSAccessPolicyAssociations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	clusterName := *h.Item.(types.Cluster).Name
	// Create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_access_policy_association.listEKSAccessPolicyAssociations", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	accessEntries, err := listAccessEntriesByClusterName(ctx, svc, clusterName)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_access_policy_association.listAccessEntriesByClusterName", "api_error", err)
		return nil, err
	}
	// Get optional filter values
	filterClusterName := d.EqualsQuals["cluster_name"].GetStringValue()
	filterPrincipalArn := d.EqualsQuals["principal_arn"].GetStringValue()

	if filterClusterName != "" && filterClusterName != clusterName {
		return nil, nil
	}

	// The accessEntry corresponds to the PrincipalArn.
	for _, accessEntry := range accessEntries {

		if filterPrincipalArn != "" && filterPrincipalArn != accessEntry {
			return nil, nil
		}

		input := &eks.ListAssociatedAccessPoliciesInput{
			ClusterName:  aws.String(clusterName),
			PrincipalArn: aws.String(accessEntry),
			MaxResults:   aws.Int32(100),
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

		paginator := eks.NewListAssociatedAccessPoliciesPaginator(svc, input, func(o *eks.ListAssociatedAccessPoliciesPaginatorOptions) {
			o.Limit = *input.MaxResults
			o.StopOnDuplicateToken = true
		})

		for paginator.HasMorePages() {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			output, err := paginator.NextPage(ctx)
			if err != nil {
				plugin.Logger(ctx).Error("aws_eks_access_policy_association.listEKSAccessPolicyAssociations", "api_error", err)
				return nil, err
			}

			for _, item := range output.AssociatedAccessPolicies {

				d.StreamListItem(ctx, &AssociatedAccessPolicyInfo{output.ClusterName, output.PrincipalArn, item})

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

func listAccessEntriesByClusterName(ctx context.Context, svc *eks.Client, clusterName string) ([]string, error) {

	accessEntries := make([]string, 0)
	input := &eks.ListAccessEntriesInput{
		ClusterName: &clusterName,
		MaxResults:  aws.Int32(100),
	}

	paginator := eks.NewListAccessEntriesPaginator(svc, input, func(o *eks.ListAccessEntriesPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_access_policy_association.listAccessEntriesByClusterName", "api_error", err)
			return nil, err
		}

		accessEntries = append(accessEntries, output.AccessEntries...)
	}

	return accessEntries, nil

}