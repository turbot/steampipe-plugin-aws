package aws

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

var permissionsData ParliamentPermissions

//// TABLE DEFINITION

func tableAwsIamAction(_ context.Context) *plugin.Table {
	permissionsData = getParliamentIamPermissions()

	return &plugin.Table{
		Name:        "aws_iam_action",
		Description: "AWS IAM Action",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("action"),
			Hydrate:    getIamAction,
		},
		List: &plugin.ListConfig{
			Hydrate: listIamActions,
		},
		Columns: []*plugin.Column{
			// "Key" Columns
			{
				Name:        "action",
				Description: "The action for this permission.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "prefix",
				Type:        proto.ColumnType_STRING,
				Description: "The prefix for this action.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "privilege",
				Type:        proto.ColumnType_STRING,
				Description: "The privilege for this action.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "access_level",
				Type:        proto.ColumnType_STRING,
				Description: "The access level for this action.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description for this action.",
				Transform:   transform.FromGo(),
			},
		},
	}
}

type awsIamPermissionData struct {
	Action      string
	Prefix      string
	Privilege   string
	AccessLevel string
	Description string
}

//// LIST FUNCTION

func listIamActions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	for _, service := range permissionsData {
		for _, privilege := range service.Privileges {
			d.StreamListItem(ctx, awsIamPermissionData{
				AccessLevel: privilege.AccessLevel,
				Action:      strings.ToLower(service.Prefix + ":" + privilege.Privilege),
				Description: privilege.Description,
				Prefix:      service.Prefix,
				Privilege:   privilege.Privilege,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamAction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	action := d.EqualsQuals["action"].GetStringValue()

	for _, service := range permissionsData {
		for _, privilege := range service.Privileges {
			a := strings.ToLower(service.Prefix + ":" + privilege.Privilege)
			if a == strings.ToLower(action) {
				return awsIamPermissionData{
					AccessLevel: privilege.AccessLevel,
					Action:      a,
					Description: privilege.Description,
					Prefix:      service.Prefix,
					Privilege:   privilege.Privilege,
				}, nil
			}
		}
	}
	return nil, nil
}
