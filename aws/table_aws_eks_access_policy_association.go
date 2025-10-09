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
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"cluster_name", "principal_arn", "policy_arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameterException", "InvalidParameter"}),
			},
			Hydrate: getEksAccessPolicyAssociation,
			Tags:    map[string]string{"service": "eks", "action": "DescribeAccessPolicyAssociation"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEKSAccessPolicyAssociations,
			Tags:    map[string]string{"service": "eks", "action": "ListAssociatedAccessPolicies"},
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
				Transform:   transform.FromField("AssociatedAccessPolicy.PolicyArn"),
			},
			{
				Name:        "associated_at",
				Description: "The date and time that the access policy was associated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("AssociatedAccessPolicy.AssociatedAt"),
			},
			{
				Name:        "modified_at",
				Description: "The date and time that the access policy association was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("AssociatedAccessPolicy.ModifiedAt"),
			},
			{
				Name:        "access_scope",
				Description: "The scope of the access policy association.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssociatedAccessPolicy.AccessScope"),
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

//// LIST FUNCTION

func listEKSAccessPolicyAssociations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_access_policy_association.listEKSAccessPolicyAssociations", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Get optional filter values
	filterClusterName := d.EqualsQuals["cluster_name"]
	filterPrincipalArn := d.EqualsQuals["principal_arn"]

	// Step 1: List all clusters (or specific cluster if filtered)
	clusterInput := &eks.ListClustersInput{
		MaxResults: aws.Int32(100),
	}

	clusterPaginator := eks.NewListClustersPaginator(svc, clusterInput, func(o *eks.ListClustersPaginatorOptions) {
		o.Limit = *clusterInput.MaxResults
		o.StopOnDuplicateToken = true
	})

	for clusterPaginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		clusterOutput, err := clusterPaginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_access_policy_association.listEKSAccessPolicyAssociations", "list_clusters_error", err)
			return nil, err
		}

		for _, clusterName := range clusterOutput.Clusters {
			// Skip if cluster name filter doesn't match
			if filterClusterName != nil && filterClusterName.GetStringValue() != clusterName {
				continue
			}

			// Step 2: List access entries for this cluster
			accessEntryInput := &eks.ListAccessEntriesInput{
				ClusterName: &clusterName,
				MaxResults:  aws.Int32(100),
			}

			accessEntryPaginator := eks.NewListAccessEntriesPaginator(svc, accessEntryInput, func(o *eks.ListAccessEntriesPaginatorOptions) {
				o.Limit = *accessEntryInput.MaxResults
				o.StopOnDuplicateToken = true
			})

			for accessEntryPaginator.HasMorePages() {
				d.WaitForListRateLimit(ctx)

				accessEntryOutput, err := accessEntryPaginator.NextPage(ctx)
				if err != nil {
					plugin.Logger(ctx).Error("aws_eks_access_policy_association.listEKSAccessPolicyAssociations", "list_access_entries_error", err)
					// Continue with next cluster if this one fails
					break
				}

				for _, principalArn := range accessEntryOutput.AccessEntries {
					// Skip if principal ARN filter doesn't match
					if filterPrincipalArn != nil && filterPrincipalArn.GetStringValue() != principalArn {
						continue
					}

					// Step 3: List associated access policies for this access entry
					policyInput := &eks.ListAssociatedAccessPoliciesInput{
						ClusterName:  &clusterName,
						PrincipalArn: &principalArn,
						MaxResults:   aws.Int32(100),
					}

					policyPaginator := eks.NewListAssociatedAccessPoliciesPaginator(svc, policyInput, func(o *eks.ListAssociatedAccessPoliciesPaginatorOptions) {
						o.Limit = *policyInput.MaxResults
						o.StopOnDuplicateToken = true
					})

					for policyPaginator.HasMorePages() {
						d.WaitForListRateLimit(ctx)

						policyOutput, err := policyPaginator.NextPage(ctx)
						if err != nil {
							plugin.Logger(ctx).Error("aws_eks_access_policy_association.listEKSAccessPolicyAssociations", "list_policies_error", err)
							// Continue with next access entry if this one fails
							break
						}

						for _, policy := range policyOutput.AssociatedAccessPolicies {
							d.StreamListItem(ctx, &AccessPolicyAssociationInfo{
								ClusterName:            &clusterName,
								PrincipalArn:           &principalArn,
								AssociatedAccessPolicy: policy,
							})

							// Context can be cancelled due to manual cancellation or the limit has been hit
							if d.RowsRemaining(ctx) == 0 {
								return nil, nil
							}
						}
					}
				}
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEksAccessPolicyAssociation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var clusterName, principalArn, policyArn string
	if h.Item != nil {
		item := h.Item.(*AccessPolicyAssociationInfo)
		clusterName = *item.ClusterName
		principalArn = *item.PrincipalArn
		policyArn = *item.AssociatedAccessPolicy.PolicyArn
	} else {
		clusterName = d.EqualsQuals["cluster_name"].GetStringValue()
		principalArn = d.EqualsQuals["principal_arn"].GetStringValue()
		policyArn = d.EqualsQuals["policy_arn"].GetStringValue()
	}

	// create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_access_policy_association.getEksAccessPolicyAssociation", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Use ListAssociatedAccessPolicies and filter for the specific policy
	// since there's no direct DescribeAccessPolicyAssociation API
	input := &eks.ListAssociatedAccessPoliciesInput{
		ClusterName:  &clusterName,
		PrincipalArn: &principalArn,
		MaxResults:   aws.Int32(100),
	}

	paginator := eks.NewListAssociatedAccessPoliciesPaginator(svc, input, func(o *eks.ListAssociatedAccessPoliciesPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_access_policy_association.getEksAccessPolicyAssociation", "api_error", err)
			return nil, err
		}

		for _, item := range output.AssociatedAccessPolicies {
			if *item.PolicyArn == policyArn {
				return &AccessPolicyAssociationInfo{
					ClusterName:            &clusterName,
					PrincipalArn:           &principalArn,
					AssociatedAccessPolicy: item,
				}, nil
			}
		}
	}

	return nil, nil
}

// AccessPolicyAssociationInfo is a struct to hold cluster name, principal ARN, and policy details
type AccessPolicyAssociationInfo struct {
	ClusterName            *string
	PrincipalArn           *string
	AssociatedAccessPolicy types.AssociatedAccessPolicy
}
