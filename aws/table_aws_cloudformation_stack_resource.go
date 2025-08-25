package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go"

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
			Tags:    map[string]string{"service": "cloudformation", "action": "DescribeStackResource"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCloudFormationStackResourcesParent,
			Hydrate:       listCloudFormationStackResources,
			Tags:          map[string]string{"service": "cloudformation", "action": "ListStackResources"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "stack_name",
					Require: plugin.Optional,
				},
				{
					Name:    "physical_resource_id",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCloudFormationStackResource,
				Tags: map[string]string{"service": "cloudformation", "action": "DescribeStackResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CLOUDFORMATION_SERVICE_ID),
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

//// PARENT HYDRATE FUNCTION

// Parent hydrate function that optimizes the query flow based on the provided qualifiers
// This function implements the performance optimization strategy described in the design:
// 1. Fast access by PhysicalResourceId: Uses DescribeStackResources with PhysicalResourceId parameter
// 2. Fast access by StackName: Falls back to listing all stacks for stack-based queries
//
// The function handles the physical_resource_id qualifier efficiently by:
// - Making a single DescribeStackResources API call with PhysicalResourceId parameter
// - Streaming results directly if < 100 resources (fast path)
// - Falling back to stack iteration if ≥ 100 resources (rare case)
// - Handling ValidationError gracefully when resource is not found
//
// This approach addresses the performance issue described in GitHub issue #2627 where
// the previous implementation was slow because it iterated through all stacks.
// The new implementation matches the CLI behavior: aws cloudformation describe-stack-resources --physical-resource-id
func listCloudFormationStackResourcesParent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	physicalResourceId := d.EqualsQualString("physical_resource_id")

	// Case 3: Fast access by PhysicalResourceId using DescribeStackResources
	// This is the optimized path that addresses the performance issue in GitHub issue #2627
	// Instead of iterating through all stacks, we make a single API call with PhysicalResourceId
	if physicalResourceId != "" {
		// Create session
		svc, err := CloudFormationClient(ctx, d)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudformation_stack_resource.listCloudFormationStackResourcesParent", "connection_error", err)
			return nil, err
		}

		// Unsupported region check
		if svc == nil {
			return nil, nil
		}

		// Apply rate limiting
		d.WaitForListRateLimit(ctx)

		// Use DescribeStackResources directly with PhysicalResourceId parameter
		// This is the key optimization: single API call instead of iterating all stacks
		// This matches the CLI behavior: aws cloudformation describe-stack-resources --physical-resource-id
		input := &cloudformation.DescribeStackResourcesInput{
			PhysicalResourceId: aws.String(physicalResourceId),
		}

		output, err := svc.DescribeStackResources(ctx, input)
		if err != nil {
			// Handle ValidationError gracefully - this occurs when the physical resource ID is not found
			// This matches the CLI behavior where it returns "Stack for my-arn-1 does not exist"
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "ValidationError" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_cloudformation_stack_resource.listCloudFormationStackResourcesParent", "api_error", err)
			return nil, err
		}

		// Fast path: If response has less than 100 resources, stream them directly
		// This covers the vast majority of cases and provides optimal performance
		if len(output.StackResources) < 100 {
			for _, resource := range output.StackResources {
				d.StreamListItem(ctx, resource)
				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
			return nil, nil
		} else {
			// Fallback path: If response has 100 or more resources (rare case),
			// return the stack list for the list function to handle with pagination
			// This ensures we don't lose the ability to handle large stacks
			return listCloudFormationStacks(ctx, d, h)
		}
	}

	// Case 1: Fast access by StackName or general stack listing
	// For queries without physical_resource_id, we need to list all stacks
	// This is the traditional approach for stack-based queries
	return listCloudFormationStacks(ctx, d, h)
}

//// LIST FUNCTION

func listCloudFormationStackResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	stackName := d.EqualsQualString("stack_name")

	// Check the type of h.Item and handle accordingly
	// This function handles different types of items returned from the parent hydrate:
	// - types.StackResource: Direct resources from physical_resource_id queries (fast path)
	// - types.StackSummary: Stack information for stack-based queries
	if h.Item != nil {
		// Case 3: If h.Item is types.StackResource, directly stream it
		if resource, ok := h.Item.(types.StackResource); ok {
			d.StreamListItem(ctx, resource)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
			return nil, nil
		}

		// Case 1: If h.Item is types.StackSummary, get stack name and make API call
		if stack, ok := h.Item.(types.StackSummary); ok {

			// If a stack name is specified in optional quals, the user is not allowed to perform API calls for other stacks.
			if stackName != "" && stackName != *stack.StackName {
				return nil, nil
			}

			// For deleted stacks, skip processing
			if stack.StackStatus == types.StackStatusDeleteComplete {
				return nil, nil
			}

			return listCloudFormationStackResourcesByStackName(ctx, d, svc, *stack.StackName)
		}
	}

	return nil, nil
}

//// HELPER FUNCTIONS

// Case 1: Fast access by StackName using ListStackResources
// This function handles stack-based queries by using ListStackResources with pagination
// It's used when the parent hydrate returns StackSummary items for stack iteration
// This approach is used for:
// - Queries with stack_name qualifier
// - General stack resource listing (when no physical_resource_id is specified)
func listCloudFormationStackResourcesByStackName(ctx context.Context, d *plugin.QueryData, svc *cloudformation.Client, stackName string) (interface{}, error) {

	// Use ListStackResources to get resources for the stack
	listInput := &cloudformation.ListStackResourcesInput{
		StackName: aws.String(stackName),
	}

	paginator := cloudformation.NewListStackResourcesPaginator(svc, listInput, func(o *cloudformation.ListStackResourcesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// Apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudformation_stack_resource.listCloudFormationStackResourcesByStackName", "api_error", err)
			return nil, err
		}

		for _, resourceSummary := range output.StackResourceSummaries {
			d.StreamListItem(ctx, types.StackResourceDetail{
				LastUpdatedTimestamp: resourceSummary.LastUpdatedTimestamp,
				LogicalResourceId:    resourceSummary.LogicalResourceId,
				ResourceStatus:       resourceSummary.ResourceStatus,
				ResourceType:         resourceSummary.ResourceType,
				DriftInformation:     (*types.StackResourceDriftInformation)(resourceSummary.DriftInformation),
				ModuleInfo:           resourceSummary.ModuleInfo,
				PhysicalResourceId:   resourceSummary.PhysicalResourceId,
				ResourceStatusReason: resourceSummary.ResourceStatusReason,
				StackName:            aws.String(stackName),
			})
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudFormationStackResource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var stackName, logicalId string
	if h.Item != nil {
		switch item := h.Item.(type) {
		case types.StackResourceDetail:
			stackName = *item.StackName
			logicalId = *item.LogicalResourceId
		case types.StackResource:
			return item, nil
		}
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
