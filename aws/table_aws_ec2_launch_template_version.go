package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2Endpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2LaunchTemplateVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_launch_template_version",
		Description: "AWS EC2 Launch Template Version",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"launch_template_id", "version_number"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidLaunchTemplateId.NotFound", "InvalidLaunchTemplateId.VersionNotFound"}),
			},
			Hydrate: getEc2LaunchTemplateVersion,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeLaunchTemplateVersions"},
		},
		List: &plugin.ListConfig{
			// IgnoreConfig: &plugin.IgnoreConfig{
			// 	ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidLaunchTemplateName.NotFoundException", "InvalidLaunchTemplateId.NotFound", "InvalidLaunchTemplateId.Malformed"}),
			// },
			ParentHydrate: listEc2LaunchTemplates,
			Hydrate:       listEc2LaunchTemplateVersions,
			Tags:          map[string]string{"service": "ec2", "action": "DescribeLaunchTemplateVersions"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "launch_template_id",
					Require: plugin.Optional,
				},
				{
					Name:    "launch_template_name",
					Require: plugin.Optional,
				},
				{
					Name:    "default_version",
					Require: plugin.Optional,
				},
				{
					Name:    "ebs_optimized",
					Require: plugin.Optional,
				},
				{
					Name:    "image_id",
					Require: plugin.Optional,
				},
				{
					Name:    "instance_type",
					Require: plugin.Optional,
				},
				{
					Name:    "kernel_id",
					Require: plugin.Optional,
				},
				{
					Name:    "ram_disk_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2Endpoint.EC2ServiceID),
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
				Description: "The time the version was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "The principal that created the version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_version",
				Description: "Indicates whether the version is the default version.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "disable_api_stop",
				Description: "Indicates whether the instance is enabled for stop protection.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LaunchTemplateData.DisableApiStop"),
			},
			{
				Name:        "disable_api_termination",
				Description: "If set to true, indicates that the instance cannot be terminated using the Amazon EC2 console, command line tool, or API.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LaunchTemplateData.DisableApiTermination"),
			},
			{
				Name:        "ebs_optimized",
				Description: "Indicates whether the instance is optimized for Amazon EBS I/O.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LaunchTemplateData.EbsOptimized"),
			},
			{
				Name:        "image_id",
				Description: "The ID of the AMI or a Systems Manager parameter.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateData.ImageId"),
			},
			{
				Name:        "instance_type",
				Description: "The instance type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateData.InstanceType"),
			},
			{
				Name:        "kernel_id",
				Description: "The ID of the kernel, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateData.KernelId"),
			},
			{
				Name:        "key_name",
				Description: "The name of the key pair.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateData.KeyName"),
			},
			{
				Name:        "ram_disk_id",
				Description: "The ID of the RAM disk, if applicable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateData.RamDiskId"),
			},
			{
				Name:        "security_groups",
				Description: "The security group names.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateData.SecurityGroups"),
			},
			{
				Name:        "security_group_ids",
				Description: "The security group IDs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateData.SecurityGroupIds"),
			},
			{
				Name:        "version_description",
				Description: "The description for the version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_number",
				Description: "The version number.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "user_data",
				Description: "The user data of the launch template.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateData.UserData").Transform(base64DecodedData),
			},
			{
				Name:        "launch_template_data",
				Description: "Information about the launch template.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LaunchTemplateName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2LaunchTemplateVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	launchTemplate := h.Item.(types.LaunchTemplate)
	launchTemplateName := d.EqualsQualString("launch_template_name")
	launchTemplateId := d.EqualsQualString("launch_template_id")

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_launch_template_version.listEc2LaunchTemplateVersions", "connection_error", err)
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

	if launchTemplateId != "" {
		if launchTemplateId != *launchTemplate.LaunchTemplateId {
			return nil, nil
		}
	}
	if launchTemplateName != "" {
		if launchTemplateName != *launchTemplate.LaunchTemplateName {
			return nil, nil
		}
	}

	if launchTemplateName != "" && launchTemplateId != "" {
		return nil, fmt.Errorf("Both LaunchtemplateName and LaunchTemplateId cannot be passed in the where clause")
	}

	// The aws_ec2_launch_template table is used as the parent hydrate because the LaunchTemplateId is not specified in the input parameter, and it will return only the latest and default version launch templates.
	input := &ec2.DescribeLaunchTemplateVersionsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if launchTemplateId != "" {
		input.LaunchTemplateId = &launchTemplateId
	}
	if launchTemplateName != "" {
		input.LaunchTemplateName = &launchTemplateName
	}
	if launchTemplateId == "" && launchTemplateName == "" {
		input.LaunchTemplateId = launchTemplate.LaunchTemplateId
	}

	filters := buildEc2LaunchTemplateVersionFilter(d.Quals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeLaunchTemplateVersionsPaginator(svc, input, func(o *ec2.DescribeLaunchTemplateVersionsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_launch_template_version.listEc2LaunchTemplateVersions", "api_error", err)
			return nil, err
		}

		for _, items := range output.LaunchTemplateVersions {
			d.StreamListItem(ctx, items)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2LaunchTemplateVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_launch_template_version.getEc2LaunchTemplateVersion", "connection_error", err)
		return nil, err
	}

	var templateId string
	var version int64
	qual := d.Quals
	if qual["launch_template_id"] != nil {
		for _, q := range qual["launch_template_id"].Quals {
			templateId = q.Value.GetStringValue()
		}
	}
	if qual["version_number"] != nil {
		for _, q := range qual["version_number"].Quals {
			version = q.Value.GetInt64Value()
		}
	}

	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(templateId),
		Versions:         []string{fmt.Sprint(version)},
	}

	op, err := svc.DescribeLaunchTemplateVersions(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_launch_template_version.getEc2LaunchTemplateVersion", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.LaunchTemplateVersions) > 0 {
		return op.LaunchTemplateVersions[0], nil
	}

	return nil, err
}

//// UTILITY FUNCTIONS

// Build ec2 launch template version list call input filter
func buildEc2LaunchTemplateVersionFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"default_version": "is-default-version",
		"ebs_optimized":   "ebs-optimized",
		"image_id":        "image-id",
		"instance_type":   "instance-type",
		"kernel_id":       "kernel-id",
		"ram_disk_id":     "ram-disk-id",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			if columnName == "default_version" || columnName == "ebs_optimized" {
				value := getQualsValueByColumn(quals, columnName, "boolean")
				filter.Values = []string{fmt.Sprint(value)}
			} else {
				value := getQualsValueByColumn(quals, columnName, "string")
				val, ok := value.(string)
				if ok {
					filter.Values = []string{val}
				}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
