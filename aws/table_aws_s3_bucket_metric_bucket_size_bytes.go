package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsS3BucketMetricBucketSizeBytes(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_bucket_metric_bucket_size_bytes",
		Description: "AWS S3 Bucket CloudWatch Metrics - BucketSizeBytes",
		List: &plugin.ListConfig{
			ParentHydrate: listS3Buckets,
			Hydrate:       listS3BucketMetricBucketSizeBytes,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns(cwMetricColumns(
			[]*plugin.Column{
				{
					Name:        "name",
					Description: "The user friendly name of the bucket.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listS3BucketMetricBucketSizeBytes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*s3.Bucket)
	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("StorageType"),
			Value: aws.String("StandardStorage"),
		},
		{
			Name:  aws.String("BucketName"),
			Value: data.Name,
		},
	}

	return listCWMetricStatistics(ctx, d, "5_MIN", "AWS/S3", "BucketSizeBytes", dimensions)
}
