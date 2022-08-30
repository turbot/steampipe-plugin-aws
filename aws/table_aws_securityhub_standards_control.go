package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsSecurityHubStandardsControl(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_standards_control",
		Description: "AWS Security Hub Standards Control",
		List: &plugin.ListConfig{
			ParentHydrate: listSecurityHubStandardsSubcriptions,
			Hydrate:       listSecurityHubStandardsControls,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "control_id",
				Description: "The identifier of the security standard control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the security standard control.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StandardsControlArn"),
			},
			{
				Name:        "control_status",
				Description: "The current status of the security standard control. Indicates whether the control is enabled or disabled. Security Hub does not check against disabled controls.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "severity_rating",
				Description: "The severity of findings generated from this security standard control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "control_status_updated_at",
				Description: "The date and time that the status of the security standard control was most recently updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The longer description of the security standard control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disabled_reason",
				Description: "The reason provided for the most recent change in status for the control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "remediation_url",
				Description: "A link to remediation information for the control in the Security Hub user documentation.",
				Type:        proto.ColumnType_STRING,
			},

			// JSON columns
			{
				Name:        "related_requirements",
				Description: "The list of requirements that are related to this control.",
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
				Transform:   transform.FromField("StandardsControlArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityHubStandardsControls(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSecurityHubStandardsControls")
	region := d.KeyColumnQualString(matrixKeyRegion)

	standardsArn := h.Item.(*securityhub.Standard).StandardsArn

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// Standards Subscription Arn format
	// arn:aws:securityhub:us-east-1:accountID:subscription/aws-foundational-security-best-practices/v/1.0.0
	// arn:aws:securityhub:::ruleset/cis-aws-foundations-benchmark/v/1.2.0
	var standardsSubscriptionArn string
	if strings.Contains(*standardsArn, "standards") {
		standardsSubscriptionArn = "arn:aws:securityhub:" + region + ":" + commonColumnData.AccountId + ":subscription" + strings.Split(*standardsArn, "standards")[1]
	} else {
		standardsSubscriptionArn = "arn:aws:securityhub:" + region + ":" + commonColumnData.AccountId + ":subscription" + strings.Split(*standardsArn, "ruleset")[1]
	}

	// Create session
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &securityhub.DescribeStandardsControlsInput{
		MaxResults:               aws.Int64(100),
		StandardsSubscriptionArn: &standardsSubscriptionArn,
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

	// List call
	err = svc.DescribeStandardsControlsPages(
		input,
		func(page *securityhub.DescribeStandardsControlsOutput, isLast bool) bool {
			for _, control := range page.Controls {
				d.StreamListItem(ctx, control)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		// Handle error for unsupported or inactive regions
		if strings.Contains(err.Error(), "ResourceNotFoundException") || strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listSecurityHubStandardsControls", "list", err)
	}

	return nil, err
}
