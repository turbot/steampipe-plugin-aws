package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSSMAssociation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_association",
		Description: "AWS SSM Association",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("association_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"AssociationDoesNotExist", "ValidationException"}),
			},
			Hydrate: getAwsSSMAssociation,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMAssociations,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "association_name", Require: plugin.Optional},
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "last_execution_date", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "association_id",
				Description: "The ID created by the system when you create an association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_name",
				Description: "The Name of association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the association.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSSMAssociationARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "document_name",
				Description: "The name of the Systems Manager document.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "date",
				Description: "The date when the association was made.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "compliance_severity",
				Description: "A cron expression that specifies a schedule when the association runs.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "apply_only_at_cron_interval",
				Description: "By default, when you create a new associations, the system runs it immediately after it is created and then according to the schedule you specified. Specify this option if you don't want an association to run immediately after you create it. This parameter is not supported for rate expressions.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "association_version",
				Description: "The association version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "automation_target_parameter_name",
				Description: "Specify the target for the association. This target is required for associations that use an Automation document and target resources by using rate controls.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "document_version",
				Description: "The version of the document used in the association.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_id",
				Description: "The ID of the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMAssociation,
			},
			{
				Name:        "last_execution_date",
				Description: "The date on which the association was last run.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_successful_execution_date",
				Description: "The last date on which the association was successfully run.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSSMAssociation,
			},
			{
				Name:        "last_update_association_date",
				Description: "The date when the association was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSSMAssociation,
			},
			{
				Name:        "schedule_expression",
				Description: "A cron expression that specifies a schedule when the association runs.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_concurrency",
				Description: "The maximum number of targets allowed to run the association at the same time.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_errors",
				Description: "The number of errors that are allowed before the system stops sending requests to run the association on additional targets.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sync_compliance",
				Description: "The mode for generating association compliance. You can specify AUTO or MANUAL.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "overview",
				Description: "Information about the association.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "output_location",
				Description: "An S3 bucket where you want to store the output details of the request.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMAssociation,
			},
			{
				Name:        "parameters",
				Description: "A description of the parameters for a document.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMAssociation,
			},
			{
				Name:        "status",
				Description: "The status of the association. Status can be: Pending, Success, or Failed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMAssociation,
				Transform:   transform.FromField("Overview.Status"),
			},
			{
				Name:        "targets",
				Description: "A cron expression that specifies a schedule when the association runs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "target_locations",
				Description: "The combination of AWS Regions and AWS accounts where you want to run the association.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMAssociation,
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSSMAssociationARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSSMAssociations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_association.listAwsSSMAssociations", "connection_error", err)
		return nil, err
	}

	maxItems := int32(50)
	input := &ssm.ListAssociationsInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	filters := buildSSMAssociationFilter(d.Quals)
	quals := d.Quals
	if quals["last_execution_date"] != nil {
		f := types.AssociationFilter{}
		for _, q := range quals["last_execution_date"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">=", ">":
				f.Key = types.AssociationFilterKeyLastExecutedAfter
				f.Value = aws.String(fmt.Sprint(timestamp))
			case "<", "<=":
				f.Key = types.AssociationFilterKeyLastExecutedBefore
				f.Value = aws.String(fmt.Sprint(timestamp))
			}
		}
		filters = append(filters, f)
	}

	if len(filters) > 0 {
		input.AssociationFilterList = filters
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := ssm.NewListAssociationsPaginator(svc, input, func(o *ssm.ListAssociationsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_association.listAwsSSMAssociations", "api_error", err)
			return nil, err
		}

		for _, association := range output.Associations {
			d.StreamListItem(ctx, association)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsSSMAssociation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var id string
	if h.Item != nil {
		id = associationID(h.Item)
	} else {
		id = d.KeyColumnQuals["association_id"].GetStringValue()
	}

	// Empty input id check
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_association.getAwsSSMAssociation", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &ssm.DescribeAssociationInput{
		AssociationId: aws.String(id),
	}

	// Get call
	data, err := svc.DescribeAssociation(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_association.getAwsSSMAssociation", "api_error", err)
		return nil, err
	}
	return *data.AssociationDescription, nil
}

func getSSMAssociationARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	associationData := associationID(h.Item)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_association.getSSMAssociationARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":association/" + associationData

	return arn, nil
}

func associationID(item interface{}) string {
	switch item := item.(type) {
	case types.Association:
		return *item.AssociationId
	case types.AssociationDescription:
		return *item.AssociationId
	}
	return ""
}

//// UTILITY FUNCTION

// Build ssm association list call input filter
func buildSSMAssociationFilter(quals plugin.KeyColumnQualMap) []types.AssociationFilter {
	filters := make([]types.AssociationFilter, 0)

	filterQuals := map[string]string{
		"association_name": string(types.AssociationFilterKeyAssociationName),
		"instance_id":      string(types.AssociationFilterKeyInstanceId),
		"status":           string(types.AssociationFilterKeyStatus),
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.AssociationFilter{
				Key: types.AssociationFilterKey(filterName),
			}

			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Value = aws.String(val)
				filters = append(filters, filter)
			}
		}
	}
	return filters
}
