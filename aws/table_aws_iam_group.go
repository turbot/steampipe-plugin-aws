package aws

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_group",
		Description: "AWS IAM Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"name", "arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationError", "NoSuchEntity", "InvalidParameter"}),
			},
			Hydrate: getIamGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listIamGroups,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "path", Require: plugin.Optional},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the group.",
				Type:        proto.ColumnType_STRING, Transform: transform.FromField("GroupName"),
			},
			{
				Name:        "group_id",
				Description: "The stable and unique string identifying the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "The path to the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The date and time, when the group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "inline_policies",
				Description: "A list of policy documents that are embedded as inline policies for the group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsIamGroupInlinePolicies,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "inline_policies_std",
				Description: "Inline policies in canonical form for the group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsIamGroupInlinePolicies,
				Transform:   transform.FromValue().Transform(inlinePoliciesToStd),
			},
			{
				Name:        "attached_policy_arns",
				Description: "A list of managed policies attached to the group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamGroupAttachedPolicies,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "users",
				Description: "A list of users in the group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamGroupUsers,
			},

			/// Standard columns for all tables
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
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listIamGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.listIamGroups", "client_error", err)
		return nil, err
	}

	input := &iam.ListGroupsInput{}

	equalQual := d.EqualsQuals
	if equalQual["path"] != nil {
		input.PathPrefix = aws.String(equalQual["path"].GetStringValue())
	}

	maxItems := int32(1000)
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
	paginator := iam.NewListGroupsPaginator(svc, input, func(o *iam.ListGroupsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aaws_iam_group.listIamGroups", "api_error", err)
			return nil, err
		}

		for _, group := range output.Groups {
			d.StreamListItem(ctx, group)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQuals["arn"].GetStringValue()
	groupName := d.EqualsQuals["name"].GetStringValue()
	if len(arn) > 0 {
		groupName = strings.Split(arn, "/")[len(strings.Split(arn, "/"))-1]
	}

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.getIamGroup", "client_error", err)
		return nil, err
	}

	params := &iam.GetGroupInput{
		GroupName: aws.String(groupName),
	}

	op, err := svc.GetGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.getIamGroup", "api_error", err)
		return nil, err
	}

	return *op.Group, nil
}

func getAwsIamGroupAttachedPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(types.Group)

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.listIamGroups", "client_error", err)
		return nil, err
	}

	params := &iam.ListAttachedGroupPoliciesInput{
		GroupName: group.GroupName,
	}

	groupData, err := svc.ListAttachedGroupPolicies(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.getAwsIamGroupAttachedPolicies", "api_error", err)
		return nil, err
	}

	var attachedPolicyArns []string

	if groupData.AttachedPolicies != nil {
		for _, policy := range groupData.AttachedPolicies {
			attachedPolicyArns = append(attachedPolicyArns, *policy.PolicyArn)
		}
	}

	return attachedPolicyArns, nil
}

func getAwsIamGroupUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(types.Group)

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.getAwsIamGroupUsers", "client_error", err)
		return nil, err
	}

	params := &iam.GetGroupInput{GroupName: group.GroupName}

	groupData, err := svc.GetGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.getAwsIamGroupUsers", "api_error", err)
		return nil, err
	}

	if len(groupData.Users) > 0 {
		return groupData, nil
	}
	return nil, nil
}

func listAwsIamGroupInlinePolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(types.Group)

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.listAwsIamGroupInlinePolicies", "client_error", err)
		return nil, err
	}

	params := &iam.ListGroupPoliciesInput{GroupName: group.GroupName}
	groupData, err := svc.ListGroupPolicies(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.listAwsIamGroupInlinePolicies", "api_error", err)
		return nil, err
	}

	var wg sync.WaitGroup
	policyCh := make(chan map[string]interface{}, len(groupData.PolicyNames))
	errorCh := make(chan error, len(groupData.PolicyNames))
	for _, policy := range groupData.PolicyNames {
		wg.Add(1)
		go getGroupPolicyDataAsync(ctx, policy, group.GroupName, svc, &wg, policyCh, errorCh)
	}

	// wait for all inline policies to be processed
	wg.Wait()

	// NOTE: close channel before ranging over results
	close(policyCh)
	close(errorCh)

	for err := range errorCh {
		// return the first error
		plugin.Logger(ctx).Error("aws_iam_group.listAwsIamGroupInlinePolicies", "channel_error", err)
		return nil, err
	}

	var groupPolicies []map[string]interface{}

	for groupPolicy := range policyCh {
		groupPolicies = append(groupPolicies, groupPolicy)
	}

	return groupPolicies, nil
}

func getGroupPolicyDataAsync(ctx context.Context, policy string, groupName *string, svc *iam.Client, wg *sync.WaitGroup, policyCh chan map[string]interface{}, errorCh chan error) {
	defer wg.Done()

	rowData, err := getGroupInlinePolicy(ctx, policy, groupName, svc)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		policyCh <- rowData
	}
}

func getGroupInlinePolicy(ctx context.Context, policyName string, groupName *string, svc *iam.Client) (map[string]interface{}, error) {
	groupPolicy := make(map[string]interface{})
	params := &iam.GetGroupPolicyInput{
		PolicyName: &policyName,
		GroupName:  groupName,
	}

	tmpPolicy, err := svc.GetGroupPolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_group.getGroupInlinePolicy", "api_error", err)
		return nil, err
	}

	if tmpPolicy != nil && tmpPolicy.PolicyDocument != nil {
		decoded, decodeErr := url.QueryUnescape(*tmpPolicy.PolicyDocument)
		if decodeErr != nil {
			plugin.Logger(ctx).Error("aws_iam_group.getGroupInlinePolicy", "decode_error", err)
			return nil, decodeErr
		}

		var rawPolicy interface{}
		err := json.Unmarshal([]byte(decoded), &rawPolicy)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_group.getGroupInlinePolicy", "unmarshal_error", err)
			return nil, err
		}

		groupPolicy = map[string]interface{}{
			"PolicyDocument": rawPolicy,
			"PolicyName":     *tmpPolicy.PolicyName,
		}
	}

	return groupPolicy, nil
}
