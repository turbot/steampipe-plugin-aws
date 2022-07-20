package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func resolveBucketRegion(ctx context.Context, d *plugin.QueryData, bucketname *string) (loc *s3.GetBucketLocationOutput, _ error) {
	if bucketname == nil {
		return nil, nil
	}
	// have we already created and cached the service?
	bucketLocationCacheKey := fmt.Sprintf("s3-bucket-location-%s", *bucketname)
	if cachedData, ok := d.ConnectionManager.Cache.Get(bucketLocationCacheKey); ok {
		return cachedData.(*s3.GetBucketLocationOutput), nil
	}
	defer func() {
		if loc != nil {
			d.ConnectionManager.Cache.SetWithTTL(bucketLocationCacheKey, loc, 1*time.Minute)
		}
	}()

	defaultRegion := GetDefaultAwsRegion(d)

	// Create Session
	svc, err := S3Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketLocationInput{
		Bucket: bucketname,
	}

	// Specifies the Region where the bucket resides. For a list of all the Amazon
	// S3 supported location constraints by Region, see Regions and Endpoints (https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region).
	location, err := svc.GetBucketLocation(params)
	if err != nil {
		return nil, err
	}

	if location != nil && location.LocationConstraint != nil {
		// Buckets in eu-west-1 created through the AWS CLI or other API driven methods can return a location of "EU",
		// so we need to convert back
		if *location.LocationConstraint == "EU" {
			return &s3.GetBucketLocationOutput{
				LocationConstraint: aws.String("eu-west-1"),
			}, nil
		}
		return location, nil
	}

	// Buckets in us-east-1 have a LocationConstraint of null
	return &s3.GetBucketLocationOutput{
		LocationConstraint: aws.String("us-east-1"),
	}, nil
}

type s3ObjectRow struct {
	StorageClass        *string
	Key                 *string
	ETag                *string
	Size                *int64
	LastModified        *time.Time
	BucketName          *string
	BucketRegion        *string
	bucketHasLockConfig bool
	Prefix              string
}

func (o *s3ObjectRow) isOutpostObject() bool {
	// S3 on Outposts provides a new storage class, OUTPOSTS
	// as in https://docs.aws.amazon.com/AmazonS3/latest/userguide/S3onOutposts.html
	return strings.EqualFold(*o.StorageClass, "OUTPOSTS")
}

type s3ObjectAttributes struct {
	s3.GetObjectAttributesOutput
	parentRow *s3ObjectRow
}

type s3ObjectContent struct {
	s3.GetObjectOutput
}

func (obj *s3ObjectContent) ReadBody() (interface{}, error) {
	body, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	if utf8.Valid(body) {
		return string(body), nil
	}

	return base64.StdEncoding.EncodeToString(body), nil
}
