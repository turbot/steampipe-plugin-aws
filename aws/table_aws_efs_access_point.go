package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEfsAccessPoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_efs_access_point",
		Description: "AWS EFS Access Point",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("access_point_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"AccessPointNotFound"}),
			},
			Hydrate: getEfsAccessPoint,
		},
		List: &plugin.ListConfig{
			Hydrate: listEfsAccessPoints,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"FileSystemNotFound"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "file_system_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the access point. This is the value of the Name tag.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "access_point_id",
				Description: "The ID of the access point, assigned by Amazon EFS.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "access_point_arn",
				Description: "The unique Amazon Resource Name (ARN) associated with the access point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "life_cycle_state",
				Description: "Identifies the lifecycle phase of the access point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "file_system_id",
				Description: "The ID of the EFS file system that the access point applies to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_token",
				Description: "The opaque string specified in the request to ensure idempotent creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "Identified the AWS account that owns the access point resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "posix_user",
				Description: "The full POSIX identity, including the user ID, group ID, and secondary group IDs on the access point that is used for all file operations by NFS clients using the access point.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "root_directory",
				Description: "The directory on the Amazon EFS file system that the access point exposes as the root directory to NFS clients using the access point.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "The tags associated with the access point, presented as an array of Tag objects.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(efsAccessPointTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(efsAccessPointTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AccessPointArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listEfsAccessPoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EFSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_efs_access_point.listEfsAccessPoints", "sconnection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxLimit := int32(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 5 {
				maxLimit = 5
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	input := &efs.DescribeAccessPointsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["file_system_id"] != nil {
		input.FileSystemId = aws.String(equalQuals["file_system_id"].GetStringValue())
	}
	paginator := efs.NewDescribeAccessPointsPaginator(svc, input, func(o *efs.DescribeAccessPointsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_efs_access_point.listEfsAccessPoints", "api_error", err)
			return nil, err
		}
		for _, accessPoint := range output.AccessPoints {
			d.StreamListItem(ctx, accessPoint)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEfsAccessPoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	accessPointID := quals["access_point_id"].GetStringValue()

	// create service
	svc, err := EFSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_efs_access_point.getEfsAccessPoint", "sconnection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &efs.DescribeAccessPointsInput{
		AccessPointId: aws.String(accessPointID),
	}

	data, err := svc.DescribeAccessPoints(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_efs_access_point.getEfsAccessPoint", "api_error", err)
		return nil, err
	}

	if data.AccessPoints != nil && len(data.AccessPoints) > 0 {
		return data.AccessPoints[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func efsAccessPointTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(types.AccessPointDescription)

	if tagList.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string

	if tagList.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

// Generate title for the resource
func efsAccessPointTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(types.AccessPointDescription)

	// If name is available, then setting name as title, else setting Access Point ID as title
	if data.Name != nil {
		return data.Name, nil
	}

	return data.AccessPointId, nil
}
