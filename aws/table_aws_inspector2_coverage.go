package aws

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"

	inspector2v1 "github.com/aws/aws-sdk-go/service/inspector2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsInspector2Coverage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector2_coverage",
		Description: "AWS Inspector2 Coverage",
		List: &plugin.ListConfig{
			Hydrate: listInspector2Coverage,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "resource_account_id", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				// {Name: "ec2_tags", Require: plugin.Optional},
				{Name: "ecr_image_tags", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "ecr_repository_name", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "lambda_function_name", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "lambda_function_runtime", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				// {Name: "lambda_function_tags", Require: plugin.Optional},
				{Name: "resource_id", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "resource_type", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "scan_status_code", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "scan_status_reason", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "scan_type", Operators: []string{"=", "<>"}, Require: plugin.Optional},
			},
		},

		GetMatrixItemFunc: SupportedRegionMatrix(inspector2v1.EndpointsID),

		// We *do not* use the common columns, because the account_id/region of
		// the default columns come from the call, *not* from the retutned data.
		// For inspector2, the account_id or region can vary within a single
		// call.
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				// The account id from the data, rather than from the call (getCommonColumns).
				Name:        "resource_account_id",
				Type:        proto.ColumnType_STRING,
				Description: "The AWS Account ID in which the resource is located.",
				Transform:   transform.FromField("AccountId"),
			},
			{
				Name:        "resource_id",
				Description: "The ID of the covered resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of the covered resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scan_type",
				Description: "The Amazon Inspector scan type covering the resource.",
				Type:        proto.ColumnType_STRING,
			},
			// the "resource_" (or "resource_metadata_") prefix seems overly
			// verbose, so we omit it; really, the "ResourceMetadata" only
			// exists to collect the union of type-specific values.
			{
				Name:        "ec2_ami_id",
				Description: "The ID of the Amazon Machine Image (AMI) used to launch the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceMetadata.Ec2.AmiId"),
			},
			{
				Name:        "ec2_platform",
				Description: "The platform of the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceMetadata.Ec2.Platform"),
			},
			{
				Name:        "ec2_tags_src",
				Description: "The tags attached to the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceMetadata.Ec2.Tags"),
			},
			{
				Name:        "ec2_tags",
				Description: "The tags attached to the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceMetadata.Ec2.Tags").Transform(jsonTags),
			},
			{
				Name:        "ecr_image_tags_src",
				Description: "Tags associated with the Amazon ECR image metadata.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceMetadata.EcrImage.Tags"),
			},
			{
				Name:        "ecr_image_tags",
				Description: "Tags associated with the Amazon ECR image metadata.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceMetadata.EcrImage.Tags"),
			},
			{
				Name:        "ecr_repository_name",
				Description: "The name of the Amazon ECR repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceMetadata.EcrRepository.Name"),
			},
			{
				Name:        "ecr_repository_scan_frequency",
				Description: "The frequency of scans for an object that contains details about the repository an Amazon ECR image resides in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceMetadata.EcrRepository.ScanFrequency"),
			},
			{
				Name:        "lambda_function_name",
				Description: "The name of a function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceMetadata.LambdaFunction.FunctionName"),
			},
			{
				Name:        "lambda_function_tags",
				Description: "The resource tags on an AWS Lambda function.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceMetadata.LambdaFunction.FunctionTags"),
			},
			{
				Name:        "lambda_function_layers",
				Description: "The layers for an AWS Lambda function. A Lambda function can have up to five layers.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceMetadata.LambdaFunction.Layers"),
			},
			{
				Name:        "lambda_function_runtime",
				Description: "An AWS Lambda function's runtime.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceMetadata.LambdaFunction.Runtime"),
			},
			{
				Name:        "scan_status_reason",
				Description: "The reason for the scan.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ScanStatus.Reason"),
			},
			{
				Name:        "scan_status_code",
				Description: "The status code of the scan.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ScanStatus.StatusCode"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceId").Transform(arnToTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceId").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

// column-to-filter mapping
type coverageStringFilterInfo struct {
	columnName  string
	filterField func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter
}

var coverageStringFilterList = []coverageStringFilterInfo{
	{
		columnName: "resource_account_id",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.AccountId)
		},
	},
	{
		columnName: "ecr_image_tags",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.EcrImageTags)
		},
	},
	{
		columnName: "ecr_repository_name",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.EcrRepositoryName)
		},
	},
	{
		columnName: "lambda_function_name",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.LambdaFunctionName)
		},
	},
	{
		columnName: "lambda_function_runtime",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.LambdaFunctionRuntime)
		},
	},
	{
		columnName: "resource_id",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.ResourceId)
		},
	},
	{
		columnName: "resource_type",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.ResourceType)
		},
	},
	{
		columnName: "scan_status_code",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.ScanStatusCode)
		},
	},
	{
		columnName: "scan_status_reason",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.ScanStatusReason)
		},
	},
	{
		columnName: "scan_type",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.ScanType)
		},
	},
}

// NOTE: I've commented out the map-filter keys, because they don't seem to work
// quite like I'd expect, and I can't find documentation that makes clear how
// they *should* work.  I'm reasonably confident that this code is a 1:1 mapping
// to the underlying map-filter feature, but don't want to expose it until I
// really understand how it's supposed to work.

// type coverageMapFilterInfo struct {
// 	columnName  string
// 	filterField func(f *types.CoverageFilterCriteria) *[]types.CoverageMapFilter
// }

// var coverageMapFilterList = []coverageMapFilterInfo{
// 	{
// 		columnName:  "ec2_tags",
// 		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageMapFilter {
//			return &(f.Ec2InstanceTags)
//		},
// 	},
// 	{
// 		columnName:  "lambda_function_tags",
// 		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageMapFilter {
//			return &(f.LambdaFunctionTags)
//		},
// 	},
// }

func listInspector2Coverage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := Inspector2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector2_coverage.listInspector2Coverage", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, etc., return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(200)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &inspector2.ListCoverageInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filter := &types.CoverageFilterCriteria{}

	// there are a bunch of keys with identically-structured syntax/types... we
	// can't easily include the filter field without reflection, but we *can* at
	// least construct the filter itself.
	for _, info := range coverageStringFilterList {
		if d.Quals[info.columnName] != nil {
			for _, q := range d.Quals[info.columnName].Quals {
				var comp types.CoverageStringComparison
				switch q.Operator {
				case "=":
					comp = types.CoverageStringComparisonEquals
				case "<>":
					comp = types.CoverageStringComparisonNotEquals
				}
				field := info.filterField(filter)
				*field = append(*field, types.CoverageStringFilter{
					Comparison: comp,
					Value:      aws.String(q.Value.GetStringValue()),
				})
			}
		}
	}

	// NOTE: I've commented out the map-filter keys, because they don't seem to
	// work quite like I'd expect, and I can't find documentation that makes
	// clear how they *should* work.  I'm reasonably confident that this code
	// is a 1:1 mapping to the underlying map-filter feature, but don't want to
	// expose it until I really understand how it's supposed to work.

	// // The CoverageMapFilter fields are JSON, which means we need to expand some
	// // number of key:value pairs into filter values.
	// for _, info := range coverageMapFilterList {
	// 	if d.Quals[info.columnName] != nil {
	// 		for _, q := range d.Quals[info.columnName].Quals {
	// 			var comp types.CoverageMapComparison
	// 			switch q.Operator {
	// 			case "=":
	// 				comp = types.CoverageMapComparisonEquals
	// 			}
	// 			parts := strings.SplitN(q.Value.GetStringValue(), ":", 2)
	// 			if len(parts) != 2 {
	// 				plugin.Logger(ctx).Error(fmt.Sprintf("filter value for %q should be in KEY:VALUE format, got %q", info.columnName, q.Value.GetStringValue()))
	// 				// assume an empty value, which is about all we can do...
	// 				parts = append(parts, "")
	// 			}
	// 			field := info.filterField(filter)
	// 			// Even though this appears to be the correct way to build a map
	// 			// filter, we get no results... does "=" imply *exact* tag
	// 			// matching? (i.e. "key:val" means *only* a tag of "key:val",
	// 			// rather than "has a tag 'key' whose value is 'val'... and may
	// 			// have other tags as well"?)
	// 			*field = append(*field, types.CoverageMapFilter{
	// 				Comparison: comp,
	// 				Key:        aws.String(parts[0]),
	// 				Value:      aws.String(parts[1]),
	// 			})
	// 		}
	// 	}
	// }

	input.FilterCriteria = filter

	paginator := inspector2.NewListCoveragePaginator(svc, input, func(o *inspector2.ListCoveragePaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_inspector2_coverage.listInspector2Coverage", "api_error", err)
			return nil, err
		}

		for _, item := range output.CoveredResources {
			// item := item
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// TRANSFORM HELPERS

func jsonTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	// always return tags in sorted order, for consistency (even though it
	// appears AWS does as well, this ensures it)
	m, ok := d.Value.(map[string]string)
	if !ok {
		return nil, errors.New("did not get expected map[string]string tags")
	}

	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	keyVals := []string{}
	for _, k := range keys {
		keyVals = append(keyVals, fmt.Sprintf("%q:%q", k, m[k]))
	}

	// the only real difference in presentation is that there's no enclosing
	// "{}".
	return strings.Join(keyVals, ","), nil
}
