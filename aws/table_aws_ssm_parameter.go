package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsSSMParameter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_parameter",
		Description: "AWS SSM Parameter",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("snapshot_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidSnapshot.NotFound", "InvalidSnapshotID.Malformed"}),
			Hydrate:           getAwsSSMParameter,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMParameters,
		},
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
				Name:        "tier",
				Description: "The parameter tier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "Policies",
				Description: "A list of policies associated with a parameter.",
				Type:        proto.ColumnType_JSON,
			},
			// {
			// 	Name:        "tags_src",
			// 	Description: "A list of tags assigned to the snapshot",
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromField("Tags"),
			// },

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			// {
			// 	Name:        "tags",
			// 	Description: resourceInterfaceDescription("tags"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.From(ec2SnapshotTurbotTags),
			// },
			// {
			// 	Name:        "akas",
			// 	Description: resourceInterfaceDescription("akas"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     getAwsEBSSnapshotAka,
			// 	Transform:   transform.FromValue(),
			// },
		}),
	}
}

//// LIST FUNCTION

func listAwsSSMParameters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listAwsSSMParameters", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := SsmService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeParametersPages(
		&ssm.DescribeParametersInput{},
		func(page *ssm.DescribeParametersOutput, isLast bool) bool {
			for _, parameter := range page.Parameters {
				d.StreamListItem(ctx, parameter)

			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSSMParameter(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSSMParameter")

	defaultRegion := GetDefaultRegion()
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Session
	svc, err := SsmService(ctx, d.ConnectionManager, defaultRegion)
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

// func getAwsEBSSnapshotAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("getAwsEBSSnapshotAka")
// 	snapshotData := h.Item.(*ec2.Snapshot)
// 	c, err := getCommonColumns(ctx, d, h)
// 	if err != nil {
// 		return nil, err
// 	}
// 	commonColumnData := c.(*awsCommonColumnData)

// 	// Get the resource akas
// 	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":snapshot/" + *snapshotData.SnapshotId}

// 	return akas, nil
// }

// //// TRANSFORM FUNCTIONS

// func ec2SnapshotTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	snapshot := d.HydrateItem.(*ec2.Snapshot)
// 	return ec2TagsToMap(snapshot.Tags)
// }
