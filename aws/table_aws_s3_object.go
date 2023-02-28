package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ettle/strcase"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func tableS3Object(ctx context.Context) *plugin.Table {
	bucketName := ctx.Value(contextKey("BucketName")).(string)
	return &plugin.Table{
		Name:        bucketName,
		Description: "Objects for " + bucketName + ".",
		List: &plugin.ListConfig{
			Hydrate: listS3Objects(ctx, bucketName),
		},
		Columns: awsAccountColumns(getS3ObjectsDynamicColumns(ctx, bucketName)),
	}
}

func getS3ObjectsDynamicColumns(_ context.Context, bucketName string) []*plugin.Column {
	columns := []*plugin.Column{}

	// default metadata columns
	allColumns := []string{"name", "uid", "kind", "api_version", "namespace", "creation_timestamp", "labels"}

	// add the spec columns
	flag := 0
	schemaSpec := versionSchemaSpec.(v1.JSONSchemaProps)
	for k, v := range schemaSpec.Properties {
		for _, specColumn := range allColumns {
			if specColumn == strcase.ToSnake(k) {
				flag = 1
				column := &plugin.Column{
					Name:        "spec_" + strcase.ToSnake(k),
					Description: v.Description,
					Transform:   transform.FromP(extractSpecProperty, k),
				}
				setDynamicColumns(v, column)
				columns = append(columns, column)
			}
		}
		if flag == 0 {
			column := &plugin.Column{
				Name:        strcase.ToSnake(k),
				Description: v.Description,
				Transform:   transform.FromP(extractSpecProperty, k),
			}
			allColumns = append(allColumns, strcase.ToSnake(k))
			setDynamicColumns(v, column)
			columns = append(columns, column)
		}
	}

	// add the status columns
	flag = 0
	schemaStatus := versionSchemaStatus.(v1.JSONSchemaProps)
	for k, v := range schemaStatus.Properties {
		for _, statusColumn := range allColumns {
			if statusColumn == strcase.ToSnake(k) {
				flag = 1
				column := &plugin.Column{
					Name:        "status_" + strcase.ToSnake(k),
					Description: v.Description,
					Transform:   transform.FromP(extractStatusProperty, k),
				}
				setDynamicColumns(v, column)
				columns = append(columns, column)
			}
		}
		if flag == 0 {
			column := &plugin.Column{
				Name:        strcase.ToSnake(k),
				Description: v.Description,
				Transform:   transform.FromP(extractStatusProperty, k),
			}
			setDynamicColumns(v, column)
			columns = append(columns, column)
		}
	}

	return columns
}

type CRDResourceInfo struct {
	Name              interface{}
	UID               interface{}
	CreationTimestamp interface{}
	Kind              interface{}
	APIVersion        interface{}
	Namespace         interface{}
	Annotations       interface{}
	Spec              interface{}
	Labels            interface{}
	Status            interface{}
}

// //// HYDRATE FUNCTIONS

func listS3Objects(ctx context.Context, crdName string, bucketName string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

		// Unlike most services, S3 buckets are a global list. They can be retrieved
		// from any single region.  We must list buckets from the default region to
		// get the actual creation_time of the bucket, in all other regions the list
		// returns the time when the bucket was last modified. See
		// https://www.marksayson.com/blog/s3-bucket-creation-dates-s3-master-regions/
		defaultRegion, err := getLastResortRegion(ctx, d, h)
		if err != nil {
			return nil, err
		}

		svc, err := S3Client(ctx, d, defaultRegion)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3_bucket.listS3Buckets", "get_client_error", err, "defaultRegion", defaultRegion)
			return nil, err
		}

		// execute list call
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		}
		objects, err := svc.ListObjects(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3_bucket.listS3Buckets", "api_error", err, "defaultRegion", defaultRegion)
			return nil, err
		}

		for _, bucket := range objects. {
			d.StreamListItem(ctx, bucket)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		return nil, nil
	}
}

func extractSpecProperty(_ context.Context, d *transform.TransformData) (interface{}, error) {
	ob := d.HydrateItem.(*CRDResourceInfo).Spec
	if ob == nil {
		return nil, nil
	}
	param := d.Param.(string)
	spec := ob.(map[string]interface{})
	if spec[param] != nil {
		return spec[param], nil
	}

	return nil, nil
}

func extractStatusProperty(_ context.Context, d *transform.TransformData) (interface{}, error) {
	ob := d.HydrateItem.(*CRDResourceInfo).Status
	if ob == nil {
		return nil, nil
	}
	param := d.Param.(string)
	status := ob.(map[string]interface{})
	if status[param] != nil {
		return status[param], nil
	}

	return nil, nil
}

func setDynamicColumns(v v1.JSONSchemaProps, column *plugin.Column) {
	switch v.Type {
	case "string":
		column.Type = proto.ColumnType_STRING
	case "integer":
		column.Type = proto.ColumnType_INT
	case "boolean":
		column.Type = proto.ColumnType_BOOL
	case "date", "dateTime":
		column.Type = proto.ColumnType_TIMESTAMP
	case "double":
		column.Type = proto.ColumnType_DOUBLE
	default:
		column.Type = proto.ColumnType_JSON
	}
}
