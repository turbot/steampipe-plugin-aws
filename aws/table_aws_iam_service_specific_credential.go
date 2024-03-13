package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamUserServiceSpecificCredential(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_service_specific_credential",
		Description: "AWS IAM User Service Specific Credential",
		List: &plugin.ListConfig{
			ParentHydrate: listIamUsers,
			Hydrate:       listAwsIamUserServiceSpecificCredentials,
			Tags:          map[string]string{"service": "iam", "action": "ListServiceSpecificCredentials"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service_name", Require: plugin.Optional},
				{Name: "user_name", Require: plugin.Optional},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "service_name",
				Description: "The name of the service associated with the service-specific credential.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_specific_credential_id",
				Description: "The unique identifier for the service-specific credential.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The date and time, in ISO 8601 date-time format (http://www.iso.org/iso/iso8601), when the service-specific credential were created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "service_user_name",
				Description: "The generated user name for the service-specific credential.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the service-specific credential. Active means that the key is valid for API calls, while Inactive means it is not.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_name",
				Description: "The name of the IAM user associated with the service-specific credential.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsIamUserServiceSpecificCredentials(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(types.User)

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_service_specific_credential.listAwsIamUserServiceSpecificCredentials", "client_error", err)
		return nil, err
	}

	if d.EqualsQuals["user_name"].GetStringValue() != "" && *user.UserName != d.EqualsQuals["user_name"].GetStringValue() {
		return nil, nil
	}

	params := &iam.ListServiceSpecificCredentialsInput{
		UserName: user.UserName,
	}

	if d.EqualsQuals["service_name"].GetStringValue() != "" {
		params.ServiceName = aws.String(d.EqualsQuals["service_name"].GetStringValue())
	}

	userData, err := svc.ListServiceSpecificCredentials(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_service_specific_credential.listAwsIamUserServiceSpecificCredentials", "api_error", err)
		return nil, err
	}

	if userData != nil && userData.ServiceSpecificCredentials != nil {
		for _, cred := range userData.ServiceSpecificCredentials {
			d.StreamListItem(ctx, cred)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
