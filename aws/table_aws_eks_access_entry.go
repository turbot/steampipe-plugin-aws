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

func tableAwsEksAccessEntry(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_access_entry",
		Description: "AWS EKS Access Entry",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"cluster_name", "principal_arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameterException", "InvalidParameter"}),
			},
			Hydrate: getEksAccessEntry,
			Tags:    map[string]string{"service": "eks", "action": "DescribeAccessEntry"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEKSClusters,
			Hydrate:       listEKSAccessEntries,
			Tags:          map[string]string{"service": "eks", "action": "ListAccessEntries"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getEksAccessEntry,
				Tags: map[string]string{"service": "eks", "action": "DescribeAccessEntry"},
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
				Name:        "access_entry_arn",
				Description: "The ARN of the access entry.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksAccessEntry,
			},
			{
				Name:        "created_at",
				Description: "The date and time that the access entry was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksAccessEntry,
			},
			{
				Name:        "modified_at",
				Description: "The date and time that the access entry was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksAccessEntry,
			},
			{
				Name:        "kubernetes_groups",
				Description: "A name that you've specified in a Kubernetes RoleBinding or ClusterRoleBinding object so that Kubernetes authorizes the principalARN access to cluster objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksAccessEntry,
			},
			{
				Name:        "type",
				Description: "The type of the access entry.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksAccessEntry,
			},
			{
				Name:        "username",
				Description: "The username to authenticate to Kubernetes with.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksAccessEntry,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksAccessEntry,
				Transform:   transform.FromField("PrincipalArn"),
			},
			{
				Name:        "tags",
				Description: "The metadata that you apply to the access entry to assist with categorization and organization.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksAccessEntry,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksAccessEntry,
				Transform:   transform.FromField("AccessEntryArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listEKSAccessEntries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get cluster details
	clusterName := *h.Item.(types.Cluster).Name

	// Create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_access_entry.listEKSAccessEntries", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &eks.ListAccessEntriesInput{
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

	paginator := eks.NewListAccessEntriesPaginator(svc, input, func(o *eks.ListAccessEntriesPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_access_entry.listEKSAccessEntries", "api_error", err)
			return nil, err
		}

		for _, item := range output.AccessEntries {
			d.StreamListItem(ctx, &AccessEntryInfo{
				ClusterName:  &clusterName,
				PrincipalArn: aws.String(item),
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

func getEksAccessEntry(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var clusterName, principalArn string
	if h.Item != nil {
		clusterName = *h.Item.(*AccessEntryInfo).ClusterName
		principalArn = *h.Item.(*AccessEntryInfo).PrincipalArn
	} else {
		clusterName = d.EqualsQuals["cluster_name"].GetStringValue()
		principalArn = d.EqualsQuals["principal_arn"].GetStringValue()
	}

	// create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_access_entry.getEksAccessEntry", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &eks.DescribeAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &principalArn,
	}

	op, err := svc.DescribeAccessEntry(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_access_entry.getEksAccessEntry", "api_error", err)
		return nil, err
	}

	return op.AccessEntry, nil
}

// AccessEntryInfo is a struct to hold cluster name and principal ARN for list operations
type AccessEntryInfo struct {
	ClusterName  *string
	PrincipalArn *string
}
