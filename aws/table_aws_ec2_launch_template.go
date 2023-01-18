package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2LaunchTemplate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_launch_template",
		Description: "AWS EC2 Launch Template",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("launch_template_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAction", "InvalidParameter", "InvalidParameterValue"}),
			},
			Hydrate: getEc2LaunchTemplate,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2LaunchTemplates,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "launch_template_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "launch_template_name",
				Description: "The name of the launch template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "launch_template_id",
				Description: "The ID of the launch template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time launch template was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "The principal that created the launch template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_version_number",
				Description: "The version number of the default version of the launch template.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "latest_version_number",
				Description: "The name of the Application-Layer Protocol Negotiation (ALPN) policy.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tags_src",
				Description: "The tags for the launch template.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(launchTemplateTurbotTags),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2LaunchTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_launch_template.listEc2LaunchTemplates", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(200)
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

	input := &ec2.DescribeLaunchTemplatesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	name := d.KeyColumnQualString("launch_template_name")
	if name != "" {
		input.LaunchTemplateNames = []string{name}
	}

	paginator := ec2.NewDescribeLaunchTemplatesPaginator(svc, input, func(o *ec2.DescribeLaunchTemplatesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_load_balancer_listener.listEc2LoadBalancers", "api_error", err)
			return nil, err
		}

		for _, items := range output.LaunchTemplates {
			d.StreamListItem(ctx, items)
		}

	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2LaunchTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of load balancer
	launchTemplateId := d.KeyColumnQuals["launch_template_id"].GetStringValue()

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getEc2Instance", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeLaunchTemplatesInput{
		LaunchTemplateIds: []string{launchTemplateId},
	}

	op, err := svc.DescribeLaunchTemplates(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_instance.getEc2Instance", "api_error", err)
		return nil, err
	}

	if op.LaunchTemplates != nil && len(op.LaunchTemplates) > 0 {

		return op.LaunchTemplates[0], nil
	}

	return nil, err
}

//// TRANSFORM FUNCTIONS

func launchTemplateTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(types.LaunchTemplate)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
