package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// column definitions for the common columns
var commonAwsRegionalColumns = []*plugin.Column{
	{
		Name:        "partition",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
		Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov)",
	},
	{
		Name:        "region",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
		Description: "The AWS Region in which the resource is located",
	},
	{
		Name:        "account_id",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
		Description: "The AWS Account ID in which the resource is located",
	},
}

// column definitions for the common columns
var commonS3Columns = []*plugin.Column{
	{
		Name:        "partition",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
		Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov)",
	},
	{
		Name:        "account_id",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
		Description: "The AWS Account ID in which the resource is located",
	},
}

var commonAwsColumns = []*plugin.Column{
	{
		Name:        "partition",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
		Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov)",
	},
	{
		Name:        "region",
		Type:        proto.ColumnType_STRING,
		Transform:   transform.FromConstant("global"),
		Description: "The AWS Region in which the resource is located",
	},
	{
		Name:        "account_id",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
		Description: "The AWS Account ID in which the resource is located",
	},
}

// append the common aws columns for REGIONAL resources onto the column list
func awsRegionalColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonAwsRegionalColumns...)
}

// append the common aws columns for GLOBAL resources onto the column list
func awsColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonAwsColumns...)
}

func awsS3Columns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonS3Columns...)
}

// struct to store the common column data
type awsCommonColumnData struct {
	Partition, Region, AccountId string
}

// get columns which are returned with all tables: region, partition and account
func getCommonColumns(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var region string
	if plugin.GetFetchMetadata(ctx)[fetchMetdataKeyRegion] != nil {
		region = plugin.GetFetchMetadata(ctx)[fetchMetdataKeyRegion].(string)
	}
	if region == "" {
		region = "global"
	}
	plugin.Logger(ctx).Trace("getCommonColumns", "region", region)

	cacheKey := "commonColumnData" + region
	var commonColumnData *awsCommonColumnData
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		commonColumnData = cachedData.(*awsCommonColumnData)
	} else {
		stsSvc, err := StsService(ctx, d.ConnectionManager)
		if err != nil {
			return nil, err
		}

		callerIdentity, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
		if err != nil {
			return nil, err
		}
		commonColumnData = &awsCommonColumnData{
			// extract partition from arn
			Partition: strings.Split(*callerIdentity.Arn, ":")[1],
			AccountId: *callerIdentity.Account,
			Region:    region,
		}

		// save to extension cache
		d.ConnectionManager.Cache.Set(cacheKey, commonColumnData)
	}

	plugin.Logger(ctx).Trace("getCommonColumns: ", "commonColumnData", commonColumnData)

	return commonColumnData, nil
}
