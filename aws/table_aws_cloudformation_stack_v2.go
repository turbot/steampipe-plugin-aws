package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsCloudFormation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudformation_stack",
		Description: "AWS CloudFormation Stack",
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.SingleColumn("name"),
		// 	IgnoreConfig: &plugin.IgnoreConfig{
		// 		ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationError", "ResourceNotFoundException"}),
		// 	},
		// 	Hydrate: getCloudFormationStack,
		// },
		List: &plugin.ListConfig{
			Hydrate: listCloudFormation,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Unique identifier of the stack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StackId"),
			},
			{
				Name:        "name",
				Description: "The name associated with the stack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StackName"),
			},
			{
				Name:        "status",
				Description: "Current status of the stack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StackStatus"),
			},
			{
				Name:        "creation_time",
				Description: "The time at which the stack was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "disable_rollback",
				Description: "Boolean to enable or disable rollback on stack creation failures.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "enable_termination_protection",
				Description: "Specifies whether termination protection is enabled for the stack.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_updated_time",
				Description: "The time the stack was last updated. This field will only be returned if the stack has been updated at least once.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "parent_id",
				Description: "ID of the direct parent of this stack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of an AWS Identity and Access Management (IAM) role that is associated with the stack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleARN"),
			},
			{
				Name:        "root_id",
				Description: "ID of the top-level stack to which the nested stack ultimately belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-defined description associated with the stack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "notification_arns",
				Description: "SNS topic ARNs to which stack related events are published.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NotificationARNs"),
			},
			{
				Name:        "outputs",
				Description: "A list of output structures.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "rollback_configuration",
				Description: "The rollback triggers for AWS CloudFormation to monitor during stack creation and updating operations, and for the specified monitoring period afterwards.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "capabilities",
				Description: "The capabilities allowed in the stack.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "stack_drift_status",
				Description: "Status of the stack's actual configuration compared to its expected template configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DriftInformation.StackDriftStatus"),
			},
			{
				Name:        "parameters",
				Description: "A list of Parameter structures.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "template_body",
				Description: "Structure containing the template body.",
				Type:        proto.ColumnType_STRING,
				// Hydrate:     getStackTemplate,
				Transform:   transform.FromField("TemplateBody"),
			},
			{
				Name:        "template_body_json",
				Description: "Structure containing the template body. Parsed into json object for better readability.",
				Type:        proto.ColumnType_JSON,
				// Hydrate:     getStackTemplate,
				Transform:   transform.FromField("TemplateBody").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "resources",
				Description: "A list of Stack resource structures.",
				Type:        proto.ColumnType_JSON,
				// Hydrate:     describeStackResources,
				Transform:   transform.FromField("StackResources"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with stack.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				// Transform:   transform.From(cfnStackTagsToTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StackName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("StackId").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudFormation(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudFormationService(ctx, d)
	if err != nil {
		return nil, err
	}

	// We can not pass the MaxResult value in param so we can't limit the result per page
	input := &cloudformation.DescribeStacksInput{}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.StackName = types.String(equalQuals["name"].GetStringValue())
	}
	cloudformation.NewDescribeStacksPaginator(svc, input ,func(o *cloudformation.DescribeStacksPaginatorOptions){
		o.StopOnDuplicateToken=true
	}){
	}

	}
  // 	err = svc.DescribeStacksPages(
	// 	input,
	// 	func(page *cloudformation.DescribeStacksOutput, lastPage bool) bool {
	// 		for _, stack := range page.Stacks {
	// 			d.StreamListItem(ctx, stack)

	// 			// Context can be cancelled due to manual cancellation or the limit has been hit
	// 			if d.QueryStatus.RowsRemaining(ctx) == 0 {
	// 				return false
	// 			}
	// 		}
	// 		return !lastPage
	// 	},
	// )
	// return nil, err
// }

//// HYDRATE FUNCTIONS

// func getCloudFormationStack(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
// 	// Create Session
// 	svc, err := CloudFormationService(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	name := d.KeyColumnQuals["name"].GetStringValue()
// 	params := &cloudformation.DescribeStacksInput{
// 		StackName: aws.String(name),
// 	}

// 	op, err := svc.DescribeStacks(params)
// 	if err != nil {
// 		plugin.Logger(ctx).Debug("getCloudFormationStack__", "ERROR", err)
// 		return nil, err
// 	}

// 	if len(op.Stacks) > 0 {
// 		return op.Stacks[0], nil
// 	}

// 	return nil, nil
// }

// func getStackTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("getStackTemplate")
// 	stack := h.Item.(*cloudformation.Stack)

// 	// Create Session
// 	svc, err := CloudFormationService(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// template_body is the template in its original string form
// 	params := &cloudformation.GetTemplateInput{
// 		StackName: stack.StackName,
// 	}
// 	stackTemplate, err := svc.GetTemplate(params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return stackTemplate, nil
// }

// func describeStackResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("getStackTemplate")
// 	stack := h.Item.(*cloudformation.Stack)

// 	// Create Session
// 	svc, err := CloudFormationService(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	params := &cloudformation.DescribeStackResourcesInput{
// 		StackName: stack.StackName,
// 	}

// 	stackResources, err := svc.DescribeStackResources(params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return stackResources, nil
// }

// //// TRANSFORM FUNCTIONS

// func cfnStackTagsToTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	stack := d.HydrateItem.(*cloudformation.Stack)
// 	var turbotTagsMap map[string]string

// 	if stack.Tags != nil {
// 		turbotTagsMap = map[string]string{}
// 		for _, i := range stack.Tags {
// 			turbotTagsMap[*i.Key] = *i.Value
// 		}
// 	}
// 	return turbotTagsMap, nil
// }
