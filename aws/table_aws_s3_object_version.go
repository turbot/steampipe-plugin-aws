package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsS3ObjectVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_object_version",
		Description: "List AWS S3 Object versions in S3 buckets by bucket name.",
		List: &plugin.ListConfig{
			Hydrate: listS3ObjectVersions,
			Tags:    map[string]string{"service": "s3", "action": "ListObjectVersions"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket_name", Require: plugin.Required, CacheMatch: "exact"},
				{Name: "prefix", Require: plugin.Optional},
				{Name: "encoding_type", Require: plugin.Optional},
				{Name: "delimiter", Require: plugin.Optional},
				{Name: "version_id_marker", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getBucketLocationForObjects,
				Tags: map[string]string{"service": "s3", "action": "GetBucketLocation"},
			},
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "bucket_name",
				Description: "The name of the container bucket of this object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("bucket_name"),
			},
			{
				Name:        "delimiter",
				Description: "The delimiter grouping the included keys.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encoding_type",
				Description: "Encoding type used by Amazon S3 to encode object key names in the XML response.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_id_marker",
				Description: "Marks the last version of the key returned in a truncated response.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "prefix",
				Description: "Selects objects that start with the value supplied by this parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_truncated",
				Description: "A flag that indicates whether Amazon S3 returned all of the results that satisfied the search criteria.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "common_prefixes",
				Description: "All of the keys rolled up into a common prefix count as a single return when calculating the number of returns.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "delete_markers",
				Description: "Specifies caching behavior along the request/reply chain.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version",
				Description: "Container for version information.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Version.Key"),
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the object is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBucketLocationForObjects,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type ObjectVersionDetails struct {
	CommonPrefixes  []types.CommonPrefix
	DeleteMarkers   []types.DeleteMarkerEntry
	Delimiter       *string
	EncodingType    types.EncodingType
	IsTruncated     bool
	BucketName      *string
	Prefix          *string
	VersionIdMarker *string
	Version         types.ObjectVersion
}

func listS3ObjectVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Bucket location will be nil if getBucketLocationForObjects returned an error but
	// was ignored through ignore_error_codes config arg
	location, err := getBucketLocationForObjects(ctx, d, h)
	if err != nil {
		return nil, err
	} else if location == "" {
		return nil, nil
	}

	svc, err := S3Client(ctx, d, fmt.Sprint(location))
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_object_version.listS3ObjectVersions", "get_client_error", err)
		return nil, err
	}

	bucketName := d.EqualsQuals["bucket_name"].GetStringValue()

	// default supported max value is 1000 by ListObjectVersions
	maxItems := int32(1000)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	input := &s3.ListObjectVersionsInput{
		Bucket:  aws.String(bucketName),
		MaxKeys: maxItems,
	}

	if d.EqualsQualString("prefix") != "" {
		input.Prefix = aws.String(d.EqualsQualString("prefix"))
	}
	if d.EqualsQualString("encoding_type") != "" {
		input.EncodingType = types.EncodingType(d.EqualsQualString("encoding_type"))
	}
	if d.EqualsQualString("delimeter") != "" {
		input.Delimiter = aws.String(d.EqualsQualString("delimeter"))
	}
	if d.EqualsQualString("version_id_marker") != "" {
		input.VersionIdMarker = aws.String(d.EqualsQualString("version_id_marker"))
	}

	// execute list call
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		objects, err := svc.ListObjectVersions(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3_object_version.ListObjectsV2", "api_error", err)
			return nil, err
		}

		for _, version := range objects.Versions {
			d.StreamListItem(ctx, ObjectVersionDetails{
				CommonPrefixes:  objects.CommonPrefixes,
				DeleteMarkers:   objects.DeleteMarkers,
				Delimiter:       objects.Delimiter,
				BucketName:      objects.Name,
				IsTruncated:     objects.IsTruncated,
				EncodingType:    objects.EncodingType,
				Prefix:          objects.Prefix,
				VersionIdMarker: objects.VersionIdMarker,
				Version:         version,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		input.KeyMarker = objects.NextVersionIdMarker
		if objects.NextVersionIdMarker == nil {
			break
		}
	}

	return nil, err
}
