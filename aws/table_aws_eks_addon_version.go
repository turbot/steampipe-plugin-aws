package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/eks"
)

//// TABLE DEFINITION

func tableAwsEksAddonVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_addon_version",
		Description: "AWS EKS Addon Version",
		List: &plugin.ListConfig{
			Hydrate: listEksAddonVersions,
		},
		GetMatrixItem: BuildRegionList,
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
	Architecture    []*string
	Compatibilities []*eks.Compatibility
	Type            *string
}

//// LIST FUNCTION

func listEksAddonVersions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := EksService(ctx, d)
	if err != nil {
		return nil, err
	}

	err = svc.DescribeAddonVersionsPages(
		&eks.DescribeAddonVersionsInput{},
		func(page *eks.DescribeAddonVersionsOutput, _ bool) bool {
			for _, addon := range page.Addons {
				for _, version := range addon.AddonVersions {
					d.StreamListItem(ctx, addonVersion{addon.AddonName, version.AddonVersion, version.Architecture, version.Compatibilities, addon.Type})
				}
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAddonVersionAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAddonVersionAkas")
	version := h.Item.(addonVersion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":eks:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":addonversion/" + *version.AddonName + "/" + *version.AddonVersion}

	return akas, nil
}
