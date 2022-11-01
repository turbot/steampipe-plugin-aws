package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

// column definitions for the common columns
func commonAwsRegionalColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "partition",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
		},
		{
			Name:        "region",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Description: "The AWS Region in which the resource is located.",
		},
		{
			Name:        "account_id",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Description: "The AWS Account ID in which the resource is located.",
			Transform:   transform.FromCamel(),
		},
	}
}

// column definitions for the common columns
func commonColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "partition",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
		},
		{
			Name:        "account_id",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Transform:   transform.FromCamel(),
			Description: "The AWS Account ID in which the resource is located.",
		},
	}
}

func commonAwsColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "partition",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
		},
		{
			Name:        "region",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromConstant("global"),
			Description: "The AWS Region in which the resource is located.",
		},
		{
			Name:        "account_id",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Description: "The AWS Account ID in which the resource is located.",
			Transform:   transform.FromCamel(),
		},
	}
}

// append the common aws columns for REGIONAL resources onto the column list
func awsRegionalColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonAwsRegionalColumns()...)
}

// append the common aws columns for GLOBAL resources onto the column list
func awsColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonAwsColumns()...)
}

func awsDefaultColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonColumns()...)
}

// struct to store the common column data
type awsCommonColumnData struct {
	Partition, Region, AccountId string
}

// get columns which are returned with all tables: region, partition and account
func getCommonColumns(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "" {
		region = "global"
	}

	var commonColumnData *awsCommonColumnData
	getCallerIdentityCached := plugin.HydrateFunc(getCallerIdentity).WithCache()
	getCallerIdentityData, err := getCallerIdentityCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	callerIdentity := getCallerIdentityData.(*sts.GetCallerIdentityOutput)
	commonColumnData = &awsCommonColumnData{
		// extract partition from arn
		Partition: strings.Split(*callerIdentity.Arn, ":")[1],
		AccountId: *callerIdentity.Account,
		Region:    region,
	}

	return commonColumnData, nil
}

func getCallerIdentity(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cacheKey := "GetCallerIdentity"

	// if found in cache, return the result
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*sts.GetCallerIdentityOutput), nil
	}

	// get the service connection for the service
	svc, err := STSClient(ctx, d)
	if err != nil {
		return nil, err
	}

	callerIdentity, err := svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		// let the cache know that we have failed to fetch this item
		return nil, err
	}

	// save to extension cache
	d.ConnectionManager.Cache.Set(cacheKey, callerIdentity)
	return callerIdentity, nil
}

func getAccountPartition(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheKey := "getAccountPartition"

	// if found in cache, return the result
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(string), nil
	}

	commonData, err := getCommonColumns(ctx, d, nil)
	if err != nil {
		plugin.Logger(ctx).Error("getAccountPartition", "common_data_error", err)
		// If error or some other issue return default partition(i.e. AWS commercial)
		return "aws", nil
	}

	// save to cache
	d.ConnectionManager.Cache.Set(cacheKey, commonData.(*awsCommonColumnData).Partition)
	return commonData.(*awsCommonColumnData).Partition, nil
}
