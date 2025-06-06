package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITION

func tableAwsInspector2Coverage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector2_coverage",
		Description: "AWS Inspector2 Coverage",
		List: &plugin.ListConfig{
			Hydrate: listInspector2Coverage,
			Tags:    map[string]string{"service": "inspector2", "action": "ListCoverage"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "source_account_id", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "ec2_instance_tags", Operators: []string{"="}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ecr_image_tag", Operators: []string{"=", "<>"}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "ecr_repository_name", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "lambda_function_name", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "lambda_function_runtime", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "lambda_function_tags", Operators: []string{"="}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "resource_id", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "resource_type", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "scan_status_code", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "scan_status_reason", Operators: []string{"=", "<>"}, Require: plugin.Optional},
				{Name: "scan_type", Operators: []string{"=", "<>"}, Require: plugin.Optional},
			},
		},

		GetMatrixItemFunc: SupportedRegionMatrix(AWS_INSPECTOR2_SERVICE_ID),

		// We *do not* use the common columns, because the account_id/region of
		// the default columns come from the call, *not* from the retutned data.
		// For inspector2, the account_id or region can vary within a single
		// call.
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				// The account id from the data, rather than from the call (getCommonColumns).
				Name:        "source_account_id",
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
			{
				Name:        "last_scanned_at",
				Description: "The date and time the resource was last checked for vulnerabilities.",
				Type:        proto.ColumnType_TIMESTAMP,
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
				Name:        "ecr_image_tag",
				Description: "Tags associated with the Amazon ECR image metadata.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ecr_image_tag"),
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
			{
				Name:        "ec2_platform",
				Description: "The platform of the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceMetadata.Ec2.Platform"),
			},
			{
				Name:        "ec2_instance_tags",
				Description: "The tags attached to the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceMetadata.Ec2.Tags"),
			},
			{
				Name:        "ecr_image_tags",
				Description: "Tags associated with the Amazon ECR image metadata.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceMetadata.EcrImage.Tags"),
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

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceId").Transform(arnToTitle),
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
		columnName: "source_account_id",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageStringFilter {
			return &(f.AccountId)
		},
	},
	{
		columnName: "ecr_image_tag",
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

type coverageMapFilterInfo struct {
	columnName  string
	filterField func(f *types.CoverageFilterCriteria) *[]types.CoverageMapFilter
}

var coverageMapFilterList = []coverageMapFilterInfo{
	{
		columnName: "ec2_instance_tags",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageMapFilter {
			return &(f.Ec2InstanceTags)
		},
	},
	{
		columnName: "lambda_function_tags",
		filterField: func(f *types.CoverageFilterCriteria) *[]types.CoverageMapFilter {
			return &(f.LambdaFunctionTags)
		},
	},
}

func buildCoverageTagFilter(d *plugin.QueryData, filter *types.CoverageFilterCriteria) {
	for _, info := range coverageMapFilterList {
		if d.Quals[info.columnName] != nil {
			field := info.filterField(filter)
			tagValue := make(map[string]string, 0)
			for _, q := range d.Quals[info.columnName].Quals {
				val := q.Value.GetJsonbValue()
				if val != "" && q.Operator == "=" {
					_ = json.Unmarshal([]byte(val), &tagValue)
					for k, v := range tagValue {
						tagfilter := types.CoverageMapFilter{
							Comparison: types.CoverageMapComparisonEquals,
							Key:        aws.String(k),
							Value:      aws.String(v),
						}
						*field = append(*field, tagfilter)
					}
				}
			}
		}
	}
}

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
			maxLimit = limit
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

	input.FilterCriteria = filter

	buildCoverageTagFilter(d, filter)

	paginator := inspector2.NewListCoveragePaginator(svc, input, func(o *inspector2.ListCoveragePaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

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
