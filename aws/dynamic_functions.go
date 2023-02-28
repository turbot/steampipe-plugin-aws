package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func listDynamicS3Buckets(ctx context.Context, d *plugin.QueryData) ([]types.Bucket, error) {

	// have we already created and cached the session?
	//cacheKey := "listDynamicS3Bucket"

	// if cachedData, ok := cn.Get(ctx, cacheKey); ok {
	// 	return cachedData.(.......), nil
	// }

	// Unlike most services, S3 buckets are a global list. They can be retrieved
	// from any single region.  We must list buckets from the default region to
	// get the actual creation_time of the bucket, in all other regions the list
	// returns the time when the bucket was last modified. See
	// https://www.marksayson.com/blog/s3-bucket-creation-dates-s3-master-regions/
	defaultRegion, err := getLastResortRegion(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	svc, err := S3Client(ctx, d, defaultRegion)
	if err != nil {
		plugin.Logger(ctx).Error("listDynamicS3Buckets", "get_client_error", err, "defaultRegion", defaultRegion)
		return nil, err
	}

	// execute list call
	input := &s3.ListBucketsInput{}
	bucketsResult, err := svc.ListBuckets(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("listDynamicS3Buckets", "api_error", err, "defaultRegion", defaultRegion)
		return nil, err
	}

	return bucketsResult.Buckets, nil
}
