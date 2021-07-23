package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type FindingInfo = struct {
	guardduty.Finding
	DetectorId string
}

//// TABLE DEFINITION

func tableAwsGuardDutyFinding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_finding",
		Description: "AWS GuardDuty Finding",
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listGuardDutyFindings,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The title of the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Title"),
			},
			{
				Name:        "id",
				Description: "The ID of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "detector_id",
				Description: "The ID of the detector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "severity",
				Description: "The severity of the finding.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "created_at",
				Description: "The time and date when the finding was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "confidence",
				Description: "The confidence score for the finding.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "description",
				Description: "The description of the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schema_version",
				Description: "The version of the schema used for the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The time and date when the finding was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "resource",
				Description: "Contains information about the AWS resource associated with the activity that prompted GuardDuty to generate a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service",
				Description: "Contains additional information about the generated finding.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

// listGuardDutyFindings handles both listing and get the details of the findings.
func listGuardDutyFindings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GuardDutyService(ctx, d)
	if err != nil {
		return nil, err
	}

	detectorId := h.Item.(detectorInfo).DetectorID

	var findingIds [][]*string

	// execute list call
	err = svc.ListFindingsPages(
		&guardduty.ListFindingsInput{
			DetectorId: aws.String(detectorId),
			MaxResults: aws.Int64(50),
		},
		func(page *guardduty.ListFindingsOutput, isLast bool) bool {
			if len(page.FindingIds) != 0 {
				findingIds = append(findingIds, page.FindingIds)
			}
			return !isLast
		},
	)
	if err != nil {
		return nil, err
	}

	for _, ids := range findingIds {
		param := &guardduty.GetFindingsInput{
			DetectorId: aws.String(detectorId),
			FindingIds: ids,
		}
		result, err := svc.GetFindings(param)
		if err != nil {
			return nil, err
		}

		for _, finding := range result.Findings {
			d.StreamListItem(ctx, FindingInfo{*finding, detectorId})
		}
	}

	return nil, nil
}
