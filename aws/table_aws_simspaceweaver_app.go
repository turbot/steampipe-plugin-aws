package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/simspaceweaver"
	"github.com/aws/aws-sdk-go-v2/service/simspaceweaver/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSimSpaceWeaverApp(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sim_space_weaver_app",
		Description: "AWS SimSpace Weaver App",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"simulation", "name", "domain"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsSimSpaceWeaverApp,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsSimSpaceWeaverSimulations,
			Hydrate:       listAwsSimSpaceWeaverApps,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "simulation",
				Description: "The name of the simulation of the app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain",
				Description: "The domain of the app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_status",
				Description: "The desired status of the app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint_info",
				Description: "Information about the network endpoint for the custom app.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSimSpaceWeaverApp,
			},
			{
				Name:        "launch_overrides",
				Description: "Options that apply when the app starts. These options override default behavior.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSimSpaceWeaverApp,
			},

			// Steampipe standard coulumns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(arnToTitle),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSimSpaceWeaverApps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	simulation := h.Item.(types.SimulationMetadata)

	// Create  Client
	svc, err := SimSpaceWeaverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sim_space_weaver_app.listAwsSimSpaceWeaverApps", "connection_error", err)
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

	params := &simspaceweaver.ListAppsInput{
		Simulation: simulation.Name,
		MaxResults: &maxLimit,
	}
	// Does not support limit
	paginator := simspaceweaver.NewListAppsPaginator(svc, params, func(o *simspaceweaver.ListAppsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Simulation must be STARTED to do this otherwise we will get ConflictException exception.
			// This could not be caught by ignore config
			if strings.Contains(err.Error(), "ConflictException") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_sim_space_weaver_app.listAwsSimSpaceWeaverApps", "api_error", err)
			return nil, err
		}
		for _, app := range output.Apps {
			d.StreamListItem(ctx, app)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSimSpaceWeaverApp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var simulationName, name, domain string
	if h.Item != nil {
		data := h.Item.(types.SimulationAppMetadata)
		name = *data.Name
		simulationName = *data.Simulation
		domain = *data.Domain
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		simulationName = d.KeyColumnQuals["simulation"].GetStringValue()
		domain = d.KeyColumnQuals["domain"].GetStringValue()
	}

	// Empty Check
	if name == "" || domain == "" || simulationName == "" {
		return nil, nil
	}

	// Create session
	svc, err := SimSpaceWeaverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sim_space_weaver_app.getAwsSimSpaceWeaverApp", "connection_err", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &simspaceweaver.DescribeAppInput{
		Simulation: aws.String(simulationName),
		Domain:     &domain,
		App:        &name,
	}

	op, err := svc.DescribeApp(ctx, params)
	if err != nil {
		// Simulation must be STARTED to do this otherwise we will get ConflictException exception.
		// This could not be caught by ignore config
		if strings.Contains(err.Error(), "ConflictException") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_sim_space_weaver_app.getAwsSimSpaceWeaverApp", "api_error", err)
		return nil, err
	}
	return op, nil
}
