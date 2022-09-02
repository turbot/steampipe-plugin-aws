package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/guardduty"
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
	guardduty.Member
	DetectorId string
}

//// LIST FUNCTION

func listGuardDutyMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listGuardDutyMembers")
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
	svc, err := GuardDutyService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &guardduty.ListMembersInput{
		MaxResults:     aws.Int64(50),
		DetectorId:     aws.String(detectorId),
		OnlyAssociated: aws.String("false"),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	err = svc.ListMembersPages(
		input,
		func(page *guardduty.ListMembersOutput, isLast bool) bool {
			for _, member := range page.Members {
				d.StreamListItem(ctx, memberInfo{*member, detectorId})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("listGuardDutyMembers", "get", err)
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGuardDutyMember(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGuardDutyMember")

	detectorId := d.KeyColumnQuals["detector_id"].GetStringValue()
	accountId := d.KeyColumnQuals["member_account_id"].GetStringValue()

	// check if detectorId or accountId is empty
	if detectorId == "" || accountId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GuardDutyService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &guardduty.GetMembersInput{
		DetectorId: &detectorId,
		AccountIds: aws.StringSlice([]string{accountId}),
	}

	op, err := svc.GetMembers(params)
	if err != nil {
		plugin.Logger(ctx).Error("getGuardDutyMember", "get", err)
		return nil, err
	}

	if len(op.Members) > 0 {
		return memberInfo{*op.Members[0], detectorId}, nil
	}

	return nil, nil
}
