package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

//// TABLE DEFINITION

func tableAwsEc2ManagedPrefixList(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_managed_prefix_list",
		Description: "AWS EC2 Managed Prefix List",
		List: &plugin.ListConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAction", "InvalidRequest", "UnsupportedOperation"}),
			},
			Hydrate: listManagedPrefixList,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "id", Require: plugin.Optional},
				{Name: "owner_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the prefix list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrefixListName"),
			},
			{
				Name:        "id",
				Description: "The ID of the prefix list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrefixListId"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the prefix list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrefixListArn"),
			},
			{
				Name:        "state",
				Description: "The current state of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "address_family",
				Description: "The IP address version of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_entries",
				Description: "The maximum number of entries for the prefix list.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the owner of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_message",
				Description: "The message regarding the current state of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The version of the prefix list.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tags_src",
				Description: "The tags for the prefix list.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PrefixListName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(prefixListTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrefixListArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listManagedPrefixList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	equalQuals := d.KeyColumnQuals
	filters := []types.Filter{}

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		logger.Error("aws_ec2_managed_prefix_list.listManagedPrefixList", "connection_error", err)
		return nil, err
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

	params := &ec2.DescribeManagedPrefixListsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if equalQuals["owner_id"] != nil {
		ownerIdFilter := types.Filter{
			Name:   aws.String("owner-id"),
			Values: []string{equalQuals["owner_id"].GetStringValue()},
		}
		filters = append(filters, ownerIdFilter)
	}

	if equalQuals["id"] != nil {
		idFilter := types.Filter{
			Name:   aws.String("prefix-list-id"),
			Values: []string{equalQuals["id"].GetStringValue()},
		}
		filters = append(filters, idFilter)
	}

	if equalQuals["name"] != nil {
		nameFilter := types.Filter{
			Name:   aws.String("prefix-list-name"),
			Values: []string{equalQuals["name"].GetStringValue()},
		}
		filters = append(filters, nameFilter)
	}

	// Add filters as request parameter when at least one filter is present
	if len(filters) > 0 {
		params.Filters = filters
	}

	paginator := ec2.NewDescribeManagedPrefixListsPaginator(svc, params, func(o *ec2.DescribeManagedPrefixListsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_managed_prefix_list.listManagedPrefixList", "api_error", err)
			return nil, err
		}

		for _, items := range output.PrefixLists {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func prefixListTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	prefixList := d.HydrateItem.(types.ManagedPrefixList)

	var turbotTagsMap map[string]string
	if prefixList.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range prefixList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
