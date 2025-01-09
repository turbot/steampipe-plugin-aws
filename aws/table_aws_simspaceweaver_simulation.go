package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/simspaceweaver"
	"github.com/aws/aws-sdk-go-v2/service/simspaceweaver/types"

	simspaceweaverEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSimSpaceWeaverSimulation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_simspaceweaver_simulation",
		Description: "AWS SimSpace Weaver Simulation",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsSimSpaceWeaverSimulation,
			Tags:    map[string]string{"service": "simspaceweaver", "action": "DescribeSimulation"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSimSpaceWeaverSimulations,
			Tags:    map[string]string{"service": "simspaceweaver", "action": "ListSimulations"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsSimSpaceWeaverSimulation,
				Tags: map[string]string{"service": "simspaceweaver", "action": "DescribeSimulation"},
			},
			{
				Func: listAwsSimSpaceWeaverSimulationTags,
				Tags: map[string]string{"service": "simspaceweaver", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(simspaceweaverEndpoint.SIMSPACEWEAVERServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the simulation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the simulation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time when the simulation was created, expressed as the number of seconds and milliseconds in UTC since the Unix epoch (0:0:0.000, January 1, 1970).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the simulation.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "status",
				Description: "The current status of the simulation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "execution_id",
				Description: "A universally unique identifier (UUID) for this simulation.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "maximum_duration",
				Description: "The maximum running time of the simulation, specified as a number of months (m or M), hours (h or H), or days (d or D).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the Identity and Access Management (IAM) role that the simulation assumes to perform actions.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "start_error",
				Description: "An error message that SimSpace Weaver returns only if a problem occurs when the simulation is in the STARTING state.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "schema_error",
				Description: "An error message that SimSpace Weaver returns only if there is a problem with the simulation schema.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "live_simulation_state",
				Description: "A collection of additional state information, such as domain and clock configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "logging_configuration",
				Description: "Settings that control how SimSpace Weaver handles your simulation log data.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "schema_s3_location",
				Description: "The location of the simulation schema in Amazon Simple Storage Service (Amazon S3).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "snapshot_s3_location",
				Description: "A location in Amazon Simple Storage Service (Amazon S3) where SimSpace Weaver stores simulation data, such as your app .zip files and schema file.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSimSpaceWeaverSimulation,
			},
			{
				Name:        "target_status",
				Description: "The desired status of the simulation.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard coulumns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(arnToTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSimSpaceWeaverSimulationTags,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSimSpaceWeaverSimulations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create  Client
	svc, err := SimSpaceWeaverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_simspaceweaver_simulation.listAwsSimSpaceWeaverSimulations", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// MaxResult value has not been specified in the docs.
	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	params := &simspaceweaver.ListSimulationsInput{
		MaxResults: &maxLimit,
	}
	// Does not support limit
	paginator := simspaceweaver.NewListSimulationsPaginator(svc, params, func(o *simspaceweaver.ListSimulationsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_simspaceweaver_simulation.listAwsSimSpaceWeaverSimulations", "api_error", err)
			return nil, err
		}
		for _, simulation := range output.Simulations {
			d.StreamListItem(ctx, simulation)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSimSpaceWeaverSimulation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		data := h.Item.(types.SimulationMetadata)
		name = *data.Name
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Empty Check
	if name == "" {
		return nil, nil
	}

	// Create session
	svc, err := SimSpaceWeaverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_simspaceweaver_simulation.getAwsSimSpaceWeaverSimulation", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &simspaceweaver.DescribeSimulationInput{
		Simulation: aws.String(name),
	}

	op, err := svc.DescribeSimulation(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_simspaceweaver_simulation.getAwsSimSpaceWeaverSimulation", "api_error", err)
		return nil, err
	}
	return op, nil
}

func listAwsSimSpaceWeaverSimulationTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn := getSimulationArn(h.Item)
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := SimSpaceWeaverClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_simspaceweaver_simulation.getAwsSimSpaceWeaverSimulationTags", "connection_error", err)
		return nil, err
	}

	params := &simspaceweaver.ListTagsForResourceInput{
		ResourceArn: aws.String(arn),
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_simspaceweaver_simulation.getAwsSimSpaceWeaverSimulationTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// UTILITY FUNCTION

func getSimulationArn(item interface{}) string {
	switch item := item.(type) {
	case types.SimulationMetadata:
		return *item.Arn
	case *simspaceweaver.DescribeSimulationOutput:
		return *item.Arn
	}
	return ""
}
