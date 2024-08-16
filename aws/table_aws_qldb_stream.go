package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/qldb"

	qldbv1 "github.com/aws/aws-sdk-go/service/qldb"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsQLDBStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_qldb_stream",
		Description: "AWS QLDB Stream",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"ledger_name", "stream_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsQldbStream,
			Tags:    map[string]string{"service": "qldb", "action": "DescribeStream"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsQldbLedgers,
			Hydrate:       listAwsQldbStreams,
			KeyColumns: plugin.OptionalColumns([]string{"ledger_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Tags:          map[string]string{"service": "qldb", "action": "ListStreams"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsQldbStream,
				Tags: map[string]string{"service": "qldb", "action": "DescribeStream"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(qldbv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "ledger_name",
				Description: "The name of the ledger.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stream_id",
				Description: "The UUID (represented in Base62-encoded text) of the QLDB journal stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stream_name",
				Description: "The user-defined name of the QLDB journal stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the QLDB journal stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time, in epoch time format, when the QLDB journal stream was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "The current state of the QLDB journal stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role that grants QLDB permissions for a journal stream to write data records to a Kinesis Data Streams resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "error_cause",
				Description: "The error message that describes the reason that a stream has a status of IMPAIRED or FAILED .",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "exclusive_end_time",
				Description: "The exclusive date and time that specifies when the stream ends.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "inclusive_start_time",
				Description: "The inclusive start date and time from which to start streaming journal data.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "kinesis_configuration",
				Description: "The configuration settings of the Amazon Kinesis Data Streams destination for a QLDB journal stream.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsQldbStreams(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get client
	if h.Item == nil {
		return nil, nil
	}
	ledger := h.Item.(qldb.DescribeLedgerOutput)

	ledgerName := d.EqualsQualString("ledger_name")
	if ledgerName != ""  && ledgerName != *ledger.Name {
		return nil, nil
	}

	svc, err := QLDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_qldb_stream.listAwsQldbStreams", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &qldb.ListJournalKinesisStreamsForLedgerInput{
		MaxResults: &maxLimit,
		LedgerName: ledger.Name,
	}

	paginator := qldb.NewListJournalKinesisStreamsForLedgerPaginator(svc, input, func(o *qldb.ListJournalKinesisStreamsForLedgerPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_qldb_stream.listAwsQldbStreams", "api_error", err)
			return nil, err
		}

		for _, item := range output.Streams {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsQldbStream(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := QLDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_qldb_stream.getAwsQldbStream", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	ledgerName := d.EqualsQualString("ledger_name")
	streamId := d.EqualsQualString("stream_id")

	// Empty Check
	if ledgerName == "" || streamId == "" {
		return nil, nil
	}

	// Build the params
	params := &qldb.DescribeJournalKinesisStreamInput{
		LedgerName: &ledgerName,
		StreamId:   &streamId,
	}

	// Get call
	data, err := svc.DescribeJournalKinesisStream(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_qldb_stream.getAwsQldbStream", "api_error", err)
		return nil, err
	}

	if data.Stream != nil {
		return data.Stream, nil
	}

	return nil, nil
}
