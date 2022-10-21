package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/serverlessapplicationrepository"
	"github.com/turbot/go-kit/types"
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
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameter", "NotFoundException"}),
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
	logger := plugin.Logger(ctx)
	logger.Trace("listServerlessApplicationRepositoryApplications")

	// Create service
	svc, err := ServerlessApplicationRepositoryService(ctx, d)
	if err != nil {
		logger.Error("listServerlessApplicationRepositoryApplications", "error_ServerlessApplicationRepositoryService", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Set MaxItems to the maximum number allowed
	input := serverlessapplicationrepository.ListApplicationsInput{
		MaxItems: types.Int64(100),
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxItems {
			if *limit < 1 {
				input.MaxItems = types.Int64(1)
			} else {
				input.MaxItems = limit
			}
		}
	}

	err = svc.ListApplicationsPages(
		&input,
		func(page *serverlessapplicationrepository.ListApplicationsOutput, lastPage bool) bool {
			for _, application := range page.Applications {
				d.StreamListItem(ctx, application)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	if err != nil {
		logger.Error("listServerlessApplicationRepositoryApplications", "error_ListApplicationsOutput", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServerlessApplicationRepositoryApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getServerlessApplicationRepositoryApplication")

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
	svc, err := ServerlessApplicationRepositoryService(ctx, d)
	if err != nil {
		logger.Error("getServerlessApplicationRepositoryApplication", "error_ServerlessApplicationRepositoryService", err)
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
	data, err := svc.GetApplication(params)
	if err != nil {
		logger.Error("getServerlessApplicationRepositoryApplication", "error_GetApplication", err)
		return nil, err
	}

	return data, nil
}

func getServerlessApplicationRepositoryApplicationPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getServerlessApplicationRepositoryApplicationPolicy")

	var arn string
	if h.Item != nil {
		arn = *serverlessApplicationRepositoryArn(h.Item)
	}

	// Create service
	svc, err := ServerlessApplicationRepositoryService(ctx, d)
	if err != nil {
		logger.Error("getServerlessApplicationRepositoryApplicationPolicy", "error_ServerlessApplicationRepositoryService", err)
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
	data, err := svc.GetApplicationPolicy(params)
	if err != nil {
		logger.Error("getServerlessApplicationRepositoryApplicationPolicy", "error_GetApplicationPolicy", err)
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ForbiddenException" {
				return nil, nil
			}
		}
		return nil, err
	}

	return data, nil
}

func serverlessApplicationRepositoryArn(item interface{}) *string {
	switch item := item.(type) {
	case *serverlessapplicationrepository.ApplicationSummary:
		return item.ApplicationId
	case *serverlessapplicationrepository.GetApplicationOutput:
		return item.ApplicationId
	}
	return nil
}
