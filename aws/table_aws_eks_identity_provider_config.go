package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/service/eks"
)

type IdentityProviderConfig struct {
	Name *string
	Type *string
	eks.OidcIdentityProviderConfig
}

//// TABLE DEFINITION

func tableAwsEksIdentityProviderConfig(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_identity_provider_config",
		Description: "AWS EKS Identity Provider Config",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "type", "cluster_name"}),
			Hydrate:    getEksIdentityProviderConfig,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEksClusters,
			Hydrate:       listEksIdentityProviderConfigs,
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Hydrate:     getEksIdentityProviderConfig,
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
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "groups_claim",
				Description: "The JSON web token (JWT) claim that the provider uses to return your groups.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "groups_prefix",
				Description: "The prefix that is prepended to group claims to prevent clashes with existing names (such as system: groups).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "issuer_url",
				Description: "The URL of the OIDC identity provider that allows the API server to discover public signing keys for verifying tokens.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "username_claim",
				Description: "The JSON Web token (JWT) claim that is used as the username.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "status",
				Description: "The status of the OIDC identity provider.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "username_prefix",
				Description: "The prefix that is prepended to username claims to prevent clashes with existing names.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "required_claims",
				Description: "The key-value pairs that describe required claims in the identity token.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "tags_src",
				Description: "The metadata to apply to the provider configuration to assist with categorization and organization.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksIdentityProviderConfig,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksIdentityProviderConfig,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IdentityProviderConfigArn").Transform(arnToAkas),
				Hydrate:     getEksIdentityProviderConfig,
			},
		}),
	}
}

//// LIST FUNCTION

func listEksIdentityProviderConfigs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get Eks Cluster details
	cluster := h.Item.(*eks.Cluster)

	// Create service
	svc, err := EksService(ctx, d)
	if err != nil {
		return nil, err
	}

	// As per the API document input parameter MaxResults should support the value of 100.
	// However with value of 100, API is throwing an error - InvalidParameterException: maxResults needs to be 1.
	// Raised an issue with AWS SDK - https://github.com/aws/aws-sdk-go/issues/4457
	param := &eks.ListIdentityProviderConfigsInput{
		ClusterName: cluster.Name,
		//MaxResults:  aws.Int64(100),
	}

	// limit := d.QueryContext.Limit
	// if d.QueryContext.Limit != nil {
	// 	if *limit < *param.MaxResults {
	// 		if *limit < 1 {
	// 			param.MaxResults = aws.Int64(1)
	// 		} else {
	// 			param.MaxResults = limit
	// 		}
	// 	}
	// }

	err = svc.ListIdentityProviderConfigsPages(
		param,
		func(page *eks.ListIdentityProviderConfigsOutput, _ bool) bool {
			for _, providerConfig := range page.IdentityProviderConfigs {
				d.StreamListItem(ctx, &IdentityProviderConfig{providerConfig.Name, providerConfig.Type, eks.OidcIdentityProviderConfig{
					ClusterName: cluster.Name,
				}})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return true
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEksIdentityProviderConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEksIdentityProviderConfig")

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
	svc, err := EksService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &eks.DescribeIdentityProviderConfigInput{
		ClusterName: &clusterName,
		IdentityProviderConfig: &eks.IdentityProviderConfig{
			Name: &providerConfigName,
			Type: &providerConfigType,
		},
	}

	op, err := svc.DescribeIdentityProviderConfig(params)
	if err != nil {
		return nil, err
	}

	return &IdentityProviderConfig{
		&providerConfigName,
		&providerConfigType,
		*op.IdentityProviderConfig.Oidc,
	}, nil
}
