package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"

	gluev1 "github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_connection",
		Description: "AWS Glue Connection",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlueConnection,
			Tags:    map[string]string{"service": "glue", "action": "GetConnection"},
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"connection_type"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
			},
			Hydrate: listGlueConnections,
			Tags:    map[string]string{"service": "glue", "action": "GetConnections"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(gluev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the connection definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the connection.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueConnectionArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "connection_type",
				Description: "The type of the connection. Currently, SFTP is not supported.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time that this connection definition was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_by",
				Description: "The user, group, or role that last updated this connection definition.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_updated_time",
				Description: "The last time that this connection definition was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "connection_properties",
				Description: "These key-value pairs define parameters for the connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "match_criteria",
				Description: "A list of criteria that can be used in selecting this connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "physical_connection_requirements",
				Description: "A map of physical connection requirements, such as virtual private cloud (VPC) and SecurityGroup, that are needed to make this connection successfully.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTagsForGlueConnection,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueConnectionArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueConnections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_connection.listGlueConnections", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	input := &glue.GetConnectionsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if d.EqualsQuals["connection_type"] != nil {
		connectionType := d.EqualsQuals["connection_type"].GetStringValue()
		if connectionType == "" || strings.EqualFold(connectionType, "SFTP") {
			return nil, nil
		}
		input.Filter = &types.GetConnectionsFilter{
			ConnectionType: types.ConnectionType(connectionType),
		}
	}

	// List call
	paginator := glue.NewGetConnectionsPaginator(svc, input, func(o *glue.GetConnectionsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glue_connection.listGlueConnections", "api_error", err)
			return nil, err
		}
		for _, connection := range output.ConnectionList {
			d.StreamListItem(ctx, connection)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueConnection(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_connection.getGlueConnection", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &glue.GetConnectionInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.GetConnection(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_connection.getGlueConnection", "api_error", err)
		return nil, err
	}
	return *data.Connection, nil
}

func getTagsForGlueConnection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn, _ := getGlueConnectionArn(ctx, d, h)
	return getTagsForGlueResource(ctx, d, arn.(string))
}

func getGlueConnectionArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	data := h.Item.(types.Connection)

	// Get common columns

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_connection.getGlueConnectionArn", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// arn format - https://docs.aws.amazon.com/glue/latest/dg/glue-specifying-resource-arns.html
	// arn:aws:glue:region:account-id:connection/connection-name
	arn := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":connection/" + *data.Name

	return arn, nil
}
