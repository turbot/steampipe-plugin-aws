package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"

	eksv1 "github.com/aws/aws-sdk-go/service/eks"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type IdentityProviderConfig struct {
	Name *string
	Type *string
	types.OidcIdentityProviderConfig
}

//// TABLE DEFINITION

func tableAwsEksIdentityProviderConfig(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_identity_provider_config",
		Description: "AWS EKS Identity Provider Config",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "type", "cluster_name"}),
			Hydrate:    getEKSIdentityProviderConfig,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEKSClusters,
			Hydrate:       listEKSIdentityProviderConfigs,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(eksv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the identity provider configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the identity provider configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_id",
				Description: "This is also known as audience. The ID of the client application that makes authentication requests to the OIDC identity provider.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "cluster_name",
				Description: "The name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IdentityProviderConfigArn"),
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "groups_claim",
				Description: "The JSON web token (JWT) claim that the provider uses to return your groups.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "groups_prefix",
				Description: "The prefix that is prepended to group claims to prevent clashes with existing names (such as system: groups).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "issuer_url",
				Description: "The URL of the OIDC identity provider that allows the API server to discover public signing keys for verifying tokens.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "username_claim",
				Description: "The JSON Web token (JWT) claim that is used as the username.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "status",
				Description: "The status of the OIDC identity provider.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "username_prefix",
				Description: "The prefix that is prepended to username claims to prevent clashes with existing names.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "required_claims",
				Description: "The key-value pairs that describe required claims in the identity token.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "tags_src",
				Description: "The metadata to apply to the provider configuration to assist with categorization and organization.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSIdentityProviderConfig,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSIdentityProviderConfig,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IdentityProviderConfigArn").Transform(arnToAkas),
				Hydrate:     getEKSIdentityProviderConfig,
			},
		}),
	}
}

//// LIST FUNCTION

func listEKSIdentityProviderConfigs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get Eks Cluster details
	cluster := h.Item.(types.Cluster)

	// Create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_identity_provider_config.listEKSIdentityProviderConfigs", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// As per the API document input parameter MaxResults should support the value of 100.
	// However with value of 100, API is throwing an error - InvalidParameterException: maxResults needs to be 1.
	// Raised an issue with AWS SDK - https://github.com/aws/aws-sdk-go/issues/4457
	// Same behaviour in AWS SDK V2 also.
	param := &eks.ListIdentityProviderConfigsInput{
		ClusterName: cluster.Name,
	}

	paginator := eks.NewListIdentityProviderConfigsPaginator(svc, param, func(o *eks.ListIdentityProviderConfigsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_identity_provider_config.listEKSIdentityProviderConfigs", "api_error", err)
			return nil, err
		}

		for _, providerConfig := range output.IdentityProviderConfigs {
			plugin.Logger(ctx).Info("providerConfig ", providerConfig)
			d.StreamListItem(ctx, &IdentityProviderConfig{providerConfig.Name, providerConfig.Type, types.OidcIdentityProviderConfig{
				ClusterName: cluster.Name,
			}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEKSIdentityProviderConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var clusterName, providerConfigName, providerConfigType string
	if h.Item != nil {
		clusterName = *h.Item.(*IdentityProviderConfig).ClusterName
		providerConfigName = *h.Item.(*IdentityProviderConfig).Name
		providerConfigType = *h.Item.(*IdentityProviderConfig).Type
	} else {
		clusterName = d.KeyColumnQuals["cluster_name"].GetStringValue()
		providerConfigName = d.KeyColumnQuals["name"].GetStringValue()
		providerConfigType = d.KeyColumnQuals["type"].GetStringValue()
	}

	// create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_identity_provider_config.getEKSIdentityProviderConfig", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &eks.DescribeIdentityProviderConfigInput{
		ClusterName: &clusterName,
		IdentityProviderConfig: &types.IdentityProviderConfig{
			Name: aws.String(providerConfigName),
			Type: aws.String(providerConfigType),
		},
	}

	op, err := svc.DescribeIdentityProviderConfig(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_identity_provider_config.getEKSIdentityProviderConfig", "api_error", err)
		return nil, err
	}

	return &IdentityProviderConfig{
		&providerConfigName,
		&providerConfigType,
		*op.IdentityProviderConfig.Oidc,
	}, nil
}
