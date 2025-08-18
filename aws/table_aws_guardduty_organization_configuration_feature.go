package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type organizationFeatureInfo struct {
	types.OrganizationFeatureConfigurationResult
	MemberAccountLimitReached     *bool
	AutoEnableOrganizationMembers types.AutoEnableMembers
	DataSources                   *types.OrganizationDataSourceConfigurationsResult
	DetectorId                    string
}

//// TABLE DEFINITION

func tableAwsGuardDutyOrganizationConfigurationFeature(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_guardduty_organization_configuration_feature",
		Description: "AWS GuardDuty Organization Configuration Feature",
		List: &plugin.ListConfig{
			ParentHydrate: listGuardDutyDetectors,
			Hydrate:       listGuardDutyOrganizationConfigurationFeatures,
			Tags:          map[string]string{"service": "guardduty", "action": "DescribeOrganizationConfiguration"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "detector_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_GUARDDUTY_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "detector_id",
				Description: "The ID of the detector.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DetectorId"),
			},
			{
				Name:        "name",
				Description: "The name of the feature that is configured for the member accounts within the organization.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "member_account_limit_reached",
				Description: "Indicates whether the maximum number of allowed member accounts are already associated with the delegated administrator account for your organization.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "auto_enable",
				Description: "Describes the status of the feature that is configured for the member accounts within the organization.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_enable_organization_members",
				Description: "Indicates the auto-enablement configuration of GuardDuty or any of the corresponding protection plans for the member accounts in the organization.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_sources",
				Description: "[DEPRECATED] Describes which data sources are enabled automatically for member accounts.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "additional_configuration",
				Description: "The additional configuration that is configured for the member accounts within the organization.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

func listGuardDutyOrganizationConfigurationFeatures(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get detector details
	detector := h.Item.(detectorInfo)

	// Check if we should filter by detector_id
	equalQuals := d.EqualsQuals
	if equalQuals["detector_id"] != nil {
		if equalQuals["detector_id"].GetStringValue() != detector.DetectorID {
			return nil, nil
		}
	}

	// Create service
	svc, err := GuardDutyClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_guardduty_organization_configuration_feature.listGuardDutyOrganizationConfigurationFeatures", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	plugin.Logger(ctx).Error("Detector ID ====>>>", "detector_id", detector.DetectorID)

	// Build the params
	maxItems := int32(50)
	params := &guardduty.DescribeOrganizationConfigurationInput{
		DetectorId: aws.String(detector.DetectorID),
		MaxResults: aws.Int32(maxItems),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			params.MaxResults = aws.Int32(limit)
		}
	}

	// Create paginator
	paginator := guardduty.NewDescribeOrganizationConfigurationPaginator(svc, params, func(o *guardduty.DescribeOrganizationConfigurationPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// We are encountering the Error: aws: operation error GuardDuty: DescribeOrganizationConfiguration, https response error StatusCode: 400, RequestID: 2f0d9365-28c2-4a86-82cc-926d2d670a9a, BadRequestException: The request is rejected because an invalid or out-of-range value is specified as an input parameter.
			// In which region the Detector is not available
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) {
				if awsErr.ErrorCode() == "BadRequestException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_guardduty_organization_configuration_feature.listGuardDutyOrganizationConfigurationFeatures", "api_error", err)
			return nil, err
		}

		if output != nil && output.Features != nil {
			for _, feature := range output.Features {
				d.StreamListItem(ctx, organizationFeatureInfo{
					MemberAccountLimitReached:              output.MemberAccountLimitReached,
					AutoEnableOrganizationMembers:          output.AutoEnableOrganizationMembers,
					DataSources:                            output.DataSources,
					OrganizationFeatureConfigurationResult: feature,
					DetectorId:                             detector.DetectorID,
				})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}
