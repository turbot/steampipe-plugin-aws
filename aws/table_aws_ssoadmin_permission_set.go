package aws

import (
	"context"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"

	ssoadminv1 "github.com/aws/aws-sdk-go/service/ssoadmin"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSsoAdminPermissionSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssoadmin_permission_set",
		Description: "AWS SSO Permission Set",
		List: &plugin.ListConfig{
			ParentHydrate: listSsoAdminInstances,
			Hydrate:       listSsoAdminPermissionSets,
			Tags:          map[string]string{"service": "sso", "action": "ListPermissionSets"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "instance_arn", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getSsoAdminPermissionSet,
				Tags: map[string]string{"service": "sso", "action": "DescribePermissionSet"},
			},
			{
				Func: getSsoAdminPermissionSetTags,
				Tags: map[string]string{"service": "sso", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssoadminv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the permission set.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSsoAdminPermissionSet,
				Transform:   transform.FromField("PermissionSet.Name"),
			},
			{
				Name:        "arn",
				Description: "The ARN of the permission set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PermissionSetArn"),
			},
			{
				Name:        "created_date",
				Description: "The date that the permission set was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getSsoAdminPermissionSet,
				Transform:   transform.FromField("PermissionSet.CreatedDate"),
			},
			{
				Name:        "description",
				Description: "The description of the permission set.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSsoAdminPermissionSet,
				Transform:   transform.FromField("PermissionSet.Description"),
			},
			{
				Name:        "instance_arn",
				Description: "The Amazon Resource Name (ARN) of the SSO Instance under which the operation will be executed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "relay_state",
				Description: "Used to redirect users within the application during the federation authentication process.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSsoAdminPermissionSet,
				Transform:   transform.FromField("PermissionSet.RelayState"),
			},
			{
				Name:        "session_duration",
				Description: "The length of time that the application user sessions are valid for in the ISO-8601 standard.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSsoAdminPermissionSet,
				Transform:   transform.FromField("PermissionSet.SessionDuration"),
			},
			{
				Name:      "tags_src",
				Type:      proto.ColumnType_JSON,
				Hydrate:   getSsoAdminPermissionSetTags,
				Transform: transform.FromValue(),
			},
			{
				Name:        "inline_policy",
				Description: "The policy document embedded inline for the permission set.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSsoAdminPermissionSetInlinePolicy,
				Transform:   transform.FromValue().Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "inline_policy_std",
				Description: "Inline policy in canonical form for the permission set.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSsoAdminPermissionSetInlinePolicy,
				Transform:   transform.FromValue().Transform(unescape).Transform(policyToCanonical),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSsoAdminPermissionSetTags,
				Transform:   transform.From(getSsoAdminResourceTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSsoAdminPermissionSet,
				Transform:   transform.FromField("PermissionSet.Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PermissionSetArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSsoAdminPermissionSets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(types.InstanceMetadata)
	instanceArn := *instance.InstanceArn

	// Create session
	svc, err := SSOAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_permission_set.listSsoAdminPermissionSets", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	equalQuals := d.EqualsQuals
	// Minimize the API call with the given layer name
	if equalQuals["instance_arn"] != nil {
		if equalQuals["instance_arn"].GetStringValue() != "" {
			if equalQuals["instance_arn"].GetStringValue() != "" && equalQuals["instance_arn"].GetStringValue() != instanceArn {
				return nil, nil
			}
		} else if len(getListValues(equalQuals["instance_arn"].GetListValue())) > 0 {
			if !slices.Contains(aws.ToStringSlice(getListValues(equalQuals["instance_arn"].GetListValue())), instanceArn) {
				return nil, nil
			}
		}
	}

	// Limiting the results
	maxLimit := int32(100)
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

	input := &ssoadmin.ListPermissionSetsInput{
		InstanceArn: aws.String(instanceArn),
		MaxResults:  aws.Int32(maxLimit),
	}

	paginator := ssoadmin.NewListPermissionSetsPaginator(svc, input, func(o *ssoadmin.ListPermissionSetsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssoadmin_permission_set.listSsoAdminPermissionSets", "api_error", err)
			return nil, err
		}

		for _, items := range output.PermissionSets {
			d.StreamListItem(ctx, &PermissionSetItem{
				InstanceArn:      aws.String(instanceArn),
				PermissionSetArn: aws.String(items),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

type PermissionSetItem struct {
	InstanceArn      *string
	PermissionSetArn *string
	PermissionSet    types.PermissionSet
}

//// HYDRATE FUNCTIONS

func getSsoAdminPermissionSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SSOAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_permission_set.getSsoAdminPermissionSet", "connection_error", err)

		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	permissionSet := h.Item.(*PermissionSetItem)
	arn := *permissionSet.PermissionSetArn
	instanceArn := *permissionSet.InstanceArn

	params := &ssoadmin.DescribePermissionSetInput{
		InstanceArn:      aws.String(instanceArn),
		PermissionSetArn: aws.String(arn),
	}

	detail, err := svc.DescribePermissionSet(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_permission_set.getSsoAdminPermissionSet", "api_error", err)
		return nil, err
	}

	item := &PermissionSetItem{
		InstanceArn:      &instanceArn,
		PermissionSetArn: &arn,
		PermissionSet:    *detail.PermissionSet,
	}
	return item, nil
}

func getSsoAdminPermissionSetTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	permissionSet := h.Item.(*PermissionSetItem)
	resourceArn := *permissionSet.PermissionSetArn
	instanceArn := *permissionSet.InstanceArn

	tags, err := getSsoAdminResourceTags(ctx, d, instanceArn, resourceArn)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_permission_set.getSsoAdminPermissionSetTags", "api_error", err)
	}
	return tags, err
}

func getSsoAdminPermissionSetInlinePolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SSOAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_permission_set.getSsoAdminPermissionSetInlinePolicy", "connection_error", err)

		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	permissionSet := h.Item.(*PermissionSetItem)
	permissionSetArn := *permissionSet.PermissionSetArn
	instanceArn := *permissionSet.InstanceArn

	params := &ssoadmin.GetInlinePolicyForPermissionSetInput{
		InstanceArn:      aws.String(instanceArn),
		PermissionSetArn: aws.String(permissionSetArn),
	}

	resp, err := svc.GetInlinePolicyForPermissionSet(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_permission_set.getSsoAdminPermissionSetInlinePolicy", "api_error", err)
		return nil, err
	}

	return resp.InlinePolicy, nil
}

func getSsoAdminResourceTags(ctx context.Context, d *plugin.QueryData, instanceArn, resourceArn string) (interface{}, error) {
	// Create session
	svc, err := SSOAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_permission_set.getSsoAdminResourceTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &ssoadmin.ListTagsForResourceInput{
		InstanceArn: aws.String(instanceArn),
		ResourceArn: aws.String(resourceArn),
	}

	tags := []types.Tag{}

	paginator := ssoadmin.NewListTagsForResourcePaginator(svc, params, func(o *ssoadmin.ListTagsForResourcePaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssoadmin_permission_set.getSsoAdminResourceTags", "api_error", err)
			return nil, err
		}

		tags = append(tags, output.Tags...)
	}

	return &tags, err
}

//// TRANSFORM FUNCTIONS

func getSsoAdminResourceTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*[]types.Tag)
	tagsMap := map[string]string{}

	for _, tag := range *tags {
		tagsMap[*tag.Key] = *tag.Value
	}

	return &tagsMap, nil
}
