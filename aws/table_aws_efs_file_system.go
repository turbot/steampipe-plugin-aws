package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// TABLE DEFINITION

func tableAwsElasticFileSystem(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_efs_file_system",
		Description: "AWS Elastic File System",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("file_system_id"),
			ShouldIgnoreError: isNotFoundError([]string{"FileSystemNotFound"}),
			Hydrate:           getElasticFileSystem,
		},
		List: &plugin.ListConfig{
			Hydrate: listElasticFileSystem,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "file_system_id",
				Description: "The ID of the file system, assigned by Amazon EFS.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "file_system_arn",
				Description: "The Amazon Resource Name (ARN) for the EFS file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The AWS account that created the file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_token",
				Description: "The opaque string specified in the request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time that the file system was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "life_cycle_state",
				Description: "The lifecycle phase of the file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "number_of_mount_targets",
				Description: "The current number of mount targets that the file system has.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "size_in_bytes",
				Description: "The latest known metered size (in bytes) of data stored in the file system.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "performance_mode",
				Description: "The performance mode of the file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "encrypted",
				Description: "A Boolean value that, if true, indicates that the file system is encrypted.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of an AWS Key Management Service (AWS KMS) customer master key (CMK) that was used to protect the encrypted file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "throughput_mode",
				Description: "The throughput mode for a file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioned_throughput_in_mibps",
				Description: "The throughput, measured in MiB/s, that you want to provision for a file system.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "policy",
				Description: "The JSON formatted FileSystemPolicy for the EFS file system.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getElasticFileSystemPolicy,
				Transform:   transform.FromField("Policy"),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getElasticFileSystemPolicy,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with Filesystem.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FileSystemId"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(elasticFileSystemTurbotData, "Tags"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FileSystemArn").Transform(arnToAkas),
			},
		}),
	}
}

// LIST FUNCTION

func listElasticFileSystem(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := EFSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeFileSystemsPages(
		&efs.DescribeFileSystemsInput{},
		func(page *efs.DescribeFileSystemsOutput, isLast bool) bool {
			for _, fileSystem := range page.FileSystems {
				d.StreamListItem(ctx, fileSystem)
			}
			return !isLast
		},
	)

	return nil, err
}

// HYDRATE FUNCTIONS

func getElasticFileSystem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// create service
	svc, err := EFSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	fileSystemID := quals["file_system_id"].GetStringValue()

	params := &efs.DescribeFileSystemsInput{
		FileSystemId: aws.String(fileSystemID),
	}

	op, err := svc.DescribeFileSystems(params)
	if err != nil {
		return nil, err
	}

	if op.FileSystems != nil && len(op.FileSystems) > 0 {
		return op.FileSystems[0], nil
	}
	return nil, nil
}

func getElasticFileSystemPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getElasticFileSystemPolicy")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	fileSystem := h.Item.(*efs.FileSystemDescription)

	// Create session
	svc, err := EFSService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &efs.DescribeFileSystemPolicyInput{
		FileSystemId: fileSystem.FileSystemId,
	}

	fileSystemPolicy, err := svc.DescribeFileSystemPolicy(param)

	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "PolicyNotFound" {
				return nil, nil
			}
			return nil, err
		}
	}
	return fileSystemPolicy, nil
}

// TRANSFORM FUNCTIONS

func elasticFileSystemTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	fileSystemTag := d.HydrateItem.(*efs.FileSystemDescription)

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
