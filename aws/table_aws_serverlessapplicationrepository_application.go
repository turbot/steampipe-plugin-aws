package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository"
	"github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository/types"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsServerlessApplicationRepositoryApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_serverlessapplicationrepository_application",
		Description: "AWS Serverless Application Repository Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"InvalidParameter", "NotFoundException"}),
			},
			Hydrate: getServerlessApplicationRepositoryApplication,
		},
		List: &plugin.ListConfig{
			Hydrate: listServerlessApplicationRepositoryApplications,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The application Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationId"),
			},
			{
				Name:        "author",
				Description: "The name of the author publishing the app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time this resource was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "home_page_url",
				Description: "A URL with more information about the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_verified_author",
				Description: "Whether the author is verified.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getServerlessApplicationRepositoryApplication,
				Default:     false,
			},
			{
				Name:        "license_url",
				Description: "The URL of the license.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getServerlessApplicationRepositoryApplication,
			},
			{
				Name:        "readme_url",
				Description: "The URL of the Readme.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getServerlessApplicationRepositoryApplication,
			},
			{
				Name:        "spdx_license_id",
				Description: "A valid identifier from https://spdx.org/licenses/.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "verified_author_url",
				Description: "The URL of the verified author.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getServerlessApplicationRepositoryApplication,
			},
			{
				Name:        "labels",
				Description: "Labels to improve discovery of apps in search results.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "statements",
				Description: "The contents of the access policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServerlessApplicationRepositoryApplicationPolicy,
			},
			{
				Name:        "version",
				Description: "The policy statement of the application.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServerlessApplicationRepositoryApplication,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationId").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listServerlessApplicationRepositoryApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := ServerlessApplicationRepositoryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_serverlessapplicationrepository_application.listServerlessApplicationRepositoryApplications", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
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

	// Set MaxItems to the maximum number allowed
	input := &serverlessapplicationrepository.ListApplicationsInput{
		MaxItems: maxLimit,
	}

	paginator := serverlessapplicationrepository.NewListApplicationsPaginator(svc, input, func(o *serverlessapplicationrepository.ListApplicationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_serverlessapplicationrepository_application.listServerlessApplicationRepositoryApplications", "api_error", err)
			return nil, err
		}

		for _, items := range output.Applications {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServerlessApplicationRepositoryApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		arn = *serverlessApplicationRepositoryArn(h.Item)
	} else {
		arn = d.KeyColumnQuals["arn"].GetStringValue()
	}

	if arn == "" {
		return nil, nil
	}

	// Create service
	svc, err := ServerlessApplicationRepositoryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_serverlessapplicationrepository_application.getServerlessApplicationRepositoryApplication", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &serverlessapplicationrepository.GetApplicationInput{
		ApplicationId: &arn,
	}

	// Get call
	data, err := svc.GetApplication(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_serverlessapplicationrepository_application.getServerlessApplicationRepositoryApplication", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getServerlessApplicationRepositoryApplicationPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		arn = *serverlessApplicationRepositoryArn(h.Item)
	}

	// Create service
	svc, err := ServerlessApplicationRepositoryClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_serverlessapplicationrepository_application.getServerlessApplicationRepositoryApplicationPolicy", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &serverlessapplicationrepository.GetApplicationPolicyInput{
		ApplicationId: &arn,
	}

	// Get call
	data, err := svc.GetApplicationPolicy(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ForbiddenException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_serverlessapplicationrepository_application.getServerlessApplicationRepositoryApplicationPolicy", "api_error", err)
		return nil, err
	}

	return data, nil
}

func serverlessApplicationRepositoryArn(item interface{}) *string {
	switch item := item.(type) {
	case types.ApplicationSummary:
		return item.ApplicationId
	case *serverlessapplicationrepository.GetApplicationOutput:
		return item.ApplicationId
	}
	return nil
}
