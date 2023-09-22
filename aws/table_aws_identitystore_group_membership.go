package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"

	identitystorev1 "github.com/aws/aws-sdk-go/service/identitystore"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsIdentityStoreGroupMembership(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_identitystore_group_membership",
		Description: "AWS Identity Store Group Membership",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "identity_store_id",
					Require: plugin.Required,
				},
				{
					Name:    "group_id",
					Require: plugin.Optional,
				},
			},
			ParentHydrate: listIdentityStoreGroups,
			Hydrate:       listIdentityStoreGroupMemberships,
			Tags:          map[string]string{"service": "identitystore", "action": "ListGroupMemberships"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(identitystorev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "membership_id",
				Description: "The identifier for a GroupMembership object in an identity store.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_store_id",
				Description: "The globally unique identifier for the identity store.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_id",
				Description: "The identifier for a group in the identity store.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "member_id",
				Description: "Specific identifier for a user indicates that the user is a member of the group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MemberId.Value"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MembershipId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listIdentityStoreGroupMemberships(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	groupId := d.EqualsQuals["group_id"].GetStringValue()

	group := h.Item.(*IdentityStoreGroup)

	if groupId != "" && groupId != *group.Group.GroupId {
		return nil, nil
	}

	// Create Session
	svc, err := IdentityStoreClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_identitystore_group_membership.listIdentityStoreGroupMemberships", "get_client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	params := &identitystore.ListGroupMembershipsInput{
		IdentityStoreId: group.IdentityStoreId,
		GroupId:         group.Group.GroupId,
		MaxResults:      aws.Int32(maxLimit),
	}

	paginator := identitystore.NewListGroupMembershipsPaginator(svc, params, func(o *identitystore.ListGroupMembershipsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_identitystore_group_membership.listIdentityStoreGroupMemberships", "api_error", err)
			return nil, err
		}
		for _, item := range output.GroupMemberships {
			d.StreamListItem(ctx, item)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
