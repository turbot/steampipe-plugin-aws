package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSSMManagedInstancePatchState(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_managed_instance_patch_state",
		Description: "AWS SSM Managed Instance Patch State",
		List: &plugin.ListConfig{
			ParentHydrate: listSsmManagedInstances,
			Hydrate:       listSsmManagedInstancePatchStates,
			Tags:          map[string]string{"service": "ssm", "action": "DescribeInstancePatchStates"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "instance_id", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SSM_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "instance_id",
				Description: "The ID of the managed node the high-level patch compliance information was collected for.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "baseline_id",
				Description: "The ID of the patch baseline used to patch the managed node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operation",
				Description: "The type of patching operation that was performed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operation_end_time",
				Description: "The time the most recent patching operation completed on the managed node.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "operation_start_time",
				Description: "The time the most recent patching operation was started on the managed node.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "patch_group",
				Description: "The name of the patch group the managed node belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "critical_non_compliant_count",
				Description: "The number of patches per node that are specified as Critical for compliance reporting in the patch baseline aren't installed. These patches might be missing, have failed installation, were rejected, or were installed but awaiting a required managed node reboot. The status of these managed nodes is NON_COMPLIANT.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "failed_count",
				Description: "The number of patches from the patch baseline that were attempted to be installed during the last patching operation, but failed to install.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "installed_count",
				Description: "The number of patches from the patch baseline that are installed on the managed node.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "installed_other_count",
				Description: "The number of patches not specified in the patch baseline that are installed on the managed node.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "installed_pending_reboot_count",
				Description: "The number of patches installed by Patch Manager since the last time the managed node was rebooted.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "installed_rejected_count",
				Description: "The number of patches installed on a managed node that are specified in a RejectedPatches list. Patches with a status of InstalledRejected were typically installed before they were added to a RejectedPatches list.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "last_no_reboot_install_operation_time",
				Description: "The time of the last attempt to patch the managed node with NoReboot specified as the reboot option.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "missing_count",
				Description: "The number of patches from the patch baseline that are applicable for the managed node but aren't currently installed.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "not_applicable_count",
				Description: "The number of patches from the patch baseline that aren't applicable for the managed node and therefore aren't installed on the node. This number may be truncated if the list of patch names is very large. The number of patches beyond this limit are reported in UnreportedNotApplicableCount.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "install_override_list",
				Description: "An https URL or an Amazon Simple Storage Service (Amazon S3) path-style URL to a list of patches to be installed.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "other_non_compliant_count",
				Description: "The number of patches per node that are specified as other than Critical or Security but aren't compliant with the patch baseline. The status of these managed nodes is NON_COMPLIANT.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "owner_information",
				Description: "Placeholder information. This field will always be empty in the current release of the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reboot_option",
				Description: "Indicates the reboot option specified in the patch baseline. Reboot options apply to Install operations only. Reboots aren't attempted for Patch Manager Scan operations.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "security_non_compliant_count",
				Description: "The number of patches per node that are specified as Security in a patch advisory aren't installed. These patches might be missing, have failed installation, were rejected, or were installed but awaiting a required managed node reboot. The status of these managed nodes is NON_COMPLIANT.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "snapshot_id",
				Description: "The ID of the patch baseline snapshot used during the patching operation when this compliance data was collected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "unreported_not_applicable_count",
				Description: "The number of patches beyond the supported limit of NotApplicableCount that aren't reported by name to Inventory. Inventory is a capability of Amazon Web Services Systems Manager.",
				Type:        proto.ColumnType_INT,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BaselineId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listSsmManagedInstancePatchStates(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	instance := h.Item.(types.InstanceInformation)

	instanceId := d.EqualsQualString("instance_id")

	// Restrict the API call for other intances if an instance ID has specified in the query parameter.
	if instanceId != "" && instanceId != *instance.InstanceId {
		return nil, nil
	}

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_managed_instance_patch_state.listSsmManagedInstancePatchStates", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Though the doc says the maxResults value can be between 10-100 but the API throws error ValidationException: 1 validation error detected: Value '100' at 'maxResults' failed to satisfy constraint: Member must have value less than or equal to 50
	maxItems := int32(50)
	input := &ssm.DescribeInstancePatchStatesInput{
		InstanceIds: []string{*instance.InstanceId},
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 10 {
				maxItems = int32(10)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := ssm.NewDescribeInstancePatchStatesPaginator(svc, input, func(o *ssm.DescribeInstancePatchStatesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_managed_instance_patch_state.listSsmManagedInstancePatchStates", "api_error", err)
			return nil, err
		}

		for _, patchState := range output.InstancePatchStates {
			d.StreamListItem(ctx, patchState)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
