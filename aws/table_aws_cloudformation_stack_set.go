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

func tableAwsCloudFormationStackSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudformation_stack_set",
		Description: "AWS CloudFormation Stack Set",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("stack_set_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"StackSetNotFoundException"}),
			},
			Hydrate: getCloudFormationStackSet,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudFormationStackSets,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "status",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudformationv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "stack_set_id",
				Description: "The ID of the stack set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stack_set_name",
				Description: "The name of the stack set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the stack set.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFormationStackSet,
				Transform:   transform.FromField("StackSetARN"),
			},
			{
				Name:        "status",
				Description: "The status of the stack set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the stack set that you specify when the stack set is created or updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "drift_status",
				Description: "Status of the stack set's actual configuration compared to its expected template and parameter configuration. A stack set is considered to have drifted if one or more of its stack instances have drifted from their expected template and parameter configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_drift_check_timestamp",
				Description: "Most recent time when CloudFormation performed a drift detection operation on the stack set.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "permission_model",
				Description: "Describes how the IAM roles required for stack set operations are created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administration_role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role used to create or update the stack set.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFormationStackSet,
				Transform:   transform.FromField("AdministrationRoleARN"),
			},
			{
				Name:        "execution_role_name",
				Description: "The name of the IAM execution role used to create or update the stack set.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFormationStackSet,
			},
			{
				Name:        "template_body",
				Description: "The structure that contains the body of the template that was used to create or update the stack set.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudFormationStackSet,
			},
			{
				Name:        "auto_deployment",
				Description: "Describes whether StackSets automatically deploys to Organizations accounts that are added to a target organizational unit (OU).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "capabilities",
				Description: "The capabilities that are allowed in the stack set.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFormationStackSet,
			},
			{
				Name:        "organizational_unit_ids",
				Description: "The organization root ID or organizational unit (OU) IDs that you specified for DeploymentTargets.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFormationStackSet,
			},
			{
				Name:        "parameters",
				Description: "A list of input parameters for a stack set.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFormationStackSet,
			},
			{
				Name:        "stack_set_drift_detection_details",
				Description: "Detailed information about the drift status of the stack set.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFormationStackSet,
			},
			{
				Name:        "managed_execution",
				Description: "Describes whether StackSets performs non-conflicting operations concurrently and queues conflicting operations.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with stack.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFormationStackSet,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFormationStackSet,
				Transform:   transform.From(cfnStackSetTagsToTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StackSetName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudFormationStackSet,
				Transform:   transform.FromField("StackSetARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudFormationStackSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack_set.listCloudFormationStackSets", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
				maxLimit = limit
		}
	}
	status := d.EqualsQualString("status")

	input := &cloudformation.ListStackSetsInput{
		MaxResults: &maxLimit,
	}

	if status != "" {
		input.Status = types.StackSetStatus(status)
	}

	paginator := cloudformation.NewListStackSetsPaginator(svc, input, func(o *cloudformation.ListStackSetsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudformation_stack_set.listCloudFormationStackSets", "api_error", err)
			return nil, err
		}

		for _, stackSet := range output.Summaries {
			d.StreamListItem(ctx, stackSet)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudFormationStackSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var name string

	if h.Item != nil {
		data := h.Item.(types.StackSetSummary)
		name = *data.StackSetName
	} else {
		name = d.EqualsQualString("stack_set_name")
	}

	// Empty check
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := CloudFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack_set.getCloudFormationStackSet", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}


	params := &cloudformation.DescribeStackSetInput{
		StackSetName: aws.String(name),
	}

	op, err := svc.DescribeStackSet(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack_set.getCloudFormationStackSet", "api_error", err)
		return nil, err
	}

	return op.StackSet, nil
}

//// TRANSFORM FUNCTIONS

func cfnStackSetTagsToTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	stack := d.HydrateItem.(*types.StackSet)
	var turbotTagsMap map[string]string
	if len(stack.Tags) > 0 {
		if stack.Tags != nil {
			turbotTagsMap = map[string]string{}
			for _, i := range stack.Tags {
				turbotTagsMap[*i.Key] = *i.Value
			}
		}
	}
	return turbotTagsMap, nil
}
