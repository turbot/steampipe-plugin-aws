package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsFsxFileSystem(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_fsx_file_system",
		Description: "AWS FSx File System",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("file_system_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"FileSystemNotFound", "ValidationException"}),
			},
			Hydrate: getFsxFileSystem,
		},
		List: &plugin.ListConfig{
			Hydrate: listFsxFileSystems,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "file_system_id",
				Description: "The system-generated, unique 17-digit ID of the file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the EFS file system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceARN"),
			},
			{
				Name:        "file_system_type",
				Description: "The type of Amazon FSx file system, which can be LUSTRE, WINDOWS, or ONTAP.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle",
				Description: "The lifecycle status of the file system, following are the possible values AVAILABLE, CREATING, DELETING, FAILED, MISCONFIGURED, UPDATING.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time that the file system was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "dns_name",
				Description: "The DNS name for the file system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DNSName"),
			},
			{
				Name:        "file_system_type_version",
				Description: "The version of your Amazon FSx for Lustre file system, either 2.10 or 2.12.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of the Key Management Service (KMS) key used to encrypt the file system's.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The AWS account that created the file system.",
				Type:        proto.ColumnType_STRING,
			},

			{
				Name:        "storage_capacity",
				Description: "The storage capacity of the file system in gibibytes (GiB).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "storage_type",
				Description: "The storage type of the file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the primary VPC for the file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administrative_actions",
				Description: "A list of administrative actions for the file system that are in process or waiting to be processed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "failure_details",
				Description: "A structure providing details of any failures that occur when creating the file system has failed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lustre_configuration",
				Description: "The configuration for the Amazon FSx for Lustre file system.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interface_ids",
				Description: "The IDs of the elastic network interface from which a specific file system is accessible.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ontap_configuration",
				Description: "The configuration for this FSx for NetApp ONTAP file system.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subnet_ids",
				Description: "Specifies the IDs of the subnets that the file system is accessible from.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with Filesystem.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "windows_configuration",
				Description: "The configuration for this Microsoft Windows file system.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getFsxFileSystemTurbotTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(fsxFileSystemTurbotData, "Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listFsxFileSystems(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := FsxService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listFsxFileSystem", "service_connection", err)
		return nil, err
	}

	// https://docs.aws.amazon.com/fsx/latest/APIReference/API_DescribeFileSystems.html
	input := &fsx.DescribeFileSystemsInput{
		MaxResults: aws.Int64(2147483647),
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
	err = svc.DescribeFileSystemsPages(
		input,
		func(page *fsx.DescribeFileSystemsOutput, isLast bool) bool {
			for _, fileSystem := range page.FileSystems {
				d.StreamListItem(ctx, fileSystem)

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

func getFsxFileSystem(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := FsxService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getFsxFileSystem", "service_connection", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	fileSystemID := quals["file_system_id"].GetStringValue()

	// Empty param check
	if fileSystemID == "" {
		return nil, nil
	}

	params := &fsx.DescribeFileSystemsInput{
		FileSystemIds: []*string{aws.String(fileSystemID)},
	}

	op, err := svc.DescribeFileSystems(params)
	if err != nil {
		plugin.Logger(ctx).Error("getFsxFileSystem", err)
		return nil, err
	}

	if op.FileSystems != nil && len(op.FileSystems) > 0 {
		return op.FileSystems[0], nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func fsxFileSystemTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	fileSystemTag := d.HydrateItem.(*fsx.FileSystem)
	if fileSystemTag.Tags == nil {
		return nil, nil
	}

	// Get the resource tags
	var turbotTagsMap map[string]string
	if fileSystemTag.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range fileSystemTag.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}

func getFsxFileSystemTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	fileSystemTitle := d.HydrateItem.(*fsx.FileSystem)

	if fileSystemTitle.Tags != nil {
		for _, i := range fileSystemTitle.Tags {
			if *i.Key == "Name" && len(*i.Value) > 0 {
				return *i.Value, nil
			}
		}
	}

	return fileSystemTitle.FileSystemId, nil
}
