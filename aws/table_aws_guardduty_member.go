package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"

	guarddutyEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGuardDutyMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_member",
		Description: "AWS GuardDuty Member",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"member_account_id", "detector_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException", "BadRequestException"}),
			},
			Hydrate: getGuardDutyMember,
			Tags:    map[string]string{"service": "guardduty", "action": "GetMembers"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listGuardDutyMembers,
			Tags:          map[string]string{"service": "guardduty", "action": "ListMembers"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(guarddutyEndpoint.AWS_GUARDDUTY_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "member_account_id",
				Description: "The ID of the member account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountId"),
			},
			{
				Name:        "detector_id",
				Description: "The detector ID of the member account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_id",
				Description: "The administrator account ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email",
				Description: "The email address of the member account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "invited_at",
				Description: "The timestamp when the invitation was sent.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "relationship_status",
				Description: "The status of the relationship between the member and the administrator.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administrator_id",
				Description: "The administrator account ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The last-updated timestamp of the member.",
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

type memberInfo = struct {
	types.Member
	DetectorId string
}

//// LIST FUNCTION

func listGuardDutyMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	detectorId := h.Item.(detectorInfo).DetectorID
	equalQuals := d.EqualsQuals

	// Minimize the API call with the given detector id
	if equalQuals["detector_id"] != nil {
		if equalQuals["detector_id"].GetStringValue() != detectorId {
			return nil, nil
		}
	}

	// Create session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_member.listGuardDutyMembers", "get_client_error", err)
		return nil, err
	}

	maxItems := int32(50)
	params := &guardduty.ListMembersInput{
		DetectorId:     aws.String(detectorId),
		OnlyAssociated: aws.String("false"),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			params.MaxResults = aws.Int32(limit)
		}
	}

	paginator := guardduty.NewListMembersPaginator(svc, params, func(o *guardduty.ListMembersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aaws_guardduty_member.listGuardDutyMembers", "api_error", err)
			return nil, err
		}

		for _, item := range output.Members {
			d.StreamListItem(ctx, memberInfo{item, detectorId})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGuardDutyMember(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	detectorId := d.EqualsQuals["detector_id"].GetStringValue()
	accountId := d.EqualsQuals["member_account_id"].GetStringValue()

	// check if detectorId or accountId is empty
	if detectorId == "" || accountId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_member.getGuardDutyMember", "get_client_error", err)
		return nil, err
	}

	params := &guardduty.GetMembersInput{
		DetectorId: &detectorId,
		AccountIds: ([]string{accountId}),
	}

	op, err := svc.GetMembers(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_member.getGuardDutyMember", "api_error", err)
		return nil, err
	}

	if len(op.Members) > 0 {
		return memberInfo{op.Members[0], detectorId}, nil
	}

	return nil, nil
}
