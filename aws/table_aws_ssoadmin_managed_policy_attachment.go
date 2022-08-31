package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsSsoAdminManagedPolicyAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssoadmin_managed_policy_attachment",
		Description: "AWS SSO Managed Policy Attachment",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"permission_set_arn"}),
			Hydrate:    listSsoAdminManagedPolicyAttachments,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listSsoAdminManagedPolicyAttachments")

	permissionSetArn := d.KeyColumnQuals["permission_set_arn"].GetStringValue()
	instanceArn, err := getSsoInstanceArnFromResourceArn(permissionSetArn)
	if err != nil {
		return nil, err
	}

	// Create session
	svc, err := SSOAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &ssoadmin.ListManagedPoliciesInPermissionSetInput{
		InstanceArn:      aws.String(instanceArn),
		PermissionSetArn: aws.String(permissionSetArn),
		MaxResults:       aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			if *limit < 1 {
				params.MaxResults = aws.Int64(1)
			} else {
				params.MaxResults = limit
			}
		}
	}

	plugin.Logger(ctx).Trace("listSsoAdminManagedPolicyAttachments:ListManagedPoliciesInPermissionSetInput", "params", params)
	err = svc.ListManagedPoliciesInPermissionSetPages(params,
		func(page *ssoadmin.ListManagedPoliciesInPermissionSetOutput, isLast bool) bool {
			for _, attachedManagedPolicy := range page.AttachedManagedPolicies {
				item := &ManagedPolicyAttachment{
					InstanceArn:           &instanceArn,
					PermissionSetArn:      &permissionSetArn,
					AttachedManagedPolicy: attachedManagedPolicy,
				}
				d.StreamListItem(ctx, item)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	plugin.Logger(ctx).Trace("listSsoAdminManagedPolicyAttachments:return", "err", err)
	return nil, err
}

type ManagedPolicyAttachment struct {
	InstanceArn           *string
	PermissionSetArn      *string
	AttachedManagedPolicy *ssoadmin.AttachedManagedPolicy
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
