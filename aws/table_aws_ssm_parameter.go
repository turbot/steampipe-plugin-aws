package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsSSMParameter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_parameter",
		Description: "AWS SSM Parameter",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException"}),
			},
			Hydrate: getAwsSSMParameter,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMParameters,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "type", Require: plugin.Optional},
				{Name: "key_id", Require: plugin.Optional},
				{Name: "tier", Require: plugin.Optional},
				{Name: "data_type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The parameter name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of parameter. Valid parameter types include the following: String, StringList, and SecureString.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "The value of parameter.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMParameterDetails,
				Transform:   transform.FromField("Parameter.Value"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the parameter.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMParameterDetails,
				Transform:   transform.FromField("Parameter.ARN"),
			},
			{
				Name:        "data_type",
				Description: "The data type of the parameter, such as text or aws:ec2:image. The default is text.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_id",
				Description: "The ID of the query key used for this parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_date",
				Description: "Date the parameter was last changed or updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_modified_user",
				Description: "Amazon Resource Name (ARN) of the AWS user who last changed the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The parameter version.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "selector",
				Description: "Either the version number or the label used to retrieve the parameter value.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMParameterDetails,
				Transform:   transform.FromField("Parameter.Selector"),
			},
			{
				Name:        "source_result",
				Description: "SourceResult is the raw result or response from the source. Applies to parameters that reference information in other AWS services.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMParameterDetails,
				Transform:   transform.FromField("Parameter.SourceResult"),
			},
			{
				Name:        "tier",
				Description: "The parameter tier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "policies",
				Description: "A list of policies associated with a parameter. Parameter policies help you manage a growing set of parameters by enabling you to assign specific criteria to a parameter such as an expiration date or time to live. Parameter policies are especially helpful in forcing you to update or delete passwords and configuration data stored in Parameter Store.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policies"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the parameter.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMParameterTags,
				Transform:   transform.FromField("TagList"),
			},

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMParameterTags,
				Transform:   transform.FromField("TagList").Transform(ssmTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMParameterAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSSMParameters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsSSMParameters")

	// Create session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &ssm.DescribeParametersInput{
		MaxResults: aws.Int64(50),
	}

	filters := buildSsmParameterFilter(d.Quals, ctx)
	if len(filters) > 0 {
		input.ParameterFilters = filters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribeParametersPages(
		input,
		func(page *ssm.DescribeParametersOutput, isLast bool) bool {
			for _, parameter := range page.Parameters {
				d.StreamListItem(ctx, parameter)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSSMParameter(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSSMParameter")

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.DescribeParametersInput{
		ParameterFilters: []*ssm.ParameterStringFilter{
			{
				Key: types.String("Name"),
				Values: []*string{
					types.String(name),
				},
			},
		},
	}

	// Get call
	data, err := svc.DescribeParameters(params)
	if err != nil {
		logger.Debug("getAwsSSMParameter", "ERROR", err)
		return nil, err
	}

	if len(data.Parameters) > 0 {
		return data.Parameters[0], nil
	}

	return nil, nil
}

func getAwsSSMParameterDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSSMParameter")

	parameterData := h.Item.(*ssm.ParameterMetadata)

	// Create Session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.GetParameterInput{
		Name:           parameterData.Name,
		WithDecryption: types.Bool(true),
	}

	// Get call
	op, err := svc.GetParameter(params)
	if err != nil {
		logger.Debug("getAwsSSMParameterDetails", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getAwsSSMParameterTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSSMParameterTags")

	parameterData := h.Item.(*ssm.ParameterMetadata)

	// Create Session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.ListTagsForResourceInput{
		ResourceType: types.String("Parameter"),
		ResourceId:   parameterData.Name,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("getAwsSSMParameterTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getAwsSSMParameterAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsSSMParameterAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	parameterData := h.Item.(*ssm.ParameterMetadata)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":parameter"

	if strings.HasPrefix(*parameterData.Name, "/") {
		aka = aka + *parameterData.Name
	} else {
		aka = aka + "/" + *parameterData.Name
	}

	return []string{aka}, nil
}

//// TRANSFORM FUNCTIONS

func ssmTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ssmTagListToTurbotTags")
	tagList := d.Value.([]*ssm.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

//// UTILITY FUNCTION

// Build ssm parameter list call input filter
func buildSsmParameterFilter(quals plugin.KeyColumnQualMap, ctx context.Context) []*ssm.ParameterStringFilter {
	filters := make([]*ssm.ParameterStringFilter, 0)

	filterQuals := map[string]string{
		"type":      "Type",
		"key_id":    "KeyId",
		"tier":      "Tier",
		"data_type": "DataType",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := ssm.ParameterStringFilter{
				Key:    aws.String(filterName),
				Option: aws.String("Equals"),
			}

			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []*string{&val}
			} else {
				filter.Values = value.([]*string)
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
