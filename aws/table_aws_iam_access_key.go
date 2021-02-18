package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/iam"
)

func tableAwsIamAccessKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_access_key",
		Description: "AWS IAM User Access Key",
		// TODO -- get call returning a list of items

		// Get: &plugin.GetConfig{
		// 	KeyColumns:  plugin.SingleColumn("user_name"),
		// 	ItemFromKey: accessKeyFromKey,
		// 	Hydrate:     getIamAccessKey,
		// },
		List: &plugin.ListConfig{
			ParentHydrate: listIamUsers,
			Hydrate:       listUserAccessKeys,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "access_key_id",
				Description: "The ID for this access key",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_name",
				Description: "The name of the IAM user that the key is associated with",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the access key. Active means that the key is valid for API calls; Inactive means it is not",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The date when the access key was created",
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

//// BUILD HYDRATE INPUT

func accessKeyFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	userName := quals["user_name"].GetStringValue()
	item := &iam.AccessKeyMetadata{
		UserName: &userName,
	}
	return item, nil
}

//// LIST FUNCTION

func listUserAccessKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listUserAccessKeys")
	user := h.Item.(*iam.User)

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &iam.ListAccessKeysInput{
		UserName: user.UserName,
	}

	// List IAM user access keys
	item, err := svc.ListAccessKeys(params)
	if err != nil {
		return nil, err
	}

	for _, key := range item.AccessKeyMetadata {
		d.StreamLeafListItem(ctx, key)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamAccessKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamAccessKey")
	key := h.Item.(*iam.AccessKeyMetadata)

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &iam.ListAccessKeysInput{
		UserName: key.UserName,
	}

	item, err := svc.ListAccessKeys(params)

	// return results as interface
	var accessKeyRowData []interface{}
	for _, item := range item.AccessKeyMetadata {
		accessKeyRowData = append(accessKeyRowData, item)
	}

	return accessKeyRowData, nil
}

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
