package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableEc2AwsManagedPrefixList(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_managed_prefix_list",
		Description: "AWS EC2 Managed Prefix List",
		List: &plugin.ListConfig{
			Hydrate: 		   listManagedPrefixList,
			ShouldIgnoreError: isNotFoundError([]string{"InvalidAction"}),
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "prefix_list_name",
				Description: "The name of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "prefix_list_id",
				Description: "The ID of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "prefix_list_arn",
				Description: "The Amazon Resource Name (ARN) for the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "address_family",
				Description: "The IP address version.",
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
				Description: "The state message.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The state message.",
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	logger.Trace("listManagedPrefixList", "AWS_REGION", region)

	ownerIds := d.KeyColumnQuals["owner_id"].GetStringValue()
	prefixListIds := d.KeyColumnQuals["prefix_list_id"].GetStringValue()
	prefixListNames := d.KeyColumnQuals["prefix_list_name"].GetStringValue()

	// Create Session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		logger.Error("listManagedPrefixList", "Ec2Service_error", err)
		return nil, err
	}

	params := &ec2.DescribeManagedPrefixListsInput{
		MaxResults: aws.Int64(1000),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("owner-id"),
				Values: []*string{aws.String(ownerIds)},
			},
			{
				Name:   aws.String("prefix-list-id"),
				Values: []*string{aws.String(prefixListIds)},
			},
			{
				Name:   aws.String("prefix-list-name"),
				Values: []*string{aws.String(prefixListNames)},
			},
		},
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			params.MaxResults = limit
		}
	}

	// List call
	err = svc.DescribeManagedPrefixListsPages(
		&ec2.DescribeManagedPrefixListsInput{},
		func(page *ec2.DescribeManagedPrefixListsOutput, isLast bool) bool {
			for _, prefixList := range page.PrefixLists {
				d.StreamListItem(ctx, prefixList)
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
	if prefixList.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range prefixList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
