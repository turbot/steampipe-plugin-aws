package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/identitystore"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsIdentityStoreGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_identitystore_group",
		Description: "AWS Identity Store Group",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"identity_store_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getIdentityStoreGroup,
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_store_id", "name"}),
			Hydrate:    listIdentityStoreGroups,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "identity_store_id",
				Description: "The globally unique identifier for the identity store.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Contains the group’s display name value.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Group.DisplayName"),
			},
			{
				Name:        "id",
				Description: "The identifier for a group in the identity store.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Group.GroupId"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Group.DisplayName"),
			},
			// TODO: Decide if GroupId is a suitable value for `akas`. It will normally be
			// a GUID, but this is determined by the underlying identity store.
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Group.GroupId").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listIdentityStoreGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listIdentityStoreGroups")

	name := d.KeyColumnQuals["name"].GetStringValue()
	identityStoreId := d.KeyColumnQuals["identity_store_id"].GetStringValue()

	// Create session
	svc, err := IdentityStoreService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &identitystore.ListGroupsInput{
		IdentityStoreId: aws.String(identityStoreId),
		Filters: []*identitystore.Filter{
			{
				AttributePath:  aws.String("DisplayName"),
				AttributeValue: aws.String(name),
			},
		},
	}

	err = svc.ListGroupsPages(
		params,
		func(page *identitystore.ListGroupsOutput, isLast bool) bool {
			for _, group := range page.Groups {
				item := &IdentityStoreGroup{
					IdentityStoreId: &identityStoreId,
					Group:           group,
				}
				d.StreamListItem(ctx, item)
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listIdentityStoreGroups", "ListGroupsPages_error", err)

	}
	return nil, err
}

type IdentityStoreGroup struct {
	IdentityStoreId *string
	Group           *identitystore.Group
}

//// HYDRATE FUNCTIONS

func getIdentityStoreGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIdentityStoreGroup")

	groupId := d.KeyColumnQuals["id"].GetStringValue()
	identityStoreId := d.KeyColumnQuals["identity_store_id"].GetStringValue()

	// Create session
	svc, err := IdentityStoreService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &identitystore.DescribeGroupInput{
		GroupId:         aws.String(groupId),
		IdentityStoreId: aws.String(identityStoreId),
	}

	op, err := svc.DescribeGroup(params)
	if err != nil {
		plugin.Logger(ctx).Error("getIdentityStoreGroup", "DescribeGroup_error", err)
		return nil, err
	}

	item := &IdentityStoreGroup{
		IdentityStoreId: &identityStoreId,
		Group: &identitystore.Group{
			DisplayName: op.DisplayName,
			GroupId:     op.GroupId,
		},
	}

	return item, nil
}
