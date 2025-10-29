package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSsoAdminCustomerPolicyAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssoadmin_customer_policy_attachment",
		Description: "AWS SSO Customer Managed Policy Attachment",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"permission_set_arn"}),
			Hydrate:    listSsoAdminCustomerPolicyAttachments,
			Tags:       map[string]string{"service": "sso", "action": "ListCustomerManagedPolicyReferencesInPermissionSet"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SSO_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "permission_set_arn",
				Description: "The ARN of the permission set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the customer managed policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CustomerManagedPolicyReference.Name"),
			},
			{
				Name:        "instance_arn",
				Description: "The Amazon Resource Name (ARN) of the SSO Instance under which the operation will be executed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "The path to the customer managed policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CustomerManagedPolicyReference.Path"),
			},
			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CustomerManagedPolicyReference.Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listSsoAdminCustomerPolicyAttachments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	permissionSetArn := d.EqualsQuals["permission_set_arn"].GetStringValue()
	instanceArn, err := getSsoInstanceArnFromResourceArn(permissionSetArn)
	if err != nil {
		return nil, err
	}

	// Create session
	svc, err := SSOAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_customer_policy_attachment.listSsoAdminCustomerPolicyAttachments", "connection_error", err)
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

	params := &ssoadmin.ListCustomerManagedPolicyReferencesInPermissionSetInput{
		InstanceArn:      aws.String(instanceArn),
		PermissionSetArn: aws.String(permissionSetArn),
		MaxResults:       aws.Int32(maxLimit),
	}

	paginator := ssoadmin.NewListCustomerManagedPolicyReferencesInPermissionSetPaginator(svc, params, func(o *ssoadmin.ListCustomerManagedPolicyReferencesInPermissionSetPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssoadmin_customer_policy_attachment.listSsoAdminCustomerPolicyAttachments", "api_error", err)
			return nil, err
		}

		for _, items := range output.CustomerManagedPolicyReferences {
			d.StreamListItem(ctx, &CustomerPolicyAttachment{
				InstanceArn:                    aws.String(instanceArn),
				PermissionSetArn:               aws.String(permissionSetArn),
				CustomerManagedPolicyReference: items,
			})
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

type CustomerPolicyAttachment struct {
	InstanceArn                    *string
	PermissionSetArn               *string
	CustomerManagedPolicyReference types.CustomerManagedPolicyReference
}
