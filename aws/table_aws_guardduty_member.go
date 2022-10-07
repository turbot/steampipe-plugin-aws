package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGuardDutyMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_member",
		Description: "AWS GuardDuty Member",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"member_account_id", "detector_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInputException", "BadRequestException"}),
			},
			Hydrate: getGuardDutyMember,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listGuardDutyMembers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	equalQuals := d.KeyColumnQuals

	// Minimize the API call with the given detector_id
	if equalQuals["detector_id"] != nil {
		if equalQuals["detector_id"].GetStringValue() != "" {
			if equalQuals["detector_id"].GetStringValue() != "" && equalQuals["detector_id"].GetStringValue() != detectorId {
				return nil, nil
			}
		} else if len(getListValues(equalQuals["detector_id"].GetListValue())) > 0 {
			if !strings.Contains(fmt.Sprint(getListValues(equalQuals["detector_id"].GetListValue())), detectorId) {
				return nil, nil
			}
		}
	}

	// Create session
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_member.listGuardDutyMembers", "get_client_error", err)
		return nil, err
	}

	input := &guardduty.ListMembersInput{
		MaxResults:     int32(50),
		DetectorId:     aws.String(detectorId),
		OnlyAssociated: aws.String("false"),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(input.MaxResults) {
			input.MaxResults = int32(*limit)
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := svc.ListMembers(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_guardduty_member.listGuardDutyMembers", "api_error", err)
			return nil, err
		}

		for _, item := range response.Members {
			d.StreamListItem(ctx, memberInfo{item, detectorId})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if response.NextToken != nil {
			pagesLeft = true
			input.NextToken = response.NextToken
		} else {
			pagesLeft = false
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_member.listGuardDutyMembers", "api_error", err)
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGuardDutyMember(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	detectorId := d.KeyColumnQuals["detector_id"].GetStringValue()
	accountId := d.KeyColumnQuals["member_account_id"].GetStringValue()

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
		return memberInfo{*&op.Members[0], detectorId}, nil
	}

	return nil, nil
}
