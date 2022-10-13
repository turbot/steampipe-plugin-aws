package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsIdentityStoreGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_identitystore_group",
		Description: "AWS Identity Store Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_store_id", "id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getIdentityStoreGroup,
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"identity_store_id"}),
			Hydrate:    listIdentityStoreGroups,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "identity_store_id",
				Description: "The globally unique identifier for the identity store.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Contains the group's display name value.",
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

	identityStoreId := d.KeyColumnQuals["identity_store_id"].GetStringValue()

	// Create Session
	svc, err := IdentityStoreClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_identitystore_group.listIdentityStoreGroups", "get_client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(50)
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

	params := &identitystore.ListGroupsInput{
		IdentityStoreId: aws.String(identityStoreId),
		MaxResults:      aws.Int32(maxLimit),
	}
	paginator := identitystore.NewListGroupsPaginator(svc, params, func(o *identitystore.ListGroupsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_identitystore_group.listIdentityStoreGroups", "api_error", err)
			return nil, err
		}
		for _, group := range output.Groups {
			item := &IdentityStoreGroup{
				IdentityStoreId: &identityStoreId,
				Group:           group,
			}
			d.StreamListItem(ctx, item)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

type IdentityStoreGroup struct {
	IdentityStoreId *string
	Group           types.Group
}

//// HYDRATE FUNCTIONS

func getIdentityStoreGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	groupId := d.KeyColumnQuals["id"].GetStringValue()
	identityStoreId := d.KeyColumnQuals["identity_store_id"].GetStringValue()

	// Create Session
	svc, err := IdentityStoreClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_identitystore_group.getIdentityStoreGroup", "get_client_error", err)
		return nil, err
	}

	params := &identitystore.DescribeGroupInput{
		GroupId:         aws.String(groupId),
		IdentityStoreId: aws.String(identityStoreId),
	}

	op, err := svc.DescribeGroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_identitystore_group.getIdentityStoreGroup", "api_error", err)
		return nil, err
	}

	item := &IdentityStoreGroup{
		IdentityStoreId: &identityStoreId,
		Group: types.Group{
			DisplayName: op.DisplayName,
			GroupId:     op.GroupId,
		},
	}

	return item, nil
}
