package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
)

//// TABLE DEFINITION

func tableAwsEksAddonVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_addon_version",
		Description: "AWS EKS Addon Version",
		List: &plugin.ListConfig{
			Hydrate: listEksAddonVersions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "addon_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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

	input := &eks.DescribeAddonVersionsInput{
		MaxResults: aws.Int64(100),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["addon_name"] != nil {
		input.AddonName = aws.String(equalQuals["addon_name"].GetStringValue())
	}

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

	err = svc.DescribeAddonVersionsPages(
		input,
		func(page *eks.DescribeAddonVersionsOutput, _ bool) bool {
			for _, addon := range page.Addons {
				for _, version := range addon.AddonVersions {
					d.StreamListItem(ctx, addonVersion{addon.AddonName, version.AddonVersion, version.Architecture, version.Compatibilities, addon.Type})

					// Context may get cancelled due to manual cancellation or if the limit has been reached
					if d.QueryStatus.RowsRemaining(ctx) == 0 {
						return false
					}
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
	region := d.KeyColumnQualString(matrixKeyRegion)
	version := h.Item.(addonVersion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":eks:" + region + ":" + commonColumnData.AccountId + ":addonversion/" + *version.AddonName + "/" + *version.AddonVersion}

	return akas, nil
}
