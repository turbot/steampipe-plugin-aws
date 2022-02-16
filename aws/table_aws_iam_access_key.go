package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
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
	plugin.Logger(ctx).Trace("listUserAccessKeys")
	user := h.Item.(*iam.User)

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
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &iam.ListAccessKeysInput{
		UserName: user.UserName,
		MaxItems: aws.Int64(1000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxItems {
			if *limit < 1 {
				params.MaxItems = aws.Int64(1)
			} else {
				params.MaxItems = limit
			}
		}
	}

	// List IAM user access keys
	err = svc.ListAccessKeysPages(
		params,
		func(page *iam.ListAccessKeysOutput, isLast bool) bool {
			for _, key := range page.AccessKeyMetadata {
				d.StreamListItem(ctx, key)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listUserAccessKeys", "ListAccessKeysPages", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamAccessKeyAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accessKey := h.Item.(*iam.AccessKeyMetadata)

	commonColumnData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	awsCommonData := commonColumnData.(*awsCommonColumnData)
	aka := []string{"arn:" + awsCommonData.Partition + ":iam::" + awsCommonData.AccountId + ":user/" + *accessKey.UserName + "/accesskey/" + *accessKey.AccessKeyId}
	return aka, nil
}
