package aws

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAwsIamRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_role",
		Description: "AWS IAM Role",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"name", "arn"}),
			Hydrate:    getIamRole,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationError", "NoSuchEntity", "InvalidParameter"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIamRoles,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "path", Require: plugin.Optional},
			},
		},
		Columns: awsColumns([]*plugin.Column{
			// "Key" Columns
			{
				Name:        "name",
				Description: "The friendly name that identifies the role.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleName"),
			},
			{
				Name:        "arn",
				Type:        proto.ColumnType_STRING,
				Description: "The Amazon Resource Name (ARN) specifying the role.",
			},
			{
				Name:        "role_id",
				Type:        proto.ColumnType_STRING,
				Description: "The stable and unique string identifying the role.",
			},

			// Other Columns
			{
				Name:        "create_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time when the role was created.",
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "A user-provided description of the role.",
			},
			{
				Name:        "instance_profile_arns",
				Description: "A list of instance profiles associated with the role.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamInstanceProfileData,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "max_session_duration",
				Description: "The maximum session duration (in seconds) for the specified role. Anyone who uses the AWS CLI, or API to assume the role can specify the duration using the optional DurationSeconds API parameter or duration-seconds CLI parameter.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "path",
				Description: "The path to the role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "permissions_boundary_arn",
				Description: "The ARN of the policy used to set the permissions boundary for the role.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIamRole,
				Transform:   transform.FromField("PermissionsBoundary.PermissionsBoundaryArn"),
			},
			{
				Name: "permissions_boundary_type",
				Description: "The permissions boundary usage type that indicates what type of IAM resource " +
					"is used as the permissions boundary for an entity. This data type can only have a value of Policy.",
				Type:      proto.ColumnType_STRING,
				Hydrate:   getIamRole,
				Transform: transform.FromField("PermissionsBoundary.PermissionsBoundaryType"),
			},
			{
				Name: "role_last_used_date",
				Type: proto.ColumnType_TIMESTAMP,
				Description: "Contains information about the last time that an IAM role was used. " +
					"Activity is only reported for the trailing 400 days. This period can be " +
					"shorter if your Region began supporting these features within the last year. " +
					"The role might have been used more than 400 days ago.",
				Hydrate:   getIamRole,
				Transform: transform.FromField("RoleLastUsed.LastUsedDate"),
			},
			{
				Name:        "role_last_used_region",
				Description: "Contains the region in which the IAM role was used.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIamRole,
				Transform:   transform.FromField("RoleLastUsed.Region"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the role.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
				Hydrate:     getIamRole,
			},
			{
				Name:        "inline_policies",
				Description: "A list of policy documents that are embedded as inline policies for the role..",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsIamRoleInlinePolicies,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "inline_policies_std",
				Description: "Inline policies in canonical form for the role.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsIamRoleInlinePolicies,
				Transform:   transform.FromValue().Transform(inlinePoliciesToStd),
			},
			{
				Name:        "attached_policy_arns",
				Description: "A list of managed policies attached to the role.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsIamRoleAttachedPolicies,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "assume_role_policy",
				Description: "The policy that grants an entity permission to assume the role.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssumeRolePolicyDocument").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "assume_role_policy_std",
				Description: "Contains the assume role policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssumeRolePolicyDocument").Transform(unescape).Transform(policyToCanonical),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleName"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getIamRoleTurbotTags),
				Hydrate:     getIamRole,
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

func listIamRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &iam.ListRolesInput{
		MaxItems: aws.Int64(1000),
	}

	equalQual := d.KeyColumnQuals
	if equalQual["path"] != nil {
		input.PathPrefix = aws.String(equalQual["path"].GetStringValue())
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			if *limit < 1 {
				input.MaxItems = aws.Int64(1)
			} else {
				input.MaxItems = limit
			}
		}
	}

	err = svc.ListRolesPages(
		input,
		func(page *iam.ListRolesOutput, lastPage bool) bool {

			for _, role := range page.Roles {
				d.StreamListItem(ctx, role)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamRole")

	// Create service
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	var name string
	if h.Item != nil {
		data := h.Item.(*iam.Role)
		name = types.SafeString(data.RoleName)
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		arn := d.KeyColumnQuals["arn"].GetStringValue()
		if len(arn) > 0 {
			name = strings.Split(arn, "/")[len(strings.Split(arn, "/"))-1]
		}
	}

	params := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}

	op, err := svc.GetRole(params)
	if err != nil {
		return nil, err
	}

	return op.Role, nil
}

func getAwsIamInstanceProfileData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsIamInstanceProfileData")
	role := h.Item.(*iam.Role)

	// create service
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	var associatedInstanceProfileArns []string

	params := &iam.ListInstanceProfilesForRoleInput{
		RoleName: role.RoleName,
	}

	err = svc.ListInstanceProfilesForRolePages(
		params,
		func(page *iam.ListInstanceProfilesForRoleOutput, lastPage bool) bool {
			for _, instanceProfile := range page.InstanceProfiles {
				associatedInstanceProfileArns = append(associatedInstanceProfileArns, *instanceProfile.Arn)
			}
			return !lastPage
		},
	)

	return associatedInstanceProfileArns, err
}

func getAwsIamRoleAttachedPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsIamRoleAttachedPolicies")
	role := h.Item.(*iam.Role)

	// create service
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &iam.ListAttachedRolePoliciesInput{
		RoleName: role.RoleName,
	}

	roleData, err := svc.ListAttachedRolePolicies(params)
	if err != nil {
		return nil, err
	}

	var attachedPolicyArns []string

	if roleData.AttachedPolicies != nil {
		for _, policy := range roleData.AttachedPolicies {
			attachedPolicyArns = append(attachedPolicyArns, *policy.PolicyArn)
		}
	}

	return attachedPolicyArns, nil
}

func listAwsIamRoleInlinePolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAwsIamRoleInlinePolicies")
	role := h.Item.(*iam.Role)

	// create service
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &iam.ListRolePoliciesInput{
		RoleName: role.RoleName,
	}

	roleData, err := svc.ListRolePolicies(params)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	policyCh := make(chan map[string]interface{}, len(roleData.PolicyNames))
	errorCh := make(chan error, len(roleData.PolicyNames))
	for _, policy := range roleData.PolicyNames {
		wg.Add(1)
		go getRolePolicyDataAsync(policy, role.RoleName, svc, &wg, policyCh, errorCh)
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

	var rolePolicies []map[string]interface{}

	for rolePolicy := range policyCh {
		rolePolicies = append(rolePolicies, rolePolicy)
	}

	return rolePolicies, nil
}

func getRolePolicyDataAsync(policy *string, roleName *string, svc *iam.IAM, wg *sync.WaitGroup, policyCh chan map[string]interface{}, errorCh chan error) {
	defer wg.Done()

	rowData, err := getRoleInlinePolicy(policy, roleName, svc)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		policyCh <- rowData
	}
}

func getRoleInlinePolicy(policyName *string, roleName *string, svc *iam.IAM) (map[string]interface{}, error) {
	rolePolicy := make(map[string]interface{})
	params := &iam.GetRolePolicyInput{
		PolicyName: policyName,
		RoleName:   roleName,
	}

	tmpPolicy, err := svc.GetRolePolicy(params)
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

		rolePolicy = map[string]interface{}{
			"PolicyDocument": rawPolicy,
			"PolicyName":     *tmpPolicy.PolicyName,
		}
	}

	return rolePolicy, nil
}

//// TRANSFORM FUNCTIONS

func getIamRoleTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*iam.Role)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
