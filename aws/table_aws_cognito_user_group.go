package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCognitoUserGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cognito_user_group",
		Description: "AWS Cognito User Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"group_name", "user_pool_id"}),
			Hydrate:    getCognitoUserGroup,
			Tags:       map[string]string{"service": "cognito-idp", "action": "GetGroup"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCognitoUserPools,
			Hydrate:       listCognitoUserGroups,
			Tags:          map[string]string{"service": "cognito-idp", "action": "ListGroups"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_pool_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_COGNITO_IDP_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "group_name",
				Description: "The name of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_pool_id",
				Description: "The user pool ID for the user pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A string containing the description of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The date and time when the group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_modified_date",
				Description: "The date and time when the group was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastModifiedDate").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "precedence",
				Description: "A non-negative integer value that specifies the precedence of this group relative to the other groups that a user can belong to in the user pool.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "role_arn",
				Description: "The role Amazon Resource Name (ARN) for the group.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCognitoUserGroupAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listCognitoUserGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get user pool details from hydrate data
	var userPoolID string
	if h.Item != nil {
		userPoolData, ok := h.Item.(types.UserPoolDescriptionType)
		if ok {
			userPoolID = *userPoolData.Id
		}
	}

	// Check if a specific user_pool_id has been provided in the query
	if d.EqualsQualString("user_pool_id") != "" && userPoolID == d.EqualsQualString("user_pool_id") {
		return nil, nil
	}

	// Return if no user pool is found
	if userPoolID == "" {
		return nil, nil
	}

	// Create session
	svc, err := CognitoIdentityProviderClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_user_group.listCognitoUserGroups", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Limiting the results
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

	input := &cognitoidentityprovider.ListGroupsInput{
		UserPoolId: aws.String(userPoolID),
		Limit:      aws.Int32(maxLimit),
	}

	// List call
	paginator := cognitoidentityprovider.NewListGroupsPaginator(svc, input, func(o *cognitoidentityprovider.ListGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cognito_user_group.listCognitoUserGroups", "api_error", err)
			return nil, err
		}

		for _, item := range output.Groups {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCognitoUserGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var groupName, userPoolID string

	if h.Item != nil {
		data := h.Item.(types.GroupType)
		groupName = *data.GroupName
		userPoolID = *data.UserPoolId
	} else {
		groupName = d.EqualsQualString("group_name")
		userPoolID = d.EqualsQualString("user_pool_id")
	}

	// Empty check for required parameters
	if groupName == "" || userPoolID == "" {
		return nil, nil
	}

	// Create service
	svc, err := CognitoIdentityProviderClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_user_group.getCognitoUserGroup", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &cognitoidentityprovider.GetGroupInput{
		GroupName:  aws.String(groupName),
		UserPoolId: aws.String(userPoolID),
	}

	// Get call
	data, err := svc.GetGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_user_group.getCognitoUserGroup", "api_error", err)
		return nil, err
	}

	return *data.Group, nil
}

func getCognitoUserGroupAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(types.GroupType)
	region := d.EqualsQualString(matrixKeyRegion)

	// Get account details
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cognito_user_group.getCognitoUserGroupAkas", "common_data_error", err)
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	accountID := commonColumnData.AccountId

	// Generate user group ARN
	userGroupArn := "arn:aws:cognito-idp:" + region + ":" + accountID + ":userpool/" + *group.UserPoolId + "/group/" + *group.GroupName

	return []string{userGroupArn}, nil
}
