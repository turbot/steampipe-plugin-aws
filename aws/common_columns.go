package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

// Columns defined on every account-level resource (e.g. aws_iam_access_key)
func commonColumnsForAccountResource() []*plugin.Column {
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

// Columns defined on every region-level resource (e.g. aws_ec2_instance)
func commonColumnsForRegionalResource() []*plugin.Column {
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

// Columns defined on every global-region-level resource (e.g. aws_waf_rule)
func commonColumnsForGlobalRegionResource() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "partition",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
		},
		{
			Name: "region",
			Type: proto.ColumnType_STRING,
			// Region is hard-coded to special global region
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

// Append columns for account-level resource (e.g. aws_iam_access_key)
func awsRegionalColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonColumnsForRegionalResource()...)
}

// Append columns for region-level resource (e.g. aws_ec2_instance)
func awsGlobalRegionColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonColumnsForGlobalRegionResource()...)
}

// Append columns for global-region-level resource (e.g. aws_waf_rule)
func awsAccountColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonColumnsForAccountResource()...)
}

// struct to store the common column data
type awsCommonColumnData struct {
	Partition, Region, AccountId string
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize
// since getCommonColumns is a multi-region call, caching should be per connection per region
var getCommonColumns = plugin.HydrateFunc(getCommonColumnsUncached).WithCache(getCommonColumnsCacheKey)

// Build a cache key for the call to getCommonColumns, including the region since this is a multi-region call.
// Notably, this may be called WITHOUT a region. In that case we just share a cache for non-region data.
func getCommonColumnsCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	key := fmt.Sprintf("getCommonColumns-%s", region)
	return key, nil
}

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
		plugin.Logger(ctx).Error("getCommonColumnsUncached", "status", "failed", "connection_name", d.Connection.Name, "region", region, "error", err)
		return nil, err
	}

	callerIdentity := getCallerIdentityData.(*sts.GetCallerIdentityOutput)
	commonColumnData = &awsCommonColumnData{
		// extract partition from arn
		Partition: strings.Split(*callerIdentity.Arn, ":")[1],
		AccountId: *callerIdentity.Account,
		Region:    region,
	}

	plugin.Logger(ctx).Trace("getCommonColumnsUncached", "status", "starting", "connection_name", d.Connection.Name, "common_column_data", *commonColumnData)

	return commonColumnData, nil
}

// define cached version of getCallerIdentity and getCommonColumns
// by default, Memoize cached the data per connection
// if no argument is passed in Memoize, the cache key will be in the format of <function_name>-<connection_name>
var getCallerIdentity = plugin.HydrateFunc(getCallerIdentityUncached).WithCache()

// returns details about the IAM user or role whose credentials are used to call the operation
func getCallerIdentityUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Trace logging to debug cache and execution flows
	plugin.Logger(ctx).Trace("getCallerIdentityUncached", "status", "starting", "connection_name", d.Connection.Name)

	// get the service connection for the service
	svc, err := STSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getCallerIdentityUncached", "status", "failed", "connection_name", d.Connection.Name, "client_error", err)
		return nil, err
	}

	callerIdentity, err := svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		plugin.Logger(ctx).Error("getCallerIdentityUncached", "status", "failed", "connection_name", d.Connection.Name, "api_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Trace("getCallerIdentityUncached", "status", "finished", "connection_name", d.Connection.Name)
	return callerIdentity, nil
}
