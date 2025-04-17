package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rolesanywhere"

	rolesanywherev1 "github.com/aws/aws-sdk-go/service/rolesanywhere"

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
			KeyColumns: plugin.AnyColumn([]string{"id"}),
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
		GetMatrixItemFunc: SupportedRegionMatrix(rolesanywherev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the profile.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProfileId"),
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
				Description: "Accept custom role session names.",
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
				Name:        "duration",
				Description: "The profile credential session duration in seconds.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DurationSeconds"),
			},
			{
				Name:        "enabled",
				Description: "If the profile is enabled or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "attribute_mappings",
				Description: "List of attribute mappings for certificate fields to session tags.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "managed_policy_arns",
				Description: "List of managed IAM boundary policy ARNs attached to the profile.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "role_arns",
				Description: "List of IAM role ARNs the profile can assume.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "require_instance_properties",
				Description: "If instance properties are required for the session creation.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "session_policy",
				Description: "Session policy applied to created sessions.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SessionPolicy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "session_policy_std",
				Description: "Contains the session policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SessionPolicy").Transform(unescape).Transform(policyToCanonical),
			},
		}),
	}
}

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

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil

}
