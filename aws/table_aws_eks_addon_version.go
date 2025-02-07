package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"

	eksEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEksAddonVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_addon_version",
		Description: "AWS EKS Addon Version",
		List: &plugin.ListConfig{
			Hydrate: listEksAddonVersions,
			Tags:    map[string]string{"service": "eks", "action": "DescribeAddonVersions"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "addon_name", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getEksAddonConfiguration,
				Tags: map[string]string{"service": "eks", "action": "DescribeAddonConfiguration"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(eksEndpoint.AWS_EKS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "addon_name",
				Description: "The name of the add-on.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "addon_version",
				Description: "The version of the add-on.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the add-on.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "addon_configuration",
				Description: "The configuration for the add-on.",
				Hydrate:     getEksAddonConfiguration,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "architecture",
				Description: "The architectures that the version supports.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "compatibilities",
				Description: "An object that represents the compatibilities of a version.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AddonVersion"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAddonVersionAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type addonVersion struct {
	AddonName       *string
	AddonVersion    *string
	Architecture    []string
	Compatibilities []types.Compatibility
	Type            *string
}

//// LIST FUNCTION

func listEksAddonVersions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_addon_version.listEksAddonVersions", "get_client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &eks.DescribeAddonVersionsInput{
		MaxResults: aws.Int32(100),
	}

	equalQuals := d.EqualsQuals
	if equalQuals["addon_name"] != nil {
		input.AddonName = aws.String(equalQuals["addon_name"].GetStringValue())
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxResults {
			if limit < 1 {
				input.MaxResults = aws.Int32(1)
			} else {
				input.MaxResults = aws.Int32(limit)
			}
		}
	}

	paginator := eks.NewDescribeAddonVersionsPaginator(svc, input, func(o *eks.DescribeAddonVersionsPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_addon_version.listEksAddonVersions", "api_error", err)
			return nil, err
		}

		for _, addon := range output.Addons {
			for _, version := range addon.AddonVersions {
				d.StreamListItem(ctx, addonVersion{addon.AddonName, version.AddonVersion, version.Architecture, version.Compatibilities, addon.Type})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEksAddonConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	data := h.Item.(addonVersion)
	version := data.AddonVersion
	addonName := data.AddonName

	// create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_addon_version.getEksAddonConfiguration", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &eks.DescribeAddonConfigurationInput{
		AddonName:    addonName,
		AddonVersion: version,
	}

	op, err := svc.DescribeAddonConfiguration(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_addon_version.getEksAddonConfiguration", "api_error", err)
		return nil, err
	}

	return op.ConfigurationSchema, nil
}

func getAddonVersionAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	version := h.Item.(addonVersion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_addon_version.getAddonVersionAkas", "api_error", err)
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":eks:" + region + ":" + commonColumnData.AccountId + ":addonversion/" + *version.AddonName + "/" + *version.AddonVersion}

	return akas, nil
}
