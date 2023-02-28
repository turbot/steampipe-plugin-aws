package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice/types"

	directoryservicev1 "github.com/aws/aws-sdk-go/service/directoryservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDirectoryServiceDirectory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_directory_service_directory",
		Description: "AWS Directory Service Directory",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("directory_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValueException", "ResourceNotFoundFault", "EntityDoesNotExistException"}),
			},
			Hydrate: getDirectoryServiceDirectory,
		},
		List: &plugin.ListConfig{
			Hydrate: listDirectoryServiceDirectories,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "directory_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(directoryservicev1.EndpointsID),
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
				Hydrate:     getDirectoryARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "stage",
				Description: "The current stage of the directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The directory type.",
				Type:        proto.ColumnType_STRING,
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
				Name:        "radius_status",
				Description: "The status of the RADIUS MFA server connection.",
				Type:        proto.ColumnType_STRING,
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
				Name:        "connect_settings",
				Description: "A DirectoryConnectSettingsDescription object that contains additional information about an AD Connector directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "dns_ip_addrs",
				Description: "he IP addresses of the DNS servers for the directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "owner_directory_description",
				Description: "Describes the AWS Managed Microsoft AD directory in the directory owner account.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "radius_settings",
				Description: "A RadiusSettings object that contains information about the RADIUS server.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "regions_info",
				Description: "Lists the Regions where the directory has replicated.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "shared_directories",
				Description: "Details about the shared directory in the directory owner account for which the share request in the directory consumer account has been accepted.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDirectoryServiceSharedDirectory,
				Transform:   transform.FromValue().NullIfZero(),
			},
			{
				Name:        "vpc_settings",
				Description: "A DirectoryVpcSettingsDescription object that contains additional information about a directory.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags currently associated with the Directory Service Directory.",
				Hydrate:     getDirectoryServiceDirectoryTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
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
				Hydrate:     getDirectoryARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listDirectoryServiceDirectories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := DirectoryServiceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directory_service_directory.listDirectoryServiceDirectories", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	// Build the params
	input := &directoryservice.DescribeDirectoriesInput{
		Limit: aws.Int32(maxLimit),
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["directory_id"] != nil {
		input.DirectoryIds = []string{equalQuals["directory_id"].GetStringValue()}
	}

	pagesLeft := true

	// List call
	for pagesLeft {
		result, err := svc.DescribeDirectories(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_directory_service_directory.listDirectoryServiceDirectories", "api_error", err)
			return nil, err
		}

		for _, directory := range result.DirectoryDescriptions {
			d.StreamListItem(ctx, directory)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				pagesLeft = false
			}
		}

		if result.NextToken != nil {
			input.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDirectoryServiceDirectory(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := DirectoryServiceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directory_service_directory.getDirectoryServiceDirectory", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	directoryID := d.EqualsQuals["directory_id"].GetStringValue()
	if directoryID == "" {
		return nil, nil
	}

	params := &directoryservice.DescribeDirectoriesInput{
		DirectoryIds: []string{directoryID},
	}

	op, err := svc.DescribeDirectories(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directory_service_directory.getDirectoryServiceDirectory", "api_error", err)
		return nil, err
	}

	if op.DirectoryDescriptions != nil && len(op.DirectoryDescriptions) > 0 {
		return op.DirectoryDescriptions[0], nil
	}
	return nil, nil
}

func getDirectoryServiceSharedDirectory(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	directory := h.Item.(types.DirectoryDescription)

	// DescribeSharedDirectories Operation is only supported for MicrosoftAD directories.
	// Ignore if not a MicrosoftAD directory
	if directory.Type != types.DirectoryTypeMicrosoftAd {
		return nil, nil
	}

	// Create service
	svc, err := DirectoryServiceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directory_service_directory.getDirectoryServiceSharedDirectory", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &directoryservice.DescribeSharedDirectoriesInput{
		OwnerDirectoryId: directory.DirectoryId,
	}

	var directories []types.SharedDirectory

	for {
		response, err := svc.DescribeSharedDirectories(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_directory_service_directory.getDirectoryServiceSharedDirectory", "api_error", err)
			return nil, err
		}
		if response.SharedDirectories != nil {
			directories = append(directories, response.SharedDirectories...)
		}
		if response.NextToken == nil {
			break
		}
		params.NextToken = response.NextToken
	}
	return directories, nil
}

func getDirectoryARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	directory := h.Item.(types.DirectoryDescription)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	region := d.EqualsQualString(matrixKeyRegion)

	arn := "arn:" + commonColumnData.Partition + ":ds:" + region + ":" + commonColumnData.AccountId + ":directory/" + *directory.DirectoryId

	return arn, nil
}

func getDirectoryServiceDirectoryTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	directoryID := h.Item.(types.DirectoryDescription).DirectoryId

	// Create service
	svc, err := DirectoryServiceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directory_service_directory.getDirectoryServiceDirectoryTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &directoryservice.ListTagsForResourceInput{
		ResourceId: directoryID,
	}

	pagesLeft := true
	tags := []types.Tag{}

	for pagesLeft {
		result, err := svc.ListTagsForResource(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_directory_service_directory.getDirectoryServiceDirectoryTags", "api_error", err)
			return nil, err
		}
		tags = append(tags, result.Tags...)

		if result.NextToken != nil {
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return tags, nil
}

//// TRANSFORM FUNCTIONS

func directoryServiceDirectoryTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	if tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
