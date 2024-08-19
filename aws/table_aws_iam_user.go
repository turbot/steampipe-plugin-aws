package aws

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_user",
		Description: "AWS IAM User",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"name", "arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationError", "NoSuchEntity", "InvalidParameter"}),
			},
			Hydrate: getIamUser,
			Tags:    map[string]string{"service": "iam", "action": "GetUser"},
		},
		List: &plugin.ListConfig{
			Hydrate: listIamUsers,
			Tags:    map[string]string{"service": "iam", "action": "ListUsers"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "path", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsIamUserLoginProfile,
				Tags: map[string]string{"service": "iam", "action": "GetLoginProfile"},
			},
			{
				Func: getAwsIamUserData,
				Tags: map[string]string{"service": "iam", "action": "GetUser"},
			},
			{
				Func: getAwsIamUserAttachedPolicies,
				Tags: map[string]string{"service": "iam", "action": "ListAttachedUserPolicies"},
			},
			{
				Func: getAwsIamUserGroups,
				Tags: map[string]string{"service": "iam", "action": "ListGroupsForUser"},
			},
			{
				Func: getAwsIamUserMfaDevices,
				Tags: map[string]string{"service": "iam", "action": "ListMFADevices"},
			},
			{
				Func: listAwsIamUserInlinePolicies,
				Tags: map[string]string{"service": "iam", "action": "ListUserPolicies"},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name identifying the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserName"),
			},
			{
				Name:        "user_id",
				Description: "The stable and unique string identifying the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "The path to the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that identifies the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The date and time, when the user was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "password_last_used",
				Description: "The date and time, when the user's password was last used to sign in to an AWS website.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "permissions_boundary_arn",
				Description: "The ARN of the policy used to set the permissions boundary for the user.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsIamUserData,
			},
			{
				Name: "permissions_boundary_type",
				Description: "The permissions boundary usage type that indicates what type of IAM resource " +
					"is used as the permissions boundary for an entity. This data type can only have a value of Policy.",
				Type:    proto.ColumnType_STRING,
				Hydrate: getAwsIamUserData,
			},
			{
				Name:        "mfa_enabled",
				Description: "The MFA status of the user.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsIamUserMfaDevices,
				Transform:   transform.From(handleEmptyUserMfaStatus),
			},
			{
				Name:        "login_profile",
				Description: "Contains the user name and password create date for a user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamUserLoginProfile,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "mfa_devices",
				Description: "A list of MFA devices attached to the user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamUserMfaDevices,
				Transform:   transform.FromField("MFADevices"),
			},
			{
				Name:        "groups",
				Description: "A list of groups attached to the user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamUserGroups,
			},
			{
				Name:        "inline_policies",
				Description: "A list of policy documents that are embedded as inline policies for the user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsIamUserInlinePolicies,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "inline_policies_std",
				Description: "Inline policies in canonical form for the user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsIamUserInlinePolicies,
				Transform:   transform.FromValue().Transform(inlinePoliciesToStd),
			},
			{
				Name:        "attached_policy_arns",
				Description: "A list of managed policies attached to the user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamUserAttachedPolicies,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamUserData,
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamUserData,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listIamUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_role.listIamRoles", "client_error", err)
		return nil, err
	}

	maxItems := int32(1000)
	input := iam.ListUsersInput{}

	equalQual := d.EqualsQuals
	if equalQual["path"] != nil {
		input.PathPrefix = aws.String(equalQual["path"].GetStringValue())
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxItems = aws.Int32(maxItems)
	paginator := iam.NewListUsersPaginator(svc, &input, func(o *iam.ListUsersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_role.listIamRoles", "api_error", err)
			return nil, err
		}

		for _, user := range output.Users {
			d.StreamListItem(ctx, user)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	arn := d.EqualsQuals["arn"].GetStringValue()
	name := d.EqualsQuals["name"].GetStringValue()
	if len(arn) > 0 {
		name = strings.Split(arn, "/")[len(strings.Split(arn, "/"))-1]
	}

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getIamUser", "client_error", err)
		return nil, err
	}

	params := &iam.GetUserInput{
		UserName: aws.String(name),
	}

	op, err := svc.GetUser(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getIamUser", "api_error", err)
		return nil, err
	}

	return *op.User, nil
}

func getAwsIamUserLoginProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := h.Item.(types.User).UserName

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserLoginProfile", "client_error", err)
		return nil, err
	}

	params := &iam.GetLoginProfileInput{
		UserName: name,
	}

	op, err := svc.GetLoginProfile(ctx, params)
	if err != nil {
		// If the user does not exist or does not have a password, the operation returns a 404 (NoSuchEntity) error.
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchEntity" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserLoginProfile", "api_error", err)
		return nil, err
	}

	return op.LoginProfile, nil
}

func getAwsIamUserData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(types.User)

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserData", "client_error", err)
		return nil, err
	}

	params := &iam.GetUserInput{
		UserName: user.UserName,
	}

	userData, _ := svc.GetUser(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserData", "api_error", err)
		return nil, err
	}

	var tags []types.Tag
	var turbotTags map[string]string
	PermissionsBoundaryArn := ""
	PermissionsBoundaryType := ""

	if userData.User != nil && userData.User.Tags != nil {
		tags = userData.User.Tags
		turbotTags = map[string]string{}
		for _, t := range userData.User.Tags {
			turbotTags[*t.Key] = *t.Value
		}
	}

	if userData.User != nil && userData.User.PermissionsBoundary != nil && userData.User.PermissionsBoundary.PermissionsBoundaryArn != nil {
		v := userData.User.PermissionsBoundary
		PermissionsBoundaryArn = *v.PermissionsBoundaryArn
		PermissionsBoundaryType = string(v.PermissionsBoundaryType)
	}

	return map[string]interface{}{
		"TagsSrc":                 tags,
		"Tags":                    turbotTags,
		"PermissionsBoundaryArn":  PermissionsBoundaryArn,
		"PermissionsBoundaryType": PermissionsBoundaryType,
	}, nil
}

func getAwsIamUserAttachedPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(types.User)

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserAttachedPolicies", "client_error", err)
		return nil, err
	}

	params := &iam.ListAttachedUserPoliciesInput{
		UserName: user.UserName,
	}

	userData, _ := svc.ListAttachedUserPolicies(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserAttachedPolicies", "api_error", err)
		return nil, err
	}

	var attachedPolicyArns []string

	if userData.AttachedPolicies != nil {
		for _, policy := range userData.AttachedPolicies {
			attachedPolicyArns = append(attachedPolicyArns, *policy.PolicyArn)
		}
	}

	return attachedPolicyArns, nil
}

func getAwsIamUserGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(types.User)

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserGroups", "client_error", err)
		return nil, err
	}

	params := &iam.ListGroupsForUserInput{
		UserName: user.UserName,
	}

	userData, _ := svc.ListGroupsForUser(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserGroups", "api_error", err)
		return nil, err
	}

	return userData, nil
}

func getAwsIamUserMfaDevices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(types.User)

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserMfaDevices", "client_error", err)
		return nil, err
	}

	params := &iam.ListMFADevicesInput{
		UserName: user.UserName,
	}

	userData, _ := svc.ListMFADevices(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getAwsIamUserMfaDevices", "api_error", err)
		return nil, err
	}

	return userData, nil
}

func listAwsIamUserInlinePolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(types.User)
	var userPolicies []map[string]interface{}
	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.listAwsIamUserInlinePolicies", "client_error", err)
		return nil, err
	}

	params := &iam.ListUserPoliciesInput{
		UserName: user.UserName,
	}

	userData, err := svc.ListUserPolicies(ctx, params)

	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.listAwsIamUserInlinePolicies", "api_error", err)
		return nil, err
	}

	var wg sync.WaitGroup
	policyCh := make(chan map[string]interface{}, len(userData.PolicyNames))
	errorCh := make(chan error, len(userData.PolicyNames))
	for _, policy := range userData.PolicyNames {
		wg.Add(1)
		go getUserPolicyDataAsync(ctx, aws.String(policy), user.UserName, svc, &wg, policyCh, errorCh)
	}

	// wait for all inline policies to be processed
	wg.Wait()

	// NOTE: close channel before ranging over results
	close(policyCh)
	close(errorCh)

	for err := range errorCh {
		// return the first error
		plugin.Logger(ctx).Error("aws_iam_user.listAwsIamUserInlinePolicies", "channel_error", err)
		return nil, err
	}

	for userPolicy := range policyCh {
		userPolicies = append(userPolicies, userPolicy)
	}

	return userPolicies, nil
}

func getUserPolicyDataAsync(ctx context.Context, policy *string, userName *string, svc *iam.Client, wg *sync.WaitGroup, policyCh chan map[string]interface{}, errorCh chan error) {
	defer wg.Done()

	rowData, err := getUserInlinePolicy(ctx, policy, userName, svc)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		policyCh <- rowData
	}
}

func getUserInlinePolicy(ctx context.Context, policyName *string, userName *string, svc *iam.Client) (map[string]interface{}, error) {
	userPolicy := make(map[string]interface{})
	params := &iam.GetUserPolicyInput{
		PolicyName: policyName,
		UserName:   userName,
	}

	tmpPolicy, err := svc.GetUserPolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_user.getUserInlinePolicy", "api_error", err)
		return nil, err
	}

	if tmpPolicy != nil && tmpPolicy.PolicyDocument != nil {
		decoded, decodeErr := url.QueryUnescape(*tmpPolicy.PolicyDocument)
		if decodeErr != nil {
			plugin.Logger(ctx).Error("aws_iam_user.getUserInlinePolicy", "decode_error", err)
			return nil, decodeErr
		}

		var rawPolicy interface{}
		err := json.Unmarshal([]byte(decoded), &rawPolicy)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_user.getUserInlinePolicy", "unmarshal_error", err)
			return nil, err
		}

		userPolicy = map[string]interface{}{
			"PolicyDocument": rawPolicy,
			"PolicyName":     *tmpPolicy.PolicyName,
		}
	}

	return userPolicy, nil
}

//// TRANSFORM FUNCTION

func handleEmptyUserMfaStatus(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*iam.ListMFADevicesOutput)
	if data != nil && len(data.MFADevices) > 0 {
		return true, nil
	}

	return false, nil
}
