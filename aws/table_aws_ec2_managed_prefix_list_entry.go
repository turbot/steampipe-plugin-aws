package aws

import (
	"context"

	ec2Endpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

//// TABLE DEFINITION

func tableAwsEc2ManagedPrefixListEntry(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_managed_prefix_list_entry",
		Description: "AWS EC2 Managed Prefix List Entry",
		List: &plugin.ListConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				// Ignore the 'UnsupportedOperation' error because while the EC2 service is supported in the 'me-south-1' region, listing managed prefix lists is not in that region
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAction", "InvalidRequest", "UnsupportedOperation"}),
			},
			ParentHydrate: listManagedPrefixList,
			Hydrate:       listManagedPrefixListEntries,
			Tags:          map[string]string{"service": "ec2", "action": "GetManagedPrefixListEntries"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "prefix_list_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2Endpoint.EC2ServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "prefix_list_id",
				Description: "The ID of the prefix list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cidr",
				Description: "The CIDR block.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the entry.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Cidr"),
			},
		}),
	}
}

type PrefixListEntryInfo struct {
	PrefixListId *string
	Cidr         *string
	Description  *string
}

//// LIST FUNCTION

func listManagedPrefixListEntries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	prefixList := h.Item.(types.ManagedPrefixList)

	if d.EqualsQualString("prefix_list_id") != "" && d.EqualsQualString("prefix_list_id") != *prefixList.PrefixListId {
		return nil, nil
	}

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		logger.Error("aws_ec2_managed_prefix_list_entry.listManagedPrefixListEntries", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	params := &ec2.GetManagedPrefixListEntriesInput{
		MaxResults:   aws.Int32(maxLimit),
		PrefixListId: prefixList.PrefixListId,
	}

	paginator := ec2.NewGetManagedPrefixListEntriesPaginator(svc, params, func(o *ec2.GetManagedPrefixListEntriesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_managed_prefix_list_entry.listManagedPrefixListEntries", "api_error", err)
			return nil, err
		}

		for _, item := range output.Entries {
			d.StreamListItem(ctx, &PrefixListEntryInfo{
				PrefixListId: prefixList.PrefixListId,
				Cidr:         item.Cidr,
				Description:  item.Description,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, nil
}
