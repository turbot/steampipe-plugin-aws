package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueConnection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_connection",
		Description: "AWS Glue Connection",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlueConnection,
		},
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"connection_type"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInputException"}),
			},
			Hydrate: listGlueConnections,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := GlueService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_connection.listGlueConnections", "service_creation_error", err)
		return nil, err
	}

	input := &glue.GetConnectionsInput{
		MaxResults: aws.Int64(100),
	}

	if d.KeyColumnQuals["connection_type"] != nil {
		connectionType := d.KeyColumnQuals["connection_type"].GetStringValue()
		if connectionType == "" || strings.EqualFold(connectionType, "SFTP") {
			return nil, nil
		}
		input.SetFilter(&glue.GetConnectionsFilter{
			ConnectionType: aws.String(connectionType),
		})
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
	err = svc.GetConnectionsPages(
		input,
		func(page *glue.GetConnectionsOutput, isLast bool) bool {
			for _, connection := range page.ConnectionList {
				d.StreamListItem(ctx, connection)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_connection.listGlueConnections", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueConnection(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_connection.getGlueConnection", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &glue.GetConnectionInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.GetConnection(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_connection.getGlueConnection", "api_error", err)
		return nil, err
	}
	return data.Connection, nil
}

func getGlueConnectionArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*glue.Connection)

	// Get common columns
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// arn format - https://docs.aws.amazon.com/glue/latest/dg/glue-specifying-resource-arns.html
	// arn:aws:glue:region:account-id:connection/connection-name
	arn := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":connection/" + *data.Name

	return arn, nil
}
