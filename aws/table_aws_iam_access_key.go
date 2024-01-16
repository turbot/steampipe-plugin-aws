package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsIamAccessKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_access_key",
		Description: "AWS IAM User Access Key",
		List: &plugin.ListConfig{
			ParentHydrate: listIamUsers,
			Hydrate:       listUserAccessKeys,
			Tags:          map[string]string{"service": "iam", "action": "ListAccessKeys"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_name", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getIamAccessKeyLastUsed,
				Tags: map[string]string{"service": "iam", "action": "GetAccessKeyLastUsed"},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "access_key_id",
				Description: "The ID for this access key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_name",
				Description: "The name of the IAM user that the key is associated with.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the access key. Active means that the key is valid for API calls; Inactive means it is not.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The date when the access key was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "access_key_last_used_date",
				Description: "The date when the access key was last used.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getIamAccessKeyLastUsed,
				Transform:   transform.FromField("LastUsedDate"),
			},
			{
				Name:        "access_key_last_used_service",
				Description: "The service last used by the access key.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIamAccessKeyLastUsed,
				Transform:   transform.FromField("ServiceName"),
			},
			{
				Name:        "access_key_last_used_region",
				Description: "The region in which the access key was last used.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIamAccessKeyLastUsed,
				Transform:   transform.FromField("Region"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccessKeyId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamAccessKeyAka,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listUserAccessKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(types.User)

	// Minimize the API call with the given user_name
	equalQuals := d.EqualsQuals
	if equalQuals["user_name"] != nil {
		if equalQuals["user_name"].GetStringValue() != "" {
			if equalQuals["user_name"].GetStringValue() != "" && equalQuals["user_name"].GetStringValue() != *user.UserName {
				return nil, nil
			}
		} else if len(getListValues(equalQuals["user_name"].GetListValue())) > 0 {
			if !strings.Contains(fmt.Sprint(getListValues(equalQuals["user_name"].GetListValue())), *user.UserName) {
				return nil, nil
			}
		}
	}

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_access_key.listUserAccessKeys", "client_error", err)
		return nil, err
	}

	params := &iam.ListAccessKeysInput{UserName: user.UserName}

	paginator := iam.NewListAccessKeysPaginator(svc, params, func(o *iam.ListAccessKeysPaginatorOptions) {
		o.Limit = 10
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "NoSuchEntity" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_iam_access_key.listUserAccessKeys", "api_error", err)
			return nil, err
		}

		for _, key := range output.AccessKeyMetadata {
			d.StreamListItem(ctx, key)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamAccessKeyAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accessKey := h.Item.(types.AccessKeyMetadata)

	commonColumnData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	awsCommonData := commonColumnData.(*awsCommonColumnData)

	aka := []string{"arn:" + awsCommonData.Partition + ":iam::" + awsCommonData.AccountId + ":user/" + *accessKey.UserName + "/accesskey/" + *accessKey.AccessKeyId}
	return aka, nil
}

func getIamAccessKeyLastUsed(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_access_key.getIamAccessKeyLastUsed", "client_error", err)
		return nil, err
	}

	accessKey := h.Item.(types.AccessKeyMetadata)

	params := iam.GetAccessKeyLastUsedInput{
		AccessKeyId: accessKey.AccessKeyId,
	}

	op, err := svc.GetAccessKeyLastUsed(ctx, &params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_access_key.getIamAccessKeyLastUsed", "api_error", err)
		return nil, err
	}

	return op.AccessKeyLastUsed, nil
}
