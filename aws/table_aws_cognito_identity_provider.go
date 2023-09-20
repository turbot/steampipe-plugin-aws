package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	cognitoidentityproviderv1 "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type identityProviderInfo = struct {
	types.ProviderDescription
	UserPoolId *string
}

//// TABLE DEFINITION

func tableAwsCognitoIdentityProvider(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cognito_identity_provider",
		Description: "AWS Cognito Identity Provider",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"provider_name", "user_pool_id"}),
			Hydrate:    getCognitoIdentityProvider,
			Tags:       map[string]string{"service": "cognito-identity", "action": "DescribeIdentityProvider"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCognitoUserPools,
			Hydrate:       listCognitoIdentityProviders,
			Tags:          map[string]string{"service": "cognito-identity", "action": "ListIdentityProviders"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_pool_id", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCognitoIdentityProvider,
				Tags: map[string]string{"service": "cognito-identity", "action": "DescribeIdentityProvider"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cognitoidentityproviderv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "provider_name",
				Description: "The IdP name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_pool_id",
				Description: "The user pool ID.",
				Hydrate:     getCognitoIdentityProvider,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attribute_mapping",
				Description: "A mapping of IdP attributes to standard and custom user pool attributes.",
				Hydrate:     getCognitoIdentityProvider,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "creation_date",
				Description: "The date the provider was added to the user pool.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "idp_identifiers",
				Description: "A list of IdP identifiers.",
				Hydrate:     getCognitoIdentityProvider,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_modified_date",
				Description: "The date the provider was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			{
				Name:        "provider_details",
				Description: "The IdP details.",
				Hydrate:     getCognitoIdentityProvider,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "provider_type",
				Description: "The IdP type.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCognitoIdentityProviderTurbotAkas,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProviderName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCognitoIdentityProviders(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CognitoIdentityProviderClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_identity_provider.listCognitoIdentityProviders", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		plugin.Logger(ctx).Debug("aws_cognito_identity_provider.listCognitoIdentityProviders", "unsupported_region")
		return nil, nil
	}

	userPoolId := h.Item.(types.UserPoolDescriptionType).Id
	plugin.Logger(ctx).Debug("aws_cognito_identity_provider.listCognitoIdentityProviders", "user_pool_id", *userPoolId)

	equalQuals := d.EqualsQuals
	// Minimize the API call with the given user_pool_id
	if equalQuals["user_pool_id"] != nil {
		if equalQuals["user_pool_id"].GetStringValue() != *userPoolId {
			plugin.Logger(ctx).Debug("aws_cognito_identity_provider.listCognitoIdentityProviders", "ignoring mismatching user pool id", userPoolId)
			return nil, nil
		}
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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

	input := &cognitoidentityprovider.ListIdentityProvidersInput{
		MaxResults: aws.Int32(maxLimit),
		UserPoolId: userPoolId,
	}
	// List call
	paginator := cognitoidentityprovider.NewListIdentityProvidersPaginator(svc, input, func(o *cognitoidentityprovider.ListIdentityProvidersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {

		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cognito_identity_provider.listCognitoIdentityProviders", "api_error", err)
			return nil, err
		}
		plugin.Logger(ctx).Debug("aws_cognito_identity_provider.listCognitoIdentityProviders", "providers", fmt.Sprintf("%#v", output.Providers))

		for _, provider := range output.Providers {
			plugin.Logger(ctx).Debug("aws_cognito_identity_provider.listCognitoIdentityProviders", "provider", fmt.Sprintf("%#v", provider))
			d.StreamListItem(ctx, identityProviderInfo{provider, userPoolId})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCognitoIdentityProvider(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var providerName string
	var userPoolId string
	if h.Item != nil {
		data := h.Item.(identityProviderInfo)
		if data.ProviderName != nil {
			providerName = *data.ProviderName
		}
		if data.UserPoolId != nil {
			userPoolId = *data.UserPoolId
		}
	} else {
		providerName = d.EqualsQualString("provider_name")
		userPoolId = d.EqualsQualString("user_pool_id")
	}

	plugin.Logger(ctx).Debug("aws_cognito_identity_provider.getCognitoIdentityProvider", "provider_name", providerName, "user_pool_id", userPoolId)

	// check if providerName or userPoolId is empty
	if providerName == "" || userPoolId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := CognitoIdentityProviderClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_identity_provider.getCognitoIdentityProvider", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &cognitoidentityprovider.DescribeIdentityProviderInput{
		ProviderName: aws.String(providerName),
		UserPoolId:   aws.String(userPoolId),
	}

	// Get call
	data, err := svc.DescribeIdentityProvider(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_identity_provider.getCognitoIdentityProvider", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("aws_cognito_identity_provider.getCognitoIdentityProvider", "identity_provider", fmt.Sprintf("%#v", *data.IdentityProvider))
	return *data.IdentityProvider, nil
}

func getCognitoIdentityProviderTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	userPoolId := d.EqualsQualString("user_pool_id")
	data := h.Item.(identityProviderInfo)
	plugin.Logger(ctx).Debug("aws_cognito_identity_provider.getCognitoIdentityProviderTurbotAkas", "user_pool_id", userPoolId, "provider_name", *data.ProviderName)

	// Get common columns

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_identity_provider.getCognitoIdentityProviderTurbotAkas", "common_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// Get data for turbot defined properties
	//arn:aws:cognito-idp:<region>:<account-id>:userpool/<id>/provider/<name>
	arn := "arn:" + commonColumnData.Partition + ":cognito-idp:" + region + ":" + commonColumnData.AccountId + ":userpool/" + userPoolId + "/provider/" + *data.ProviderName

	return []string{arn}, nil
}
