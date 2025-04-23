package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rolesanywhere"
	"github.com/aws/aws-sdk-go-v2/service/rolesanywhere/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRolesAnywhereProfile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rolesanywhere_profile",
		Description: "AWS Roles Anywhere Profile",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"profile_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getProfile,
			Tags:    map[string]string{"service": "rolesanywhere", "action": "GetProfile"},
		},
		List: &plugin.ListConfig{
			Hydrate: listProfiles,
			Tags:    map[string]string{"service": "rolesanywhere", "action": "ListProfiles"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listTagsForProfile,
				Tags: map[string]string{"service": "rolesanywhere", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ROLESANYWHERE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "profile_id",
				Description: "The unique identifier of the profile.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the profile.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProfileArn"),
			},
			{
				Name:        "name",
				Description: "The name of the profile.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "accept_role_session_name",
				Description: "Used to determine if a custom role session name will be accepted in a temporary credential request.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "updated_at",
				Description: "The date and time when the profile was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_at",
				Description: "The date and time when the profile was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_by",
				Description: "The account that created the profile.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "duration_seconds",
				Description: "Used to determine how long sessions vended using this profile are valid for.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "enabled",
				Description: "If the profile is enabled or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "attribute_mappings",
				Description: "A mapping applied to the authenticating end-entity certificate.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "managed_policy_arns",
				Description: "A list of managed policy ARNs that apply to the vended session credentials.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "role_arns",
				Description: "A list of IAM roles that this profile can assume in a temporary credential request.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "require_instance_properties",
				Description: "Specifies whether instance properties are required in temporary credential requests with this profile.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "session_policy",
				Description: "A session policy that applies to the trust boundary of the vended session credentials.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SessionPolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "session_policy_std",
				Description: "Contains the session policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SessionPolicy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the profile",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForProfile,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProfileId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForProfile,
				Transform:   transform.From(profileTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ProfileArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listProfiles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := RolesAnywhereClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_profile.listProfiles", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	input := rolesanywhere.ListProfilesInput{}
	paginator := rolesanywhere.NewListProfilesPaginator(svc, &input, func(o *rolesanywhere.ListProfilesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rolesanywhere_profile.listProfiles", "api_error", err)
			return nil, err
		}

		for _, profile := range output.Profiles {
			d.StreamListItem(ctx, profile)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil

}

//// HYDRATE FUNCTIONS

func getProfile(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	profile_id := d.EqualsQuals["id"].GetStringValue()

	svc, err := RolesAnywhereClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_profile.getProfile", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	params := &rolesanywhere.GetProfileInput{
		ProfileId: aws.String(profile_id),
	}
	op, err := svc.GetProfile(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_profile.getProfile", "api_error", err)
		return nil, err
	}
	return *op.Profile, nil
}

func listTagsForProfile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	profile_arn := h.Item.(types.ProfileDetail).ProfileArn

	svc, err := RolesAnywhereClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_profile.listTagsForProfile", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	params := &rolesanywhere.ListTagsForResourceInput{
		ResourceArn: profile_arn,
	}
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_profile.listTagsForProfile", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func profileTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*rolesanywhere.ListTagsForResourceOutput)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
