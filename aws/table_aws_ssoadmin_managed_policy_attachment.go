package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSsoAdminManagedPolicyAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssoadmin_managed_policy_attachment",
		Description: "AWS SSO Managed Policy Attachment",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"permission_set_arn"}),
			Hydrate:    listSsoAdminManagedPolicyAttachments,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "permission_set_arn",
				Description: "The ARN of the permission set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the IAM managed policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AttachedManagedPolicy.Name"),
			},
			{
				Name:        "instance_arn",
				Description: "The Amazon Resource Name (ARN) of the SSO Instance under which the operation will be executed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "managed_policy_arn",
				Description: "The ARN of the IAM managed policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AttachedManagedPolicy.Arn"),
			},
			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AttachedManagedPolicy.Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listSsoAdminManagedPolicyAttachments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	permissionSetArn := d.EqualsQuals["permission_set_arn"].GetStringValue()
	instanceArn, err := getSsoInstanceArnFromResourceArn(permissionSetArn)
	if err != nil {
		return nil, err
	}

	// Create session
	svc, err := SSOAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_managed_policy_attachment.listSsoAdminManagedPolicyAttachments", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	params := &ssoadmin.ListManagedPoliciesInPermissionSetInput{
		InstanceArn:      aws.String(instanceArn),
		PermissionSetArn: aws.String(permissionSetArn),
		MaxResults:       aws.Int32(maxLimit),
	}

	paginator := ssoadmin.NewListManagedPoliciesInPermissionSetPaginator(svc, params, func(o *ssoadmin.ListManagedPoliciesInPermissionSetPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssoadmin_managed_policy_attachment.listSsoAdminManagedPolicyAttachments", "api_error", err)
			return nil, err
		}

		for _, items := range output.AttachedManagedPolicies {
			d.StreamListItem(ctx, &ManagedPolicyAttachment{
				InstanceArn:           aws.String(instanceArn),
				PermissionSetArn:      aws.String(permissionSetArn),
				AttachedManagedPolicy: items,
			})
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

type ManagedPolicyAttachment struct {
	InstanceArn           *string
	PermissionSetArn      *string
	AttachedManagedPolicy types.AttachedManagedPolicy
}

//// UTILITY FUNCTIONS

func getSsoInstanceArnFromResourceArn(resourceArn string) (string, error) {
	arnParts := strings.Split(resourceArn, ":")
	if len(arnParts) < 6 {
		return "", fmt.Errorf("resourceArn must meet the format of an ARN")
	}

	resourceIdParts := strings.Split(arnParts[5], "/")
	if len(resourceIdParts) < 2 {
		return "", fmt.Errorf("resource ID part of resourceArn must contain an instance ID")
	}

	instanceResourceId := fmt.Sprintf("instance/%s", resourceIdParts[1])
	instanceArnParts := append(arnParts[0:5], instanceResourceId)
	instanceArn := strings.Join(instanceArnParts, ":")

	return instanceArn, nil
}
