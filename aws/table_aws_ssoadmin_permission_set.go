package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsSsoAdminPermissionSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssoadmin_permission_set",
		Description: "AWS SSO Permission Set",
		List: &plugin.ListConfig{
			ParentHydrate: listSsoAdminInstances,
			Hydrate:       listSsoAdminPermissionSets,
		},
		GetMatrixItem: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listSsoAdminPermissionSets")

	instance := h.Item.(*ssoadmin.InstanceMetadata)
	instanceArn := *instance.InstanceArn

	// Create session
	svc, err := SSOAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	err = svc.ListPermissionSetsPages(
		&ssoadmin.ListPermissionSetsInput{
			InstanceArn: aws.String(instanceArn),
		},
		func(page *ssoadmin.ListPermissionSetsOutput, isLast bool) bool {
			for _, arn := range page.PermissionSets {
				item := &PermissionSetItem{
					InstanceArn:      &instanceArn,
					PermissionSetArn: arn,
				}
				d.StreamListItem(ctx, item)
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listSsoAdminPermissionSets", "ListPermissionSetsPages_error", err)
	}

	return nil, err
}

type PermissionSetItem struct {
	InstanceArn      *string
	PermissionSetArn *string
	PermissionSet    *ssoadmin.PermissionSet
}

//// HYDRATE FUNCTIONS

func getSsoAdminPermissionSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSsoAdminPermissionSet")

	// Create session
	svc, err := SSOAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	permissionSet := h.Item.(*PermissionSetItem)
	arn := *permissionSet.PermissionSetArn
	instanceArn := *permissionSet.InstanceArn

	params := &ssoadmin.DescribePermissionSetInput{
		InstanceArn:      aws.String(instanceArn),
		PermissionSetArn: aws.String(arn),
	}
	plugin.Logger(ctx).Trace("getSsoAdminPermissionSet", "DescribePermissionSet_input", params)

	detail, err := svc.DescribePermissionSet(params)
	if err != nil {
		plugin.Logger(ctx).Error("getSsoAdminPermissionSet", "DescribePermissionSet_error", err)
		return nil, err
	}

	item := &PermissionSetItem{
		InstanceArn:      &instanceArn,
		PermissionSetArn: &arn,
		PermissionSet:    detail.PermissionSet,
	}
	return item, nil
}

func getSsoAdminPermissionSetTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSsoAdminPermissionSetTags")

	permissionSet := h.Item.(*PermissionSetItem)
	resourceArn := *permissionSet.PermissionSetArn
	instanceArn := *permissionSet.InstanceArn

	tags, err := getSsoAdminResourceTags(ctx, d, instanceArn, resourceArn)
	return tags, err
}

func getSsoAdminResourceTags(ctx context.Context, d *plugin.QueryData, instanceArn, resourceArn string) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSsoAdminResourceTags")

	// Create session
	svc, err := SSOAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &ssoadmin.ListTagsForResourceInput{
		InstanceArn: aws.String(instanceArn),
		ResourceArn: aws.String(resourceArn),
	}

	tags := []*ssoadmin.Tag{}

	err = svc.ListTagsForResourcePages(
		params,
		func(page *ssoadmin.ListTagsForResourceOutput, isLast bool) bool {
			tags = append(tags, page.Tags...)
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("getSsoAdminResourceTags", "ListTagsForResourcePages_error", err)
		return nil, err
	}

	return &tags, err
}

//// TRANSFORM FUNCTIONS

func getSsoAdminResourceTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSsoAdminResourceTurbotTags")

	tags := d.HydrateItem.(*[]*ssoadmin.Tag)
	tagsMap := map[string]string{}

	for _, tag := range *tags {
		tagsMap[*tag.Key] = *tag.Value
	}

	return &tagsMap, nil
}
