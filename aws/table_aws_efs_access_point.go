package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEfsAccessPoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_efs_access_point",
		Description: "AWS EFS Access Point",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("access_point_id"),
			ShouldIgnoreError: isNotFoundError([]string{"AccessPointNotFound"}),
			Hydrate:           getEfsAccessPoint,
		},
		List: &plugin.ListConfig{
			Hydrate: listEfsAccessPoints,
		},
		GetMatrixItem: BuildRegionList,
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
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := EfsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeAccessPointsPages(
		&efs.DescribeAccessPointsInput{},
		func(page *efs.DescribeAccessPointsOutput, isLast bool) bool {
			for _, accessPoint := range page.AccessPoints {
				d.StreamListItem(ctx, accessPoint)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEfsAccessPoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	quals := d.KeyColumnQuals
	accessPointID := quals["access_point_id"].GetStringValue()

	// create service
	svc, err := EfsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &efs.DescribeAccessPointsInput{
		AccessPointId: aws.String(accessPointID),
	}

	data, err := svc.DescribeAccessPoints(params)
	if err != nil {
		return nil, err
	}

	if data.AccessPoints != nil && len(data.AccessPoints) > 0 {
		return data.AccessPoints[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func efsAccessPointTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("efsAccessPointTurbotTags")
	tagList := d.HydrateItem.(*efs.AccessPointDescription)

	if tagList.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

// Generate title for the resource
func efsAccessPointTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("efsAccessPointTitle")
	data := d.HydrateItem.(*efs.AccessPointDescription)

	// If name is available, then setting name as title, else setting Access Point ID as title
	if data.Name != nil {
		return data.Name, nil
	}

	return data.AccessPointId, nil
}
