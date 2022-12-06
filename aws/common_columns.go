package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

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
			Hydrate:     getCommonColumnsCached,
			Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
		},
		{
			Name:        "region",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumnsCached,
			Description: "The AWS Region in which the resource is located.",
		},
		{
			Name:        "account_id",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumnsCached,
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
			Hydrate:     getCommonColumnsCached,
			Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
		},
		{
			Name:        "account_id",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumnsCached,
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
			Hydrate:     getCommonColumnsCached,
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
			Hydrate:     getCommonColumnsCached,
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

// build a cache key for the call to getCommonColumns, including the region since this is a multi-region call
func getCommonColumnsCacheKey() plugin.HydrateFunc {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		region := d.KeyColumnQualString(matrixKeyRegion)
		key := fmt.Sprintf("getCommonColumns-%s", region)
		return key, nil
	}
}

// if the caching is required other than per connection, build a cache key for the call and use it in WithCache
// since getCommonColumns is a multi-region call, caching should be per connection per region
var getCommonColumnsCached = plugin.HydrateFunc(getCommonColumnsUncached).WithCache(getCommonColumnsCacheKey())

// get columns which are returned with all tables: region, partition and account
func getCommonColumnsUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "" {
		region = "global"
	}

	// Trace logging to debug cache and execution flows
	plugin.Logger(ctx).Trace("getCommonColumnsUncached", "status", "starting", "connection_name", d.Connection.Name, "region", region)

	// use the cached version of the getCallerIdentity to reduce the number of request
	var commonColumnData *awsCommonColumnData
	getCallerIdentityData, err := getCallerIdentity(ctx, d, h)
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

// define cached version of getCallerIdentity and getCommonColumns
// by default, WithCache cached the data per connection
// if no argument is passed in WithCache, the cache key will be in the format of <function_name>-<connection_name>
var getCallerIdentity = plugin.HydrateFunc(getCallerIdentityUncached).WithCache()

// returns details about the IAM user or role whose credentials are used to call the operation
func getCallerIdentityUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Trace logging to debug cache and execution flows
	plugin.Logger(ctx).Trace("getCallerIdentityUncached", "status", "starting", "connection_name", d.Connection.Name)

	// get the service connection for the service
	svc, err := STSClient(ctx, d)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	callerIdentity, err := svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	plugin.Logger(ctx).Trace("getCallerIdentityUncached", "STS.GetCallerIdentity response time", time.Since(now), "connection_name", d.Connection.Name)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("getCallerIdentityUncached", "status", "finished", "connection_name", d.Connection.Name)
	return callerIdentity, nil
}
