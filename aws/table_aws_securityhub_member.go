package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsSecurityHubMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_member",
		Description: "AWS Securityhub Member",
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubMembers,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"BadRequestException"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "member_account_id",
				Description: "The Amazon Web Services account ID of the member account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountId"),
			},
			{
				Name:        "administrator_id",
				Description: "The Amazon Web Services account ID of the Security Hub administrator account associated with this member account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email",
				Description: "The email address of the member account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "invited_at",
				Description: "A timestamp for the date and time when the invitation was sent to the member account.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "master_id",
				Description: "The Amazon Web Services account ID of the Security Hub administrator account associated with this member account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "member_status",
				Description: "The status of the relationship between the member account and its administrator account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The timestamp for the date and time when the member account was updated.",
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

func listSecurityHubMembers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_member.listSecurityHubMembers", "client_error", err)
		return nil, err
	}


// Limiting the results
	maxLimit := int32(500)
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

	input := &securityhub.ListMembersInput{
		MaxResults:     int32(maxLimit),
		OnlyAssociated: false,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	// limit := d.QueryContext.Limit
	// if d.QueryContext.Limit != nil {
	// 	if *limit < *input.MaxResults {
	// 		if *limit < 1 {
	// 			input.MaxResults = aws.Int64(1)
	// 		} else {
	// 			input.MaxResults = limit
	// 		}
	// 	}
	// }

	// Reduce the basic request limit down if the user has only requested a small number of rows

	// if d.QueryContext.Limit != nil {
	// 	limit := int32(*d.QueryContext.Limit)
	// 	if limit < input.MaxResults {
	// 		if limit < 1 {
	// 			input.MaxResults = int32(1)
	// 		} else {
	// 			input.MaxResults = int32(limit)
	// 		}
	// 	}
	// }





	// List call
// 	err = svc.ListMembersPages(
// 		input,
// 		func(page *securityhub.ListMembersOutput, isLast bool) bool {
// 			for _, member := range page.Members {
// 				d.StreamListItem(ctx, member)

// 				// Context may get cancelled due to manual cancellation or if the limit has been reached
// 				if d.QueryStatus.RowsRemaining(ctx) == 0 {
// 					return false
// 				}
// 			}
// 			return !isLast
// 		},
// 	)

// 	if err != nil {
// 		// Handle error for accounts that are not subscribed to AWS Security Hub
// 		if strings.Contains(err.Error(), "not subscribed") {
// 			return nil, nil
// 		}
// 		plugin.Logger(ctx).Error("listSecurityHubMembers", "list", err)
// 	}
// 	return nil, err
// }

	// maxItems := int32(100)
	// input.MaxResults = int32(maxItems)
	paginator := securityhub.NewListMembersPaginator(svc, input, func(o *securityhub.ListMembersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_securityhub_member.listSecurityHubMembers", "api_error", err)
			return nil, err
		}

		for _, member := range output.Members {
			d.StreamListItem(ctx, member)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	// if err != nil {
	// 	// Handle error for accounts that are not subscribed to AWS Security Hub
	// 	if strings.Contains(err.Error(), "not subscribed") {
	// 		return nil, nil
	// 	}
	// 	plugin.Logger(ctx).Error("listSecurityHubMembers", "list", err)
	// }
	return nil, nil
}