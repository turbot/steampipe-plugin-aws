package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func tableAwsIamAccessKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_access_key",
		Description: "AWS IAM User Access Key",
		List: &plugin.ListConfig{
			ParentHydrate: listIamUsers,
			Hydrate:       listUserAccessKeys,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_name", Require: plugin.Optional},
			},
		},
		Columns: awsColumns([]*plugin.Column{
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
	equalQuals := d.KeyColumnQuals
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
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_access_key.listUserAccessKeys", "api_error", err)
			return nil, err
		}

		for _, key := range output.AccessKeyMetadata {
			d.StreamListItem(ctx, key)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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
