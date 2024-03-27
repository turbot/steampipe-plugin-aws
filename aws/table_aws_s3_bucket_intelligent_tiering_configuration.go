package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsS3BucketIntelligentTieringConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_bucket_intelligent_tiering_configuration",
		Description: "AWS S3 Bucket Intelligent Tiering Configuration",
		Get: &plugin.GetConfig{
			Hydrate:    getBucketIntelligentTieringConfiguration,
			Tags:       map[string]string{"service": "s3", "action": "GetIntelligentTieringConfiguration"},
			KeyColumns: plugin.AllColumns([]string{"bucket_name", "id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchConfiguration"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listS3Buckets,
			Hydrate:       listBucketIntelligentTieringConfigurations,
			Tags:          map[string]string{"service": "s3", "action": "ListBucketIntelligentTieringConfigurations"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket_name", Require: plugin.Optional, CacheMatch: "exact"},
			},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "bucket_name",
				Description: "The name of the container bucket of this object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID used to identify the S3 Intelligent-Tiering configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Specifies the status of the configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tierings",
				Description: "Specifies the S3 Intelligent-Tiering storage class tier of the configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "filter",
				Description: "Specifies a bucket filter. The configuration only includes objects that meet the filter's criteria.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
		}),
	}
}

type IntelligentTieringConfigurationInfo struct {
	BucketName *string
	types.IntelligentTieringConfiguration
}

//// LIST FUNCTION

func listBucketIntelligentTieringConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucket := h.Item.(types.Bucket)
	if d.EqualsQualString("bucket_name") != "" && d.EqualsQualString("bucket_name") != *bucket.Name {
		return nil, nil
	}

	bucketRegion, err := doGetBucketRegion(ctx, d, h, *bucket.Name)
	if err != nil {
		return nil, err
	}

	// Create client
	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket_intelligent_tiering_configuration.listBucketIntelligentTieringConfigurations", "client_error", err)
		return nil, err
	}

	params := &s3.ListBucketIntelligentTieringConfigurationsInput{
		Bucket: bucket.Name,
	}

	pageLeft := true
	for pageLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		op, err := svc.ListBucketIntelligentTieringConfigurations(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3_bucket_intelligent_tiering_configuration.listBucketIntelligentTieringConfigurations", "api_error", err)
			return nil, err
		}

		for _, op := range op.IntelligentTieringConfigurationList {
			d.StreamListItem(ctx, &IntelligentTieringConfigurationInfo{bucket.Name, op})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if op.NextContinuationToken != nil {
			params.ContinuationToken = op.NextContinuationToken
		} else {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBucketIntelligentTieringConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucketName := d.EqualsQualString("bucket_name")
	id := d.EqualsQualString("id")

	if bucketName == "" || id == "" {
		return nil, nil
	}

	bucketRegion, err := doGetBucketRegion(ctx, d, h, bucketName)
	if err != nil {
		return nil, err
	}

	// Create client
	svc, err := S3Client(ctx, d, bucketRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket_intelligent_tiering_configuration.getBucketIntelligentTieringConfiguration", "client_error", err)
		return nil, err
	}

	params := &s3.GetBucketIntelligentTieringConfigurationInput{
		Bucket: &bucketName,
		Id:     &id,
	}

	op, err := svc.GetBucketIntelligentTieringConfiguration(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket_intelligent_tiering_configuration.getBucketIntelligentTieringConfiguration", "api_error", err)
		return nil, err
	}

	if op != nil && op.IntelligentTieringConfiguration != nil {
		return &IntelligentTieringConfigurationInfo{&bucketName, *op.IntelligentTieringConfiguration}, nil
	}

	return nil, nil
}
