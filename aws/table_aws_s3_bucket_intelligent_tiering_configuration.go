package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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

	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	bucket := h.Item.(types.Bucket)
	location, err := getIntelligentTieringBucketLocation(ctx, d, h)
	if err != nil {
		return nil, nil
	} else if location == nil {
		return nil, nil
	}

	if d.EqualsQualString("bucket_name") != "" && d.EqualsQualString("bucket_name") != *bucket.Name {
		return nil, nil
	}

	// Create client
	svc, err := S3Client(ctx, d, fmt.Sprint(location))
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

	// Bucket location will be nil if getBucketLocation returned an error but
	// was ignored through ignore_error_codes config arg
	location, err := getIntelligentTieringBucketLocation(ctx, d, h)
	if err != nil {
		return nil, nil
	} else if location == nil {
		return nil, nil
	}

	// Create client
	svc, err := S3Client(ctx, d, fmt.Sprint(location))
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

func getIntelligentTieringBucketLocation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var bucketName string
	if h.Item != nil {
		bucketName = *h.Item.(types.Bucket).Name
	} else {
		bucketName = d.EqualsQuals["bucket_name"].GetStringValue()
	}

	if bucketName == "" {
		return nil, nil
	}

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket_intelligent_tiering_configuration.getIntelligentTieringBucketLocation", "get_common_columns_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// have we already created and cached the session?
	cacheKey := "getIntelligentTieringBucketLocation" + bucketName + commonColumnData.AccountId

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(string), nil
	}

	// Unlike most services, S3 buckets are a global list. They can be retrieved
	// from any single region. It's best to use the client region of the user
	// (e.g. closest to them).
	clientRegion, err := getDefaultRegion(ctx, d, h)
	if err != nil {
		return "", err
	}

	svc, err := S3Client(ctx, d, clientRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket_intelligent_tiering_configuration.getIntelligentTieringBucketLocation", "get_client_error", err, "clientRegion", clientRegion)
		return "", err
	}

	params := &s3.GetBucketLocationInput{Bucket: aws.String(bucketName), ExpectedBucketOwner: aws.String(commonColumnData.AccountId)}

	// Specifies the Region where the bucket resides. For a list of all the Amazon
	// S3 supported location constraints by Region, see Regions and Endpoints (https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region).
	location, err := svc.GetBucketLocation(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_bucket_intelligent_tiering_configuration.getIntelligentTieringBucketLocation", "bucket_name", bucketName, "clientRegion", clientRegion, "api_error", err)
		return "", err
	}
	var locationConstraint string
	if location != nil && location.LocationConstraint != "" {
		// Buckets in eu-west-1 created through the AWS CLI or other API driven methods can return a location of "EU",
		// so we need to convert back
		if location.LocationConstraint == "EU" {
			locationConstraint = "eu-west-1"
			d.ConnectionManager.Cache.Set(cacheKey, locationConstraint)
			return locationConstraint, nil
		}
		d.ConnectionManager.Cache.Set(cacheKey, string(location.LocationConstraint))
		return string(location.LocationConstraint), nil
	}

	// Buckets in us-east-1 have a LocationConstraint of null
	locationConstraint = "us-east-1"
	d.ConnectionManager.Cache.Set(cacheKey, locationConstraint)
	return locationConstraint, nil
}
