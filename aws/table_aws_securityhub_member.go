package aws

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

// AWS SDK Migration from V1 to V2 using middleware
// due to https://github.com/aws/aws-sdk-go-v2/issues/1884

//// TABLE DEFINITION

func tableAwsSecurityHubMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_member",
		Description: "AWS Securityhub Member",
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubMembers,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInputException", "BadRequestException"}),
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
	cfg, err := SecurityHubClientConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_member.listSecurityHubMembers", "service_error", err)
		return nil, err
	}
	if cfg == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &securityhub.ListMembersInput{
		OnlyAssociated: false,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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

	myMiddleware := middleware.SerializeMiddlewareFunc(
		"AssociatedMembers",
		func(ctx context.Context, input middleware.SerializeInput, next middleware.SerializeHandler) (
			output middleware.SerializeOutput,
			metadata middleware.Metadata,
			err error) {
			req, ok := input.Request.(*smithyhttp.Request)
			if !ok {
				return output, metadata, fmt.Errorf("unexpected transport: %T", input.Request)
			}

			params, ok = input.Parameters.(*securityhub.ListMembersInput)
			if !ok {
				return output, metadata, fmt.Errorf("unexpected input type: %T", input.Parameters)
			}

			query := req.URL.Query()
			query.Set("OnlyAssociated", strconv.FormatBool(false))
			req.URL.RawQuery = query.Encode()
			return next.HandleSerialize(ctx, input)
		},
	)

	client := securityhub.NewFromConfig(*cfg, func(options *securityhub.Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			return stack.Serialize.Insert(myMiddleware, "OperationSerializer", middleware.After)
		})
	})

	paginator := securityhub.NewListMembersPaginator(client, params, func(o *securityhub.ListMembersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
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

	return nil, nil
}
