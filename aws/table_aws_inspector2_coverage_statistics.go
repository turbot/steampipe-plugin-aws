package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

// A note about the coverage statistics: Because steampipe is modelling cloud
// resources as SQL tables, things that are less like resources do not map
// particularly well to a table.  The inspector2 coverage statistics is one such
// thing. Other inspector2 values have a 1:1 relationship with a resource, and
// thus the mapping makes sense, but the coverage statistics are a summary, and
// not tied to any one resource.  Further, the filtering criteria (which is the
// same as the coverage table) and grouping keys are *not* part of the result
// data, and thus cannot be defined as columns... and thus steampipe doesn't
// really have a way to represent/map the underlying API call.  For now, this
// table is present for completeness, but there is no way to get counts_by_group
// to be anything other than an empty value; only the total_counts column has
// value.
//
// Note that *if* it is decided to add dummy columns so as to enable filtering
// (and thus counts_by_group), the filter columns are exactly the same ones as
// in inspector2_coverage, and could be refactored/shared.
func tableAwsInspector2CoverageStatistics(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_inspector2_coverage_statistics",
		Description: "AWS Inspector2 Coverage Statistics",
		List: &plugin.ListConfig{
			Hydrate: listInspector2CoverageStatistics,
			Tags:    map[string]string{"service": "inspector2", "action": "ListCoverageStatistics"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_INSPECTOR2_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "total_counts",
				Description: "The total number for all groups",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "counts_by_group",
				Description: "An array with the number for each group",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

//// LIST FUNCTION

func listInspector2CoverageStatistics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := Inspector2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_inspector2_coveragestatistics.listInspector2CoverageStatistics", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &inspector2.ListCoverageStatisticsInput{
		// no FilterCriteria, no GroupBy
	}

	paginator := inspector2.NewListCoverageStatisticsPaginator(svc, input, func(o *inspector2.ListCoverageStatisticsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_inspector2_coverage_statistics.listInspector2CoverageStatistics", "api_error", err)
			return nil, err
		}

		d.StreamListItem(ctx, output)
	}

	return nil, err
}
