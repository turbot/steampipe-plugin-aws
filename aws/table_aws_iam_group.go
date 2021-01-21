package aws

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsIamGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_group",
		Description: "AWS IAM Group",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("name"),
			ItemFromKey: groupFromKey,
			Hydrate:     getIamGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listIamGroups,
		},
		HydrateDependencies: []plugin.HydrateDependencies{
			{
				Func:    getAwsIamGroupInlinePolicies,
				Depends: []plugin.HydrateFunc{listAwsIamGroupInlinePolicies},
			},
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the group",
				Type:        proto.ColumnType_STRING, Transform: transform.FromField("GroupName"),
			},
			{
				Name:        "group_id",
				Description: "The stable and unique string identifying the group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "The path to the group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the group",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The date and time, when the group was created",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "inline_policies",
				Description: "A list of policy documents that are embedded as inline policies for the group",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamGroupInlinePolicies,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "attached_policy_arns",
				Description: "A list of managed policies attached to the group",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamGroupAttachedPolicies,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "users",
				Description: "A list of users in the group",
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

//// ITEM FROM KEY

func groupFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &iam.Group{
		GroupName: &name,
	}
	return item, nil
}

//// LIST FUNCTION

func listIamGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listIamGroups")

	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	err = svc.ListGroupsPages(
		&iam.ListGroupsInput{},
		func(page *iam.ListGroupsOutput, lastPage bool) bool {
			for _, group := range page.Groups {
				d.StreamListItem(ctx, group)
			}
			return true
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getIamGroup")
	group := h.Item.(*iam.Group)

	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	params := &iam.GetGroupInput{
		GroupName: group.GroupName,
	}

	op, err := svc.GetGroup(params)
	if err != nil {
		logger.Debug("getIamGroup__", "ERROR", err)
		return nil, err
	}

	return op.Group, nil
}

func getAwsIamGroupAttachedPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsIamGroupAttachedPolicies")
	group := h.Item.(*iam.Group)

	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	params := &iam.ListAttachedGroupPoliciesInput{
		GroupName: group.GroupName,
	}

	groupData, err := svc.ListAttachedGroupPolicies(params)
	if err != nil {
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
	plugin.Logger(ctx).Trace("getAwsIamGroupUsers")
	group := h.Item.(*iam.Group)

	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	params := &iam.GetGroupInput{
		GroupName: group.GroupName,
	}

	groupData, err := svc.GetGroup(params)
	if err != nil {
		return nil, err
	}

	if groupData.Users != nil {
		return groupData, nil
	}
	return iam.GetGroupOutput{}, nil
}

func listAwsIamGroupInlinePolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsIamGroupInlinePolicies")
	group := h.Item.(*iam.Group)

	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	params := &iam.ListGroupPoliciesInput{
		GroupName: group.GroupName,
	}

	groupData, err := svc.ListGroupPolicies(params)
	if err != nil {
		return nil, err
	}

	return groupData, nil
}

func getAwsIamGroupInlinePolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsIamGroupInlinePolicies")
	group := h.Item.(*iam.Group)
	listGroupPoliciesOutput := h.HydrateResults["listAwsIamGroupInlinePolicies"].(*iam.ListGroupPoliciesOutput)

	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	policyCh := make(chan map[string]interface{}, len(listGroupPoliciesOutput.PolicyNames))
	errorCh := make(chan error, len(listGroupPoliciesOutput.PolicyNames))
	for _, policy := range listGroupPoliciesOutput.PolicyNames {
		wg.Add(1)
		go getGroupPolicyDataAsync(policy, group.GroupName, svc, &wg, policyCh, errorCh)
	}

	// wait for all inline policies to be processed
	wg.Wait()
	// NOTE: close channel before ranging over results
	close(policyCh)
	close(errorCh)

	for err := range errorCh {
		// return the first error
		return nil, err
	}

	var groupPolicies []map[string]interface{}

	for groupPolicy := range policyCh {
		groupPolicies = append(groupPolicies, groupPolicy)
	}

	return groupPolicies, nil
}

func getGroupPolicyDataAsync(policy *string, groupName *string, svc *iam.IAM, wg *sync.WaitGroup, policyCh chan map[string]interface{}, errorCh chan error) {
	defer wg.Done()

	rowData, err := getGroupInlinePolicy(policy, groupName, svc)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		policyCh <- rowData
	}
}

func getGroupInlinePolicy(policyName *string, groupName *string, svc *iam.IAM) (map[string]interface{}, error) {
	groupPolicy := make(map[string]interface{})
	params := &iam.GetGroupPolicyInput{
		PolicyName: policyName,
		GroupName:  groupName,
	}

	tmpPolicy, err := svc.GetGroupPolicy(params)
	if err != nil {
		return nil, err
	}

	if tmpPolicy != nil && tmpPolicy.PolicyDocument != nil {
		decoded, decodeErr := url.QueryUnescape(*tmpPolicy.PolicyDocument)
		if decodeErr != nil {
			return nil, decodeErr
		}

		var rawPolicy interface{}
		err := json.Unmarshal([]byte(decoded), &rawPolicy)
		if err != nil {
			return nil, err
		}

		groupPolicy = map[string]interface{}{
			"PolicyDocument": rawPolicy,
			"PolicyName":     *tmpPolicy.PolicyName,
		}
	}

	return groupPolicy, nil
}
