package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2LaunchTemplate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_launch_template",
		Description: "AWS EC2 Launch Template",
		List: &plugin.ListConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidLaunchTemplateName.NotFoundException", "InvalidLaunchTemplateId.NotFound", "InvalidLaunchTemplateId.Malformed"}),
			},
			Hydrate: listEc2LaunchTemplates,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeLaunchTemplates"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "launch_template_name",
					Require: plugin.Optional,
				},
				{
					Name:    "launch_template_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_EC2_SERVICE_ID),
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
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     launchTemplateAkas,
				Transform:   transform.FromValue(),
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
			maxLimit = limit
		}
	}

	input := &ec2.DescribeLaunchTemplatesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	name := d.EqualsQualString("launch_template_name")
	if name != "" {
		input.LaunchTemplateNames = []string{name}
	}
	if d.EqualsQualString("launch_template_id") != "" {
		input.LaunchTemplateIds = []string{d.EqualsQualString("launch_template_id")}
	}

	paginator := ec2.NewDescribeLaunchTemplatesPaginator(svc, input, func(o *ec2.DescribeLaunchTemplatesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// When a table is used as the parent hydrate in another table, the ignore config does not function as expected. Consequently, it becomes necessary to handle this situation manually.
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "InvalidLaunchTemplateId.NotFound" || ae.ErrorCode() == "InvalidLaunchTemplateName.NotFoundException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_ec2_launch_template.listEc2LaunchTemplates", "api_error", err)
			return nil, err
		}

		for _, items := range output.LaunchTemplates {
			d.StreamListItem(ctx, items)
		}

	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func launchTemplateAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	launchTemplate := h.Item.(types.LaunchTemplate)
	region := d.EqualsQualString(matrixKeyRegion)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_launch_template.launchTemplateAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for Turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":launch-template/" + *launchTemplate.LaunchTemplateId}

	return akas, nil
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
