package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"

	cognitoidentityv1 "github.com/aws/aws-sdk-go/service/cognitoidentity"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCognitoIdentityPool(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cognito_identity_pool",
		Description: "AWS Cognito Identity Pool",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("identity_pool_id"),
			Hydrate:    getCognitoIdentityPool,
			Tags:       map[string]string{"service": "cognito-identity", "action": "DescribeIdentityPool"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCognitoIdentityPools,
			Tags:    map[string]string{"service": "cognito-identity", "action": "ListIdentityPools"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCognitoIdentityPool,
				Tags: map[string]string{"service": "cognito-identity", "action": "DescribeIdentityPool"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cognitoidentityv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "identity_pool_id",
				Description: "An identity pool ID in the format REGION:GUID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allow_unauthenticated_identities",
				Description: "TRUE if the identity pool supports unauthenticated logins.",
				Hydrate:     getCognitoIdentityPool,
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "identity_pool_name",
				Description: "A string that you provide.",
				Hydrate:     getCognitoIdentityPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allow_classic_flow",
				Description: "Enables or disables the Basic (Classic) authentication flow.",
				Hydrate:     getCognitoIdentityPool,
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "cognito_identity_providers",
				Description: "A list representing an Amazon Cognito user pool and its client ID.",
				Hydrate:     getCognitoIdentityPool,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "developer_provider_name",
				Description: "The 'domain' by which Cognito will refer to your users.",
				Hydrate:     getCognitoIdentityPool,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "open_id_connect_provider_arns",
				Description: "The ARNs of the OpenID Connect providers.",
				Hydrate:     getCognitoIdentityPool,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OpenIdConnectProviderARNs"),
			},
			{
				Name:        "saml_provider_arns",
				Description: "An array of Amazon Resource Names (ARNs) of the SAML provider for your identity pool.",
				Hydrate:     getCognitoIdentityPool,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SamlProviderARNs"),
			},
			{
				Name:        "supported_login_providers",
				Description: "Optional key:value pairs mapping provider names to provider app IDs.",
				Hydrate:     getCognitoIdentityPool,
				Type:        proto.ColumnType_JSON,
			},
			// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCognitoIdentityPoolTurbotAkas,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCognitoIdentityPool,
				Transform:   transform.FromField("IdentityPoolTags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IdentityPoolName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCognitoIdentityPools(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CognitoIdentityClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_identity_pool.listCognitoIdentityPools", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		plugin.Logger(ctx).Debug("aws_cognito_identity_pool.listCognitoIdentityPools", "unsupported_region")
		return nil, nil
	}

	// Reduce the basic request limit down if the identity has only requested a small number of rows
	maxLimit := int32(60)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	input := &cognitoidentity.ListIdentityPoolsInput{
		MaxResults: maxLimit,
	}
	// List call
	paginator := cognitoidentity.NewListIdentityPoolsPaginator(svc, input, func(o *cognitoidentity.ListIdentityPoolsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {

		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cognito_identity_pool.listCognitoIdentityPools", "api_error", err)
			return nil, err
		}
		plugin.Logger(ctx).Debug("aws_cognito_identity_pool.listCognitoIdentityPools", "identity_pools", fmt.Sprintf("%#v", output.IdentityPools))

		for _, identityPool := range output.IdentityPools {
			plugin.Logger(ctx).Debug("aws_cognito_identity_pool.listCognitoIdentityPools", "identity_pool", fmt.Sprintf("%#v", identityPool))
			d.StreamListItem(ctx, identityPool)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCognitoIdentityPool(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var identityPoolId string
	if h.Item != nil {
		data := h.Item.(types.IdentityPoolShortDescription)
		if data.IdentityPoolId != nil {
			identityPoolId = *data.IdentityPoolId
		}
	} else {
		identityPoolId = d.EqualsQualString("identity_pool_id")
	}
	plugin.Logger(ctx).Debug("aws_cognito_identity_pool.getCognitoIdentityPool", "identity_pool_id", identityPoolId)

	// check if id is empty
	if identityPoolId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := CognitoIdentityClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_identity_pool.getCognitoIdentityPool", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &cognitoidentity.DescribeIdentityPoolInput{
		IdentityPoolId: aws.String(identityPoolId),
	}

	// Get call
	data, err := svc.DescribeIdentityPool(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_identity_pool.getCognitoIdentityPool", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("aws_cognito_identity_pool.getCognitoIdentityPool", "identity_pool", fmt.Sprintf("%#v", *data))
	return *data, nil
}

func getCognitoIdentityPoolTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	data := h.Item.(types.IdentityPoolShortDescription)
	plugin.Logger(ctx).Debug("aws_cognito_identity_pool.getCognitoIdentityPoolTurbotAkas", "identity_pool_id", *data.IdentityPoolId)

	// Get common columns

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_identity_pool.getCognitoIdentityPoolTurbotAkas", "common_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// Get data for turbot defined properties
	//arn:aws:cognito-identity:<region>:<account-id>:identitypool/<id>
	arn := "arn:" + commonColumnData.Partition + ":cognito-identity:" + region + ":" + commonColumnData.AccountId + ":identitypool/" + *data.IdentityPoolId

	return []string{arn}, nil
}
