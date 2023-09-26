package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice"

	directoryservicev1 "github.com/aws/aws-sdk-go/service/directoryservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDirectoryServiceLogSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_directory_service_log_subscription",
		Description: "AWS Directory Service Log Subscription",
		List: &plugin.ListConfig{
			Hydrate: listDirectoryServiceLogSubscription,
			Tags:    map[string]string{"service": "directoryservice", "action": "ListLogSubscriptions"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityDoesNotExistException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "directory_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(directoryservicev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "directory_id",
				Description: "Identifier (ID) of the directory that you want to associate with the log subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_group_name",
				Description: "The name of the log group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subscription_created_date_time",
				Description: "The date and time that the log subscription was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogGroupName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDirectoryServiceLogSubscription(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DirectoryServiceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directory_service_log_subscription.listDirectoryServiceLogSubscription", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	// Build the params
	input := &directoryservice.ListLogSubscriptionsInput{
		Limit: aws.Int32(maxLimit),
	}

	if d.EqualsQualString("directory_id") != "" {
		input.DirectoryId = aws.String(d.EqualsQualString("directory_id"))
	}

	pagesLeft := true

	// List call
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		result, err := svc.ListLogSubscriptions(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_directory_service_log_subscription.listDirectoryServiceLogSubscription", "api_error", err)
			return nil, err
		}

		for _, subscription := range result.LogSubscriptions {
			d.StreamListItem(ctx, subscription)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				pagesLeft = false
			}
		}

		if result.NextToken != nil {
			input.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}
