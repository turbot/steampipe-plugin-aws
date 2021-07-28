package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/directoryservice"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsDirectoryServiceDirectory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_directory_service_directory",
		Description: "AWS Directory Service Directory",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("directory_id"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidParameterValueException", "ResourceNotFoundFault"}),
			Hydrate:           getDirectoryServiceDirectory,
		},
		List: &plugin.ListConfig{
			Hydrate: listDirectiryServiceDirectorties,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The fully qualified name of the directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "directory_id",
				Description: "The directory identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies the directory.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDirectoryArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "access_url",
				Description: "The access URL for the directory, such as http://<alias>.awsapps.com.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "alias",
				Description: "The alias for the directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connect_settings",
				Description: "A DirectoryConnectSettingsDescription object that contains additional information about an AD Connector directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "description",
				Description: "The description for the directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "desired_number_of_domain_controllers",
				Description: "The desired number of domain controllers in the directory if the directory is Microsoft AD.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "dns_ip_addrs",
				Description: "he IP addresses of the DNS servers for the directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "edition",
				Description: "The edition associated with this directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "launch_time",
				Description: "Specifies when the directory was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "owner_directory_description",
				Description: "Describes the AWS Managed Microsoft AD directory in the directory owner account.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "radius_settings",
				Description: "A RadiusSettings object that contains information about the RADIUS server",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "radius_status",
				Description: "The status of the RADIUS MFA server connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "regions_info",
				Description: "Lists the Regions where the directory has replicated.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "share_method",
				Description: "The method used when sharing a directory to determine whether the directory should be shared within your AWS organization (ORGANIZATIONS) or with any AWS account by sending a shared directory request (HANDSHAKE).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "share_notes",
				Description: "A directory share request that is sent by the directory owner to the directory consumer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "share_status",
				Description: "Current directory status of the shared AWS Managed Microsoft AD directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "short_name",
				Description: "The short name of the directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "The directory size.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sso_enabled",
				Description: "Indicates if single sign-on is enabled for the directory. For more information, see EnableSso and DisableSso.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "stage",
				Description: "The current stage of the directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stage_last_updated_date_time",
				Description: "The date and time that the stage was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "stage_reason",
				Description: "Additional information about the directory stage.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The directory size.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_settings",
				Description: "A DirectoryVpcSettingsDescription object that contains additional information about a directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tag_src",
				Description: "A list of tags currently associated with the Directory Service Directory",
				Hydrate:     getDirectoryServiceDirectoryTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe Standard columns
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
				Hydrate:     getDirectoryServiceDirectoryTags,
				Transform:   transform.From(directoryServiceDirectoryTurbotData),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDirectoryArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDirectiryServiceDirectorties(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DirectoryServiceService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &directoryservice.DescribeDirectoriesInput{}

	// List call
	result, err := svc.DescribeDirectories(params)
	if err != nil {
		return nil, err
	}

	for _, directory := range result.DirectoryDescriptions {
		d.StreamListItem(ctx, directory)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDirectoryServiceDirectory(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := DirectoryServiceService(ctx, d)
	if err != nil {
		return nil, err
	}

	directoryID := d.KeyColumnQuals["directory_id"].GetStringValue()

	params := &directoryservice.DescribeDirectoriesInput{
		DirectoryIds: []*string{&directoryID},
	}

	op, err := svc.DescribeDirectories(params)
	if err != nil {
		return nil, err
	}

	if op.DirectoryDescriptions != nil && len(op.DirectoryDescriptions) > 0 {
		return op.DirectoryDescriptions[0], nil
	}
	return nil, nil
}

func getDirectoryArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDirectoryArn")
	directory := h.Item.(*directoryservice.DirectoryDescription)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ds:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":directory/" + *directory.DirectoryId

	return arn, nil
}

func getDirectoryServiceDirectoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDirectoryServiceDirectoryTags")

	directoryID := h.Item.(*directoryservice.DirectoryDescription).DirectoryId

	// Create service
	svc, err := DirectoryServiceService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &directoryservice.ListTagsForResourceInput{
		ResourceId: directoryID,
	}

	tags, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

//// TRANSFORMATION FUNCTION

func directoryServiceDirectoryTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*directoryservice.ListTagsForResourceOutput)
	if data.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	if data.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range data.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
