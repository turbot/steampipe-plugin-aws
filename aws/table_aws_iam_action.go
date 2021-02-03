package aws

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

var permissionsData ParliamentPermissions

//// TABLE DEFINITION

func tableAwsIamAction(_ context.Context) *plugin.Table {
	permissionsData = getParliamentIamPermissions()

	return &plugin.Table{
		Name:        "aws_iam_action",
		Description: "AWS IAM Action",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("action"),
			Hydrate:     getIamAction,
			ItemFromKey: permissionFromKey,
		},
		List: &plugin.ListConfig{
			Hydrate: listIamActions,
		},
		Columns: []*plugin.Column{
			// "Key" Columns
			{
				Name:        "action",
				Description: "The action for this permission",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "prefix",
				Type:        proto.ColumnType_STRING,
				Description: "The prefix for this action",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "privilege",
				Type:        proto.ColumnType_STRING,
				Description: "The privilege for this action",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description for this action",
				Transform:   transform.FromGo(),
			},
		},
	}
}

type awsIamPermissionData struct {
	Prefix      string
	Privilege   string
	Action      string
	Description string
}

//// ITEM FROM KEY

func permissionFromKey(_ context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	action := quals["action"].GetStringValue()
	item := &awsIamPermissionData{
		Action: action,
	}
	return item, nil
}

//// LIST FUNCTION

func listIamActions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	for _, service := range permissionsData {
		for _, privilege := range service.Privileges {
			d.StreamListItem(ctx, awsIamPermissionData{
				Prefix:      service.Prefix,
				Privilege:   privilege.Privilege,
				Action:      service.Prefix + ":" + privilege.Privilege,
				Description: privilege.Description,
			})
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamAction(ctx context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Info("Item", h.Item)
	action := h.Item.(*awsIamPermissionData)
	for _, service := range permissionsData {
		for _, privilege := range service.Privileges {
			a := service.Prefix + ":" + privilege.Privilege
			if a == action.Action {
				return awsIamPermissionData{
					Prefix:      service.Prefix,
					Privilege:   privilege.Privilege,
					Action:      a,
					Description: privilege.Description,
				}, nil
			}
		}

	}
	return nil, nil
}
