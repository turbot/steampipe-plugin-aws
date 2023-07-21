package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/fms"
	"github.com/aws/aws-sdk-go-v2/service/fms/types"

	fmsv1 "github.com/aws/aws-sdk-go/service/fms"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsFMSPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_fms_policy",
		Description: "AWS FMS Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("file_system_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"FileSystemNotFound", "ValidationException"}),
			},
			Hydrate: getFmsPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listFmsPolicies,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(fmsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "PolicyName",
				Description: "The name of the specified policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "PolicyId",
				Description: "The ID of the specified policy.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("PolicyArn"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the specified policy.",
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

type PolicyInfo struct {
    ExcludeResourceTags bool
		PolicyArn *string
    PolicyName *string
    RemediationEnabled bool
    ResourceType *string
    SecurityServicePolicyData *types.SecurityServicePolicyData
    DeleteUnusedFMManagedResources bool
    ExcludeMap map[string][]string
    IncludeMap map[string][]string
    PolicyDescription *string
    PolicyId *string
    PolicyStatus types.CustomerPolicyStatus
    PolicyUpdateToken *string
    ResourceSetIds []string
    ResourceTags []types.ResourceTag
    ResourceTypeList []string
}

//// LIST FUNCTION

func listFmsPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := FMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fms_policy.listFmsPolicies", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(100)
	input := fms.ListPoliciesInput{
		MaxResults: aws.Int32(maxItems),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	paginator := fms.NewListPoliciesPaginator(svc, &input, func(o *fms.ListPoliciesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_fms_policy.listFmsPolicies", "api_error", err)
			return nil, err
		}

		for _, policy := range output.PolicyList {
			d.StreamListItem(ctx, policy)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getFmsPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	policyId := ""

	if h.Item != nil {
		data := h.Item.(types.PolicySummary)
		policyId = *data.PolicyId
	} else {
		policyId = d.EqualsQualString("policy_id")
	}

	if policyId == "" {
		return nil, nil
	}
	// Create service
	svc, err := FMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fms_policy.getFmsPolicy", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &fms.GetPolicyInput{
		PolicyId: &policyId,
	}

	op, err := svc.GetPolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_fsx_file_system.getFsxFileSystem", "api_error", err)
		return nil, err
	}

	return op, nil
}

// //// TRANSFORM FUNCTIONS

// func fsxFileSystemTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	fileSystemTag := d.HydrateItem.(types.FileSystem)
// 	if fileSystemTag.Tags == nil {
// 		return nil, nil
// 	}

// 	// Get the resource tags
// 	var turbotTagsMap map[string]string
// 	if fileSystemTag.Tags != nil {
// 		turbotTagsMap = map[string]string{}
// 		for _, i := range fileSystemTag.Tags {
// 			turbotTagsMap[*i.Key] = *i.Value
// 		}
// 	}
// 	return turbotTagsMap, nil
// }

// func getFsxFileSystemTurbotTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	fileSystemTitle := d.HydrateItem.(types.FileSystem)

// 	if fileSystemTitle.Tags != nil {
// 		for _, i := range fileSystemTitle.Tags {
// 			if *i.Key == "Name" && len(*i.Value) > 0 {
// 				return *i.Value, nil
// 			}
// 		}
// 	}

// 	return fileSystemTitle.FileSystemId, nil
// }
