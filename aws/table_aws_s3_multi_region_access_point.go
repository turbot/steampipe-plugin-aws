package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/s3control/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsS3MultiRegionAccessPoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_multi_region_access_point",
		Description: "AWS S3 Multi Region Access Point",
		List: &plugin.ListConfig{
			Hydrate: listS3MultiRegionAccessPoints,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameter", "InvalidRequest"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "account_id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "account_id"}),
			Hydrate:    getS3MultiRegionAccessPoint,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchMultiRegionAccessPoint", "InvalidParameter", "InvalidRequest"}),
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Multi-Region Access Point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "alias",
				Description: "The alias for the Multi-Region Access Point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "When the Multi-Region Access Point create request was received.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "The current status of the Multi-Region Access Point. CREATING and DELETING are temporary states that exist while the request is propagating and being completed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_access_block",
				Description: "The PublicAccessBlock configuration that you want to apply to this Amazon S3 account.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "regions",
				Description: "A collection of the Regions and buckets associated with the Multi-Region Access Point.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMultiRegionAccessPointArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listS3MultiRegionAccessPoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get account details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_multi_region_access_point.listS3MultiRegionAccessPoints", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	region := d.KeyColumnQualString(matrixKeyRegion)
	// Create Session
	svc, err := S3ControlMultiRegionAccessClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_multi_region_access_point.listS3MultiRegionAccessPoints", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	ownerAccountId := d.KeyColumnQuals["account_id"].GetStringValue()
	if ownerAccountId != "" && ownerAccountId != commonColumnData.AccountId {
		return nil, nil
	}

	maxItems := int32(100)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	input := &s3control.ListMultiRegionAccessPointsInput{
		AccountId:  aws.String(commonColumnData.AccountId),
		MaxResults: maxItems,
	}

	paginator := s3control.NewListMultiRegionAccessPointsPaginator(svc, input, func(o *s3control.ListMultiRegionAccessPointsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3_multi_region_access_point.listS3MultiRegionAccessPoints", "api_error", err)
			return nil, err
		}

		for _, accessPoint := range output.AccessPoints {
			d.StreamListItem(ctx, accessPoint)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getS3MultiRegionAccessPoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_multi_region_access_point.getS3MultiRegionAccessPoint", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := S3ControlMultiRegionAccessClient(ctx, d, matrixRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_multi_region_access_point.getS3MultiRegionAccessPoint", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	var name, ownerAccountId string
	if h.Item != nil {
		multiRegionAccessPoint := h.Item.(types.MultiRegionAccessPointReport)
		name = *multiRegionAccessPoint.Name
		ownerAccountId = commonColumnData.AccountId
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		ownerAccountId = d.KeyColumnQuals["account_id"].GetStringValue()
	}

	// Build params
	params := &s3control.GetMultiRegionAccessPointInput{
		Name:      aws.String(name),
		AccountId: aws.String(ownerAccountId),
	}

	// execute list call
	item, err := svc.GetMultiRegionAccessPoint(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_multi_region_access_point.getS3MultiRegionAccessPoint", "api_error", err)
		return nil, err
	}

	return item.AccessPoint, nil
}

func getMultiRegionAccessPointArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accessPointName := multiRegionAccessPointName(h.Item)

	// Get account details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_multi_region_access_point.getMultiRegionAccessPointArn", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":s3::" + commonColumnData.AccountId + ":accesspoint/" + accessPointName

	return arn, nil
}

func multiRegionAccessPointName(item interface{}) string {
	switch item := item.(type) {
	case types.MultiRegionAccessPointReport:
		return *item.Alias
	case *types.MultiRegionAccessPointReport:
		return *item.Name
	}
	return ""
}
