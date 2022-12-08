package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codecatalyst"
	"github.com/aws/aws-sdk-go-v2/service/codecatalyst/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeCatalystSpace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codecatalyst_space",
		Description: "AWS Code Catalyst Space",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceDoesNotExistException"}),
			},
			Hydrate: getCodeCatalystSpace,
		},
		List: &plugin.ListConfig{
			Hydrate: listCodeCatalystSpaces,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCodeDeployApplication,
			},
			{
				Name:        "region_name",
				Description: "The Amazon Web Services Region where the space exists.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the space.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The friendly name of the space displayed to users..",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			// {
			// 	Name:        "akas",
			// 	Description: resourceInterfaceDescription("akas"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     getCodeDeployApplicationArn,
			// 	Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			// },
		}),
	}
}

//// LIST FUNCTION

func listCodeCatalystSpaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)

	plugin.Logger(ctx).Error("====>>>>> ", "11111111")

	// Create session
	svc, err := CodeCatalystClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codecatalyst_project.listCodeCatalystSpaces", "service_creation_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Error("====>>>>> ", "22222222")
	if svc == nil {
		return nil, nil
	}
	plugin.Logger(ctx).Error("====>>>>> ", "33333333")
	input := &codecatalyst.ListSpacesInput{}
	plugin.Logger(ctx).Error("====>>>>> ", "input", input)


	paginator := codecatalyst.NewListSpacesPaginator(svc, input)

	// plugin.Logger(ctx).Error("====>>>>> ", paginator.HasMorePages())
	plugin.Logger(ctx).Error("====>>>>> ", "44444444")

	for paginator.HasMorePages() {
		output, _ := paginator.NextPage(ctx)
		plugin.Logger(ctx).Error("====>>>>> ", "OUTPUT", output)

		// if err != nil {
		// 	plugin.Logger(ctx).Error("aws_codecatalyst_project.listCodeCatalystSpaces", "api_error", err)
		// 	return nil, err
		// }
		// plugin.Logger(ctx).Error("====>>>>> ", "66666666")
		// if output != nil {
		// 	for _, space := range output.Items {
		// 		d.StreamListItem(ctx, space)

		// 		// Context may get cancelled due to manual cancellation or if the limit has been reached
		// 		if d.QueryStatus.RowsRemaining(ctx) == 0 {
		// 			return nil, nil
		// 		}
		// 	}
		// }
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCodeCatalystSpace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)

	var name string
	if h.Item != nil {
		name = *h.Item.(types.SpaceSummary).Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &codecatalyst.GetSpaceInput{
		Name: &name,
	}

	// Create session
	svc, err := CodeCatalystClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codecatalyst_project.getCodeCatalystSpace", "service_creation_error", err)
		return nil, err
	}

	// Get call
	data, err := svc.GetSpace(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_codecatalyst_project.getCodeCatalystSpace", "api_error", err)
		return nil, err
	}
	return data, nil
}
