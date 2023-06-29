package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamOpenIdConnectProvider(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_open_id_connect_provider",
		Description: "AWS IAM OpenID Connect Provider",
		List: &plugin.ListConfig{
			Hydrate: listIamOpenIdConnectProviders,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"arn"}),
			Hydrate:    getIamOpenIdConnectProvider,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchEntity"}),
			},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the OIDC provider resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_id_list",
				Description: "A list of client IDs (also known as audiences) that are associated with the specified IAM OIDC provider resource object",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamOpenIdConnectProvider,
				Transform:   transform.FromField("ClientIDList"),
			},
			{
				Name:        "create_date",
				Description: "The date and time when the IAM OIDC provider resource object was created in the Amazon Web Services account.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getIamOpenIdConnectProvider,
			},
			{
				Name:        "thumbprint_list",
				Description: "A list of certificate thumbprints that are associated with the specified IAM OIDC provider resource object.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamOpenIdConnectProvider,
			},
			{
				Name:        "url",
				Description: "The URL that the IAM OIDC provider resource object is associated with.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIamOpenIdConnectProvider,
			},
			// Get and Tag API not listing tags for providers
			// {
			// 	Name:        "tags_src",
			// 	Description: "A list of tags that are attached to the specified IAM OIDC provider.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromField("Tags"),
			// 	Hydrate:     getIamOpenIdConnectProvider,
			// },

			// Standard columns for all tables
			// {
			// 	Name:        "tags",
			// 	Description: resourceInterfaceDescription("tags"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.From(openIDConnectTurbotTags),
			// 	Hydrate:     getIamOpenIdConnectProvider,
			// },
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
			},
		}),
	}
}

type OpenIDConnectProvider struct {
	Arn string `type:"string"`
	iam.GetOpenIDConnectProviderOutput
}

//// LIST FUNCTION

func listIamOpenIdConnectProviders(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_open_id_connect_provider.listIamOpenIdConnectProviders", "client_error", err)
		return nil, err
	}

	// SDK doesn't have new paginator for ListSAMLProviders action
	output, err := svc.ListOpenIDConnectProviders(ctx, &iam.ListOpenIDConnectProvidersInput{}, func(o *iam.Options) {
	})
	for _, provider := range output.OpenIDConnectProviderList {
		d.StreamListItem(ctx, provider)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamOpenIdConnectProvider(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		provider := h.Item.(types.OpenIDConnectProviderListEntry)
		arn = *provider.Arn
	} else {
		arn = d.EqualsQuals["arn"].GetStringValue()
	}

	if arn == "" {
		return nil, nil
	}

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_open_id_connect_provider.getIamOpenIdConnectProvider", "client_error", err)
		return nil, err
	}

	params := &iam.GetOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: aws.String(arn),
	}

	op, err := svc.GetOpenIDConnectProvider(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_open_id_connect_provider.getIamOpenIdConnectProvider", "api_error", err)
		return nil, err
	}

	return OpenIDConnectProvider{arn, *op}, nil
}

//// TRANSFORM FUNCTION

// func openIDConnectTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	tags := d.HydrateItem.(OpenIDConnectProvider)
// 	var turbotTagsMap map[string]string
// 	if len(tags.Tags) > 0 {
// 		turbotTagsMap = map[string]string{}
// 		for _, i := range tags.Tags {
// 			turbotTagsMap[*i.Key] = *i.Value
// 		}
// 	}
// 	return turbotTagsMap, nil
// }
