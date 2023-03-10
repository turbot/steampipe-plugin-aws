package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	cloudformationv1 "github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFormationStackResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudformation_stack_resource",
		Description: "AWS CloudFormation Stack Resource",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"stack_name", "logical_resource_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationError", "ResourceNotFoundException"}),
			},
			Hydrate: getCloudFormationStackResource,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCloudFormationStacks,
			Hydrate:       listCloudFormationStackResources,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "stack_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudformationv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "logical_resource_id",
				Description: "The logical name of the resource specified in the template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stack_name",
				Description: "The name associated with the stack.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFormationStackResource,
				Transform:   transform.FromField("StackName"),
			},
			{
				Name:        "stack_id",
				Description: "Unique identifier of the stack.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFormationStackResource,
				Transform:   transform.FromField("StackId"),
			},
			{
				Name:        "last_updated_timestamp",
				Description: "Time the status was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "resource_status",
				Description: "Current status of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "Type of resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "User defined description associated with the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFormationStackResource,
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "physical_resource_id",
				Description: "The name or unique identifier that corresponds to a physical instance ID of a resource supported by CloudFormation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_status_reason",
				Description: "Success/failure message associated with the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "drift_information",
				Description: "Information about whether the resource's actual configuration differs, or has drifted, from its expected configuration, as defined in the stack template and any values specified as template parameters. For more information, see Detecting Unregulated Configuration Changes to Stacks and Resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "module_info",
				Description: "Contains information about the module from which the resource was created, if the resource was created from a module included in the stack template.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogicalResourceId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudFormationStackResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	stack := h.Item.(types.Stack)
	stackName := d.EqualsQualString("stack_name")

	// If a stack name is specified in optional quals, the user is not allowed to perform API calls for other stacks.
	if stackName == "" || stackName != *stack.StackName {
		return nil, nil
	}

	// Create session
	svc, err := CloudFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack_resource.listCloudFormationStackResources", "connection_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	// We can not pass the MaxResult value in param so we can't limit the result per page
	input := &cloudformation.ListStackResourcesInput{
		StackName: stack.StackName,
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["stack_name"] != nil {
		input.StackName = aws.String(equalQuals["stack_name"].GetStringValue())
	}
	paginator := cloudformation.NewListStackResourcesPaginator(svc, input, func(o *cloudformation.ListStackResourcesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudformation_stack_resource.listCloudFormationStackResources", "api_error", err)
			return nil, err
		}

		for _, resource := range output.StackResourceSummaries {
			d.StreamListItem(ctx, types.StackResourceDetail{
				LastUpdatedTimestamp: resource.LastUpdatedTimestamp,
				LogicalResourceId:    resource.LogicalResourceId,
				ResourceStatus:       resource.ResourceStatus,
				ResourceType:         resource.ResourceType,
				DriftInformation:     (*types.StackResourceDriftInformation)(resource.DriftInformation),
				ModuleInfo:           resource.ModuleInfo,
				PhysicalResourceId:   resource.PhysicalResourceId,
				ResourceStatusReason: resource.ResourceStatusReason,
				StackName:            stack.StackName,
				StackId:              stack.StackId,
			})
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudFormationStackResource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var stackName, logicalId string
	if h.Item != nil {
		data := h.Item.(types.StackResourceDetail)
		stackName = *data.StackName
		logicalId = *data.LogicalResourceId
	} else {
		stackName = d.EqualsQualString("stack_name")
		logicalId = d.EqualsQualString("logical_resource_id")
	}

	// Empty param check
	if stackName == "" || logicalId == "" {
		return nil, nil
	}

	// Create Session
	svc, err := CloudFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack_resource.getCloudFormationStackResource", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &cloudformation.DescribeStackResourceInput{
		StackName:         aws.String(stackName),
		LogicalResourceId: aws.String(logicalId),
	}

	op, err := svc.DescribeStackResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack_resource.getCloudFormationStackResource", "api_error", err)
		return nil, err
	}

	return op.StackResourceDetail, nil
}
