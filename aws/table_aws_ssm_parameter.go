package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	ssmv1 "github.com/aws/aws-sdk-go/service/ssm"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSSMParameter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_parameter",
		Description: "AWS SSM Parameter",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException"}),
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
		GetMatrixItemFunc: SupportedRegionMatrix(ssmv1.EndpointsID),
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

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_parameter.listAwsSSMParameters", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(50)
	input := &ssm.DescribeParametersInput{}

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

	filters := buildSSMParameterFilter(d.Quals, ctx)
	if len(filters) > 0 {
		input.ParameterFilters = filters
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := ssm.NewDescribeParametersPaginator(svc, input, func(o *ssm.DescribeParametersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_parameter.listAwsSSMParameters", "api_error", err)
			return nil, err
		}

		for _, parameter := range output.Parameters {
			d.StreamListItem(ctx, parameter)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsSSMParameter(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_parameter.getAwsSSMParameter", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &ssm.DescribeParametersInput{
		ParameterFilters: []types.ParameterStringFilter{
			{
				Key:    aws.String("Name"),
				Values: []string{name},
			},
		},
	}

	// Get call
	data, err := svc.DescribeParameters(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_parameter.getAwsSSMParameter", "api_error", err)
		return nil, err
	}

	if len(data.Parameters) > 0 {
		return data.Parameters[0], nil
	}

	return nil, nil
}

func getAwsSSMParameterDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	parameterData := h.Item.(types.ParameterMetadata)

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_parameter.getAwsSSMParameterDetails", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &ssm.GetParameterInput{
		Name:           parameterData.Name,
		WithDecryption: aws.Bool(true),
	}

	// Get call
	op, err := svc.GetParameter(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			// In case the KMS key encrypting the SSM Parameter value is disabled, below error is thrown
			// operation error SSM: GetParameter, https response error StatusCode: 400, RequestID: 0965014b-77ab-4847-98d4-2b9e09a68385, InvalidKeyId: arn:aws:kms:us-east-1:111122223333:key/1a2b3c4d-f6b4-4c5b-97e7-123456ab210c is disabled. (Service: AWSKMS; Status Code: 400; Error Code: DisabledException; Request ID: 7b6ae355-c99a-4cad-b2c3-4b40c0abdda9; Proxy: null)
			if ae.ErrorCode() == "InvalidKeyId" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_ssm_parameter.getAwsSSMParameterDetails", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsSSMParameterTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	parameterData := h.Item.(types.ParameterMetadata)

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_parameter.getAwsSSMParameterTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &ssm.ListTagsForResourceInput{
		ResourceType: types.ResourceTypeForTagging("Parameter"),
		ResourceId:   parameterData.Name,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_parameter.getAwsSSMParameterTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsSSMParameterAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	parameterData := h.Item.(types.ParameterMetadata)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_parameter.getAwsSSMParameterTags", "common_data_error", err)
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
	tagList := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if len(tagList) > 0 {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

//// UTILITY FUNCTION

// Build ssm parameter list call input filter
func buildSSMParameterFilter(quals plugin.KeyColumnQualMap, ctx context.Context) []types.ParameterStringFilter {
	filters := make([]types.ParameterStringFilter, 0)

	filterQuals := map[string]string{
		"type":      "Type",
		"key_id":    "KeyId",
		"tier":      "Tier",
		"data_type": "DataType",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.ParameterStringFilter{
				Key:    aws.String(filterName),
				Option: aws.String("Equals"),
			}

			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			} else {
				filter.Values = value.([]string)
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
