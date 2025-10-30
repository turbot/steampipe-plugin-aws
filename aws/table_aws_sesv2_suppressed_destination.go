package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSesV2SuppressedDestination(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sesv2_suppressed_destination",
		Description: "AWS SESv2 Suppressed Destination",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("email_address"),
			Hydrate:    getSuppressedDestination,
			Tags:       map[string]string{"service": "sesv2", "action": "GetSuppressedDestination"},
		},
		List: &plugin.ListConfig{
			Hydrate: listSesSuppressedDestinations,
			Tags:    map[string]string{"service": "sesv2", "action": "ListSuppressedDestinations"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "reason",
					Require: plugin.Optional,
				},
				{
					Name:    "start_date",
					Require: plugin.Optional,
					Operators: []string{"="},
				},
				{
					Name:    "end_date",
					Require: plugin.Optional,
					Operators: []string{"="},
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getSuppressedDestination,
				Tags: map[string]string{"service": "sesv2", "action": "GetSuppressedDestination"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EMAIL_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "email_address",
				Description: "The email address that is on the suppression list for your account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reason",
				Description: "The reason that the address was added to the suppression list for your account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_update_time",
				Description: "The date and time when the suppressed destination was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "start_date",
				Description: "Used to filter the list of suppressed email addresses so that it only includes addresses that were added to the list after the specified date.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("start_date"),
			},
			{
				Name:        "end_date",
				Description: "Used to filter the list of suppressed email addresses so that it only includes addresses that were added to the list before the specified date.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("end_date"),
			},

			// Hydrated columns from GetSuppressedDestination
			{
				Name:        "suppressed_destination_attributes",
				Description: "An object that contains additional attributes that are related to a suppressed destination.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSuppressedDestination,
			},
			{
				Name:        "message_tag",
				Description: "A unique identifier that's generated when an email address is added to the suppression list for your account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSuppressedDestination,
				Transform:   transform.FromField("Attributes.MessageTag"),
			},
			{
				Name:        "feedback_id",
				Description: "The unique identifier of the email message that caused the email address to be added to the suppression list for your account.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSuppressedDestination,
				Transform:   transform.FromField("Attributes.FeedbackId"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("EmailAddress"),
			},
		}),
	}
}

//// LIST FUNCTION
func listSesSuppressedDestinations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := SESV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sesv2_suppressed_destination.listSesSuppressedDestinations", "client_error", err)
		return nil, err
	}

	maxItems := int32(1000)
	params := &sesv2.ListSuppressedDestinationsInput{}

	// Handle optional qualifiers
	if d.EqualsQuals["reason"] != nil {
		reasonStr := d.EqualsQualString("reason")
		if reasonStr != "" {
			// Convert string to SuppressionListReason enum
			switch reasonStr {
			case "BOUNCE":
				params.Reasons = []types.SuppressionListReason{types.SuppressionListReasonBounce}
			case "COMPLAINT":
				params.Reasons = []types.SuppressionListReason{types.SuppressionListReasonComplaint}
			default:
				plugin.Logger(ctx).Warn("aws_sesv2_suppressed_destination.listSesSuppressedDestinations", "invalid_reason", reasonStr)
			}
		}
	}

	if d.EqualsQuals["start_date"] != nil {
		if d.EqualsQuals["start_date"].GetTimestampValue() != nil {
			startDate := d.EqualsQuals["start_date"].GetTimestampValue().AsTime()
			params.StartDate = &startDate
		}
	}

	if d.EqualsQuals["end_date"] != nil {
		if d.EqualsQuals["end_date"].GetTimestampValue() != nil {
			endDate := d.EqualsQuals["end_date"].GetTimestampValue().AsTime()
			params.EndDate = &endDate
		}
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			params.PageSize = &limit
		}
	}

	paginator := sesv2.NewListSuppressedDestinationsPaginator(svc, params, func(o *sesv2.ListSuppressedDestinationsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sesv2_suppressed_destination.listSesSuppressedDestinations", "api_error", err)
			return nil, err
		}

		for _, item := range output.SuppressedDestinationSummaries {
			d.StreamListItem(ctx, item)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION
func getSuppressedDestination(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	svc, err := SESV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sesv2_suppressed_destination.getSuppressedDestination", "client_error", err)
		return nil, err
	}

	var emailAddress string
	if h.Item != nil {
		// Called from list - extract email address from the list item
		if item, ok := h.Item.(types.SuppressedDestinationSummary); ok {
			if item.EmailAddress != nil {
				emailAddress = *item.EmailAddress
			}
		}
	} else {
		// Called from get - extract email address from key columns
		emailAddress = d.EqualsQualString("email_address")
	}

	if emailAddress == "" {
		plugin.Logger(ctx).Error("aws_sesv2_suppressed_destination.getSuppressedDestination", "missing_email_address")
		return nil, nil
	}

	params := &sesv2.GetSuppressedDestinationInput{
		EmailAddress: &emailAddress,
	}

	output, err := svc.GetSuppressedDestination(ctx, params)
	if err != nil {
		var notFoundErr *types.NotFoundException
		if errors.As(err, &notFoundErr) {
			plugin.Logger(ctx).Debug("aws_sesv2_suppressed_destination.getSuppressedDestination", "email_not_found", emailAddress)
		}

		plugin.Logger(ctx).Error("aws_sesv2_suppressed_destination.getSuppressedDestination", "api_error", err)
		return nil, err
	}

	return output.SuppressedDestination, nil
}
