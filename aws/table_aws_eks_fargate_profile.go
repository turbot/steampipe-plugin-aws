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

//// TABLE DEFINITION

func tableAwsEksFargateProfile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_fargate_profile",
		Description: "AWS Elastic Kubernetes Service Fargate Profile",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"cluster_name", "fargate_profile_name"}),
			Hydrate:    getEKSFargateProfile,
			Tags:       map[string]string{"service": "eks", "action": "DescribeFargateProfile"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEKSClusters,
			Hydrate:       listEKSFargateProfiles,
			Tags:          map[string]string{"service": "eks", "action": "ListFargateProfiles"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "cluster_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(eksv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "fargate_profile_name",
				Description: "The name of the Fargate profile.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_name",
				Description: "The name of the Amazon EKS cluster that the Fargate profile belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fargate_profile_arn",
				Description: "The full Amazon Resource Name (ARN) of the Fargate profile.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSFargateProfile,
			},
			{
				Name:        "created_at",
				Description: "The Unix epoch timestamp in seconds for when the Fargate profile was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEKSFargateProfile,
			},
			{
				Name:        "pod_execution_role_arn",
				Description: "The Amazon Resource Name (ARN) of the pod execution role to use for pods that match the selectors in the Fargate profile.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSFargateProfile,
			},
			{
				Name:        "status",
				Description: "The current status of the Fargate profile.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSFargateProfile,
			},
			{
				Name:        "selectors",
				Description: "The selectors to match for pods to use this Fargate profile.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSFargateProfile,
			},
			{
				Name:        "subnets",
				Description: "The subnets used by the Fargate profile.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSFargateProfile,
			},
			{
				Name:        "tags",
				Description: "A list of tags assigned to the Fargate profile.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSFargateProfile,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FargateProfileName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FargateProfileArn").Transform(transform.EnsureStringArray),
				Hydrate:     getEKSFargateProfile,
			},
		}),
	}
}

//// LIST FUNCTION

func listEKSFargateProfiles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(types.Cluster)
	clusterName := cluster.Name

	if d.EqualsQuals["cluster_name"] != nil {
		if *clusterName != d.EqualsQualString("cluster_name") {
			return nil, nil
		}
	}

	if clusterName == nil {
		return nil, nil
	}

	// Create client
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_fargate_profile.listEKSFargateProfiles", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &eks.ListFargateProfilesInput{
		ClusterName: clusterName,
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

	paginator := eks.NewListFargateProfilesPaginator(svc, input, func(o *eks.ListFargateProfilesPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_fargate_profile.listEKSFargateProfiles", "api_error", err)
			return nil, err
		}

		for _, profile := range output.FargateProfileNames {
			d.StreamListItem(ctx, types.FargateProfile{
				ClusterName:        clusterName,
				FargateProfileName: aws.String(profile),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEKSFargateProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var clusterName, fargateProfileName string
	if h.Item != nil {
		clusterName = *h.Item.(types.FargateProfile).ClusterName
		fargateProfileName = *h.Item.(types.FargateProfile).FargateProfileName
	} else {
		clusterName = d.EqualsQuals["cluster_name"].GetStringValue()
		fargateProfileName = d.EqualsQuals["fargate_profile_name"].GetStringValue()
	}

	// create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_fargate_profile.listEKSFargateProfiles", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	if clusterName == "" || fargateProfileName == "" {
		return nil, nil
	}

	params := &eks.DescribeFargateProfileInput{
		ClusterName:        &clusterName,
		FargateProfileName: aws.String(fargateProfileName),
	}

	op, err := svc.DescribeFargateProfile(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_fargate_profile.listEKSFargateProfiles", "api_error", err)
		return nil, err
	}

	return op.FargateProfile, nil
}
