package aws

import (
	"context"
	"encoding/json"
	"net/url"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	go_kit "github.com/turbot/go-kit/types"

	cloudformationv1 "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/goccy/go-yaml"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudFormationStack(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudformation_stack",
		Description: "AWS CloudFormation Stack",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationError", "ResourceNotFoundException"}),
			},
			Hydrate: getCloudFormationStack,
			Tags:    map[string]string{"service": "cloudformation", "action": "DescribeStacks"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudFormationStacks,
			Tags:    map[string]string{"service": "cloudformation", "action": "DescribeStacks"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getStackTemplate,
				Tags: map[string]string{"service": "cloudformation", "action": "GetTemplate"},
			},
			{
				Func: describeStackResources,
				Tags: map[string]string{"service": "cloudformation", "action": "DescribeStackResources"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudformationv1.EndpointsID),
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
				Hydrate:     getStackTemplate,
				Transform:   transform.FromField("TemplateBody").Transform(transform.ToString),
			},
			{
				Name:        "template_body_json",
				Description: "Structure containing the template body. Parsed into json object for better readability.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getStackTemplate,
				Transform:   transform.FromField("TemplateBody").Transform(formatJsonBody),
			},
			{
				Name:        "resources",
				Description: "A list of Stack resource structures.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeStackResources,
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
				Transform:   transform.From(cfnStackTagsToTurbotTags),
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

func listCloudFormationStacks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack.listCloudFormationStacks", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// We can not pass the MaxResult value in param so we can't limit the result per page
	input := &cloudformation.DescribeStacksInput{}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		input.StackName = aws.String(equalQuals["name"].GetStringValue())
	}
	paginator := cloudformation.NewDescribeStacksPaginator(svc, input, func(o *cloudformation.DescribeStacksPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudformation_stack.listCloudFormationStacks", "api_error", err)
			return nil, err
		}

		for _, stack := range output.Stacks {
			d.StreamListItem(ctx, stack)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudFormationStack(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := CloudFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack.getCloudFormationStack", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	name := d.EqualsQuals["name"].GetStringValue()
	params := &cloudformation.DescribeStacksInput{
		StackName: aws.String(name),
	}

	op, err := svc.DescribeStacks(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack.getCloudFormationStack", err)
		return nil, err
	}

	if len(op.Stacks) > 0 {
		return op.Stacks[0], nil
	}

	return nil, nil
}

func getStackTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	stack := h.Item.(types.Stack)

	// Create Session
	svc, err := CloudFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack.getStackTemplate", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// template_body is the template in its original string form
	params := &cloudformation.GetTemplateInput{
		StackName: stack.StackName,
	}
	stackTemplate, err := svc.GetTemplate(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack.getStackTemplate", "api_error", err)
		return nil, err
	}

	return stackTemplate, nil
}

func describeStackResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	stack := h.Item.(types.Stack)

	// Create Session
	svc, err := CloudFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack.describeStackResources", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &cloudformation.DescribeStackResourcesInput{
		StackName: stack.StackName,
	}

	stackResources, err := svc.DescribeStackResources(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudformation_stack.describeStackResources", "api_error", err)
		return nil, err
	}

	return stackResources, nil
}

// // TRANSFORM FUNCTIONS
func cfnStackTagsToTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	stack := d.HydrateItem.(types.Stack)
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

/*
* Encountering errors when parsing YAML to JSON with `transform.UnmarshalYAML` in the Steampipe SDK, primarily due to `url.QueryUnescape()` usage.
* The API may return the template body in various formats like JSON or YAML, based on configuration made by user. Consequently, converting YAML to JSON isn't always required.
* To address this, a dedicated function has been implemented. It ensures the template body is correctly formatted to JSON as per the API's output, adapting the transformation process accordingly.
* This approach is specifically designed for this table, providing a tailored and accurate response handling as per the API specifications.
*/

// Functionality Overview
// Identifies and decodes URLs within the template body. (URLs are decoded selectively to avoid issues with '%' characters not part of a valid escaped sequence.)
// Integrates the decoded URLs back into the original template body.
// Determines if the template body is in JSON or YAML format.
// Directly unmarshal the content if it's a JSON string.
// Converts YAML format to JSON if the template body is in YAML.
func formatJsonBody(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	inputStr := go_kit.SafeString(d.Value)
	var result interface{}
	if inputStr != "" {

		// Unescape only URLs instead of checking if any '%' character is not followed by two hexadecimal digits in the template body.
		// QueryUnescape does the inverse transformation of QueryEscape, converting each 3-byte encoded substring of the form "%AB" into the hex-decoded byte 0xAB. It returns an error if any % is not followed by two hexadecimal digits.
		// Regex to match URLs
		regex := regexp.MustCompile(`(https?://[^\s]+)`)

		// Find all URLs in the input string
		matches := regex.FindAllString(inputStr, -1)

		// Iterate the URLs present in the template body.
		for _, match := range matches {
			// The `QueryUnescape()` function returns an error if any '%' character is not followed by two hexadecimal digits while unescaping the URL.
			// The template body may contain instances where '%' is not followed by two hexadecimal digits.
			// We should Unescape only the URLs not all the template body.
			decoded, err := url.QueryUnescape(match)
			if err != nil {
				return nil, err
			}
			inputStr = strings.ReplaceAll(inputStr, match, decoded)
		}

		// Check if the template body return by API is in JSON/YAML format
		if isJSON(inputStr) {
			err := json.Unmarshal([]byte(inputStr), &result)
			if err != nil {
				return nil, err
			}
		} else {
			// Occasionally, the API response includes the template body with carriage return characters (`\r`).
			// This leads to an "Unmarshal error [2:4] unexpected key name" when attempting to parse the YAML string to JSON.
			// Therefore, it's necessary to substitute the carriage return character `\r` with `\n` to resolve this issue.
			inputStr = strings.ReplaceAll(inputStr, "\r", "\n")
			err := yaml.Unmarshal([]byte(inputStr), &result)
			if err != nil {
				return nil, err
			}
		}

	}
	return result, nil
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
