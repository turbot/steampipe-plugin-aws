package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITION

func tableAwsInspector2Member(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector2_member",
		Description: "AWS Inspector2 Member",
		List: &plugin.ListConfig{
			Hydrate: listInspector2Member,
			Tags:    map[string]string{"service": "inspector2", "action": "ListMembers"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "only_associated",
					Require:    plugin.Optional,
					CacheMatch: query_cache.CacheMatchExact,
				},
			},
		},

		GetMatrixItemFunc: SupportedRegionMatrix(AWS_INSPECTOR2_SERVICE_ID),

		// We *do not* use the common columns, because the account_id/region of
		// the default columns come from the call, *not* from the returned data.
		// For inspector2, the account_id or region can vary within a single
		// call.
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				// The account id from the data, rather than from the call (getCommonColumns).
				Name:        "member_account_id",
				Type:        proto.ColumnType_STRING,
				Description: "The AWS Account ID in which the resource is located.",
				Transform:   transform.FromField("AccountId"),
			},
			{
				Name:        "delegated_admin_account_id",
				Description: "The Amazon Web Services account ID of the Amazon Inspector delegated administrator for this member account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "only_associated",
				Description: "Specifies whether to list only currently associated members if True or to list all members within the organization if False.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("only_associated"),
			},
			{
				Name:        "relationship_status",
				Description: "The status of the member account. Valid values are: CREATED | INVITED | DISABLED | ENABLED | REMOVED | RESIGNED | DELETED | EMAIL_VERIFICATION_IN_PROGRESS | EMAIL_VERIFICATION_FAILED | REGION_DISABLED | ACCOUNT_SUSPENDED | CANNOT_CREATE_DETECTOR_IN_ORG_MASTER.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "A timestamp showing when the status of this member was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listInspector2Member(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := Inspector2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector2_member.listInspector2Member", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &inspector2.ListMembersInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if d.EqualsQuals["only_associated"] != nil {
		onlyAssociated := getQualsValueByColumn(d.Quals, "only_associated", "boolean")
		if onlyAssociated.(string) == "true" {
			input.OnlyAssociated = aws.Bool(true)
		} else {
			input.OnlyAssociated = aws.Bool(false)
		}
	}

	paginator := inspector2.NewListMembersPaginator(svc, input, func(o *inspector2.ListMembersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_inspector2_member.listInspector2Member", "api_error", err)
			return nil, err
		}

		for _, item := range output.Members {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}
