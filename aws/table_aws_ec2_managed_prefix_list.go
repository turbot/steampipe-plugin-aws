package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2ManagedPrefixList(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_managed_prefix_list",
		Description: "AWS EC2 Managed Prefix List",
		List: &plugin.ListConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidAction", "InvalidRequest"}),
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
				Transform:   transform.FromField("Tags").Transform(handlePrefixListTagsEmptyResult),
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	logger.Trace("listManagedPrefixList", "AWS_REGION", region)

	equalQuals := d.KeyColumnQuals
	filters := []*ec2.Filter{}

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		logger.Error("listManagedPrefixList", "Ec2Service_error", err)
		return nil, err
	}

	params := &ec2.DescribeManagedPrefixListsInput{
		MaxResults: aws.Int64(100),
	}

	if equalQuals["owner_id"] != nil {
		ownerIdFilter := ec2.Filter{
			Name:   aws.String("owner-id"),
			Values: []*string{aws.String(equalQuals["owner_id"].GetStringValue())},
		}
		filters = append(filters, &ownerIdFilter)
	}

	if equalQuals["id"] != nil {
		idFilter := ec2.Filter{
			Name:   aws.String("prefix-list-id"),
			Values: []*string{aws.String(equalQuals["id"].GetStringValue())},
		}
		filters = append(filters, &idFilter)
	}

	if equalQuals["name"] != nil {
		nameFilter := ec2.Filter{
			Name:   aws.String("prefix-list-name"),
			Values: []*string{aws.String(equalQuals["name"].GetStringValue())},
		}
		filters = append(filters, &nameFilter)
	}

	// Add filters as request parameter when at least one filter is present
	if len(filters) > 0 {
		params.Filters = filters
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			if *limit < 1 {
				params.MaxResults = aws.Int64(1)
			} else {
				params.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribeManagedPrefixListsPages(
		params,
		func(page *ec2.DescribeManagedPrefixListsOutput, isLast bool) bool {
			for _, prefix := range page.PrefixLists {
				d.StreamListItem(ctx, prefix)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		logger.Error("listManagedPrefixList", "DescribeManagedPrefixListsPages_error", err)
		return nil, err
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func prefixListTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	prefixList := d.HydrateItem.(*ec2.ManagedPrefixList)

	var turbotTagsMap map[string]string

	if len(prefixList.Tags) < 1 {
		return nil, nil
	}

	if prefixList.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range prefixList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}

func handlePrefixListTagsEmptyResult(_ context.Context, d *transform.TransformData) (interface{}, error) {
	prefixList := d.HydrateItem.(*ec2.ManagedPrefixList)
	if len(prefixList.Tags) > 0  {
		return prefixList.Tags, nil
	}
	return nil, nil
}
