package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamInstanceProfile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_instance_profile",
		Description: "AWS IAM Instance Profile",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"name", "arn"}),
			Hydrate:    getIamInstanceProfile,
			Tags:       map[string]string{"service": "iam", "action": "GetInstanceProfile"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationError", "NoSuchEntity", "InvalidParameter"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIamInstanceProfiles,
			Tags:    map[string]string{"service": "iam", "action": "ListInstanceProfiles"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "path", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getIamInstanceProfileTags,
				Tags: map[string]string{"service": "iam", "action": "ListInstanceProfileTags"},
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			// "Key" Columns
			{
				Name:        "name",
				Description: "The friendly name that identifies the instance profile.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceProfileName"),
			},
			{
				Name:        "arn",
				Type:        proto.ColumnType_STRING,
				Description: "The Amazon Resource Name (ARN) specifying the instance profile.",
			},
			{
				Name:        "instance_profile_id",
				Type:        proto.ColumnType_STRING,
				Description: "The stable and unique string identifying the instance profile.",
			},

			// Other Columns
			{
				Name:        "create_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time when the instance profile was created.",
			},
			{
				Name:        "path",
				Description: "The path to the instance profile. For more information about paths, see IAM identifiers in the IAM User Guide.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "roles",
				Description: "The role associated with the instance profile.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Roles").Transform(getInstanceProfileRoles),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the instance profile.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamInstanceProfileTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceProfileName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamInstanceProfileTags,
				Transform:   transform.From(getIamInstanceProfileTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listIamInstanceProfiles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_instance_profile.listIamInstanceProfiles", "client_error", err)
		return nil, err
	}

	maxItems := int32(1000)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	input := iam.ListInstanceProfilesInput{
		MaxItems: aws.Int32(maxItems),
	}

	// Additional filters
	if d.EqualsQualString("path") != "" {
		input.PathPrefix = aws.String(d.EqualsQualString("path"))
	}

	paginator := iam.NewListInstanceProfilesPaginator(svc, &input, func(o *iam.ListInstanceProfilesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_instance_profile.listIamInstanceProfiles", "api_error", err)
			return nil, err
		}

		for _, instanceProfile := range output.InstanceProfiles {
			d.StreamListItem(ctx, instanceProfile)

			// Context cancelled
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamInstanceProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_instance_profile.getIamInstanceProfile", "client_error", err)
		return nil, err
	}

	var name string
	if h.Item != nil {
		instanceProfile := h.Item.(types.InstanceProfile)
		name = *instanceProfile.InstanceProfileName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
		if name == "" {
			arn := d.EqualsQuals["arn"].GetStringValue()
			name = iamInstanceProfileNameFromArn(arn)
		}
	}

	params := &iam.GetInstanceProfileInput{InstanceProfileName: aws.String(name)}

	op, err := svc.GetInstanceProfile(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_instance_profile.getIamInstanceProfile", "api_error", err)
		return nil, err
	}

	return *op.InstanceProfile, nil
}

func getIamInstanceProfileTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceProfile := h.Item.(types.InstanceProfile)

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_instance_profile.getIamInstanceProfileTags", "client_error", err)
		return nil, err
	}

	params := &iam.ListInstanceProfileTagsInput{
		InstanceProfileName: instanceProfile.InstanceProfileName,
	}

	op, err := svc.ListInstanceProfileTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_instance_profile.getIamInstanceProfileTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// UTILITY FUNCTIONS

func getInstanceProfileRoles(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	roles := d.Value.([]types.Role)
	
	var rolesData []map[string]interface{}
	for _, role := range roles {
		roleData := map[string]interface{}{
			"arn":                   role.Arn,
			"name":                  role.RoleName,
			"role_id":               role.RoleId,
			"path":                  role.Path,
			"create_date":           role.CreateDate,
			"max_session_duration":  role.MaxSessionDuration,
			"description":           role.Description,
		}
		rolesData = append(rolesData, roleData)
	}
	
	return rolesData, nil
}

// Extract instance profile name from ARN
func iamInstanceProfileNameFromArn(arn string) string {
	// ARN format: arn:aws:iam::account-id:instance-profile/instance-profile-name
	parts := strings.Split(arn, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

//// TRANSFORM FUNCTIONS

func getIamInstanceProfileTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*iam.ListInstanceProfileTagsOutput)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
} 