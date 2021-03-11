package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/eks"
)

func tableAwsEksFargateProfile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_fargate_profile",
		Description: "AWS Elastic Kubernetes Service Fargate Profile",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "cluster_name"}),
			Hydrate:    getEksFargateProfile,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEksClusters,
			Hydrate:       ListEksFargateProfile,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Fargate profile.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FargateProfileName"),
			},
			{
				Name:        "fargate_profile_arn",
				Description: "The full Amazon Resource Name (ARN) of the Fargate profile.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksFargateProfile,
			},
			{
				Name:        "cluster_name",
				Description: "The name of the Amazon EKS cluster that the Fargate profile belongs to.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksFargateProfile,
			},
			{
				Name:        "created_at",
				Description: "The Unix epoch timestamp in seconds for when the Fargate profile was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksFargateProfile,
			},
			{
				Name:        "pod_execution_role_arn",
				Description: "The Amazon Resource Name (ARN) of the pod execution role to use for pods that match the selectors in the Fargate profile.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksFargateProfile,
			},
			{
				Name:        "subnets",
				Description: "The IDs of subnets to launch pods into.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksFargateProfile,
			},
			{
				Name:        "selectors",
				Description: "The selectors to match for pods to use this Fargate profile.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksFargateProfile,
			},
			{
				Name:        "status",
				Description: "The current status of the Fargate profile.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksFargateProfile,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FargateProfileName"),
				Hydrate:     getEksFargateProfile,
			},
			{
				Name:        "tags",
				Description: "The metadata applied to the Fargate profile to assist with categorization and organization.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksFargateProfile,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FargateProfileArn").Transform(arnToAkas),
				Hydrate:     getEksFargateProfile,
			},
		}),
	}
}

//// LIST FUNCTION

func ListEksFargateProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("ListEksFargateProfile", "AWS_REGION", region)

	// Create service
	svc, err := EksService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	clusterName := *h.Item.(*eks.Cluster).Name

	resp, err := svc.ListFargateProfiles(&eks.ListFargateProfilesInput{
		ClusterName: &clusterName,
	})
	for _, fargateProfile := range resp.FargateProfileNames {
		d.StreamLeafListItem(ctx, &eks.FargateProfile{
			ClusterName:        &clusterName,
			FargateProfileName: fargateProfile,
		})
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEksFargateProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEksFargateProfile")
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var clusterName string
	var fargateProfileName string
	if h.Item != nil {
		clusterName = *h.Item.(*eks.FargateProfile).ClusterName
		fargateProfileName = *h.Item.(*eks.FargateProfile).FargateProfileName
	} else {
		clusterName = d.KeyColumnQuals["cluster_name"].GetStringValue()
		fargateProfileName = d.KeyColumnQuals["name"].GetStringValue()
	}

	// create service
	svc, err := EksService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &eks.DescribeFargateProfileInput{
		ClusterName:        &clusterName,
		FargateProfileName: &fargateProfileName,
	}

	op, err := svc.DescribeFargateProfile(params)
	if err != nil {
		return nil, err
	}

	return op.FargateProfile, nil
}
