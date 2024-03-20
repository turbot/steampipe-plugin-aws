package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"

	ssoadminv1 "github.com/aws/aws-sdk-go/service/ssoadmin"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSsoAdminAccountAssignment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssoadmin_account_assignment",
		Description: "AWS SSO Account Assignment",
		List: &plugin.ListConfig{
			KeyColumns: append(
				plugin.AllColumns([]string{"permission_set_arn", "target_account_id"}),
				plugin.OptionalColumns([]string{"instance_arn"})...,
			),
			Hydrate: listSsoAdminAccountAssignments,
			Tags:    map[string]string{"service": "sso", "action": "ListAccountAssignments"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssoadminv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "target_account_id",
				Description: "The identifier of the AWS account from which to list the assignments.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountAssignment.AccountId"),
			},
			{
				Name:        "instance_arn",
				Description: "The Amazon Resource Name (ARN) of the SSO Instance under which the operation will be executed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceArn"),
			},
			{
				Name:        "permission_set_arn",
				Description: "The ARN of the permission set from which to list assignments.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountAssignment.PermissionSetArn"),
			},
			{
				Name:        "principal_type",
				Description: "The entity type for which the assignment will be created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountAssignment.PrincipalType"),
			},
			{
				Name:        "principal_id",
				Description: "An identifier for an object in IAM Identity Center, such as a user or group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountAssignment.PrincipalId"),
			},
		}),
	}
}

// // LIST FUNCTION

func listSsoAdminAccountAssignments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SSOAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error(d.Table.Name+".listSsoAdminAccountAssignments", "connection_error", err)
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
			maxLimit = limit
		}
	}

	permissionSetArn := d.EqualsQualString("permission_set_arn")
	input := &ssoadmin.ListAccountAssignmentsInput{
		AccountId:        aws.String(d.EqualsQualString("target_account_id")),
		PermissionSetArn: aws.String(permissionSetArn),
	}

	if v, ok := d.EqualsQuals["instance_arn"]; ok {
		input.InstanceArn = aws.String(v.GetStringValue())
	} else if i, err := ssoAdminPermissionSetToInstanceArn(permissionSetArn); err != nil {
		return nil, fmt.Errorf("failed to extract instance ARN from %q: %w", permissionSetArn, err)
	} else {
		input.InstanceArn = aws.String(i)
	}

	paginator := ssoadmin.NewListAccountAssignmentsPaginator(svc, input, func(o *ssoadmin.ListAccountAssignmentsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error(d.Table.Name+".listSsoAdminAccountAssignments", "api_error", err)
			return nil, err
		}

		for _, item := range output.AccountAssignments {
			d.StreamListItem(ctx, &SsoAdminAccountAssignmentItem{
				InstanceArn:       input.InstanceArn,
				AccountAssignment: item,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

type SsoAdminAccountAssignmentItem struct {
	InstanceArn       *string
	AccountAssignment types.AccountAssignment
}

func ssoAdminPermissionSetToInstanceArn(ps string) (string, error) {
	// Convert "arn:${Partition}:sso:::permissionSet/${InstanceId}/${PermissionSetId}"
	// to "arn:${Partition}:sso:::instance/${InstanceId}"
	// See https://docs.aws.amazon.com/service-authorization/latest/reference/list_awsiamidentitycentersuccessortoawssinglesign-on.html#awsiamidentitycentersuccessortoawssinglesign-on-resources-for-iam-policies
	a, err := arn.Parse(ps)
	if err != nil {
		return "", fmt.Errorf("failed to parse ARN: %w", err)
	} else if a.Service != "sso" {
		return "", fmt.Errorf("not an SSO ARN")
	}
	parts := strings.Split(a.Resource, "/")
	if len(parts) != 3 || parts[0] != "permissionSet" {
		return "", fmt.Errorf("not a permission set ARN")
	}
	return fmt.Sprintf("arn:aws:sso:::instance/%s", parts[1]), nil
}
