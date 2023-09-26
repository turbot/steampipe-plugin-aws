package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamSamlProvider(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_saml_provider",
		Description: "AWS IAM Saml Provider",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"arn"}),
			Hydrate:    getIamSamlProvider,
			Tags:       map[string]string{"service": "iam", "action": "GetSAMLProvider"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchEntity"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listIamSamlProviders,
			Tags:    map[string]string{"service": "iam", "action": "ListSAMLProviders"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the IAM policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The date and time when the SAML provider was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "valid_until",
				Description: "The expiration date and time for the SAML provider.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "saml_metadata_document",
				Description: "The XML metadata document that includes information about an identity provider.",
				Hydrate:     getIamSamlProvider,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SAMLMetadataDocument"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the specified IAM SAML provider.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamSamlProvider,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamSamlProvider,
				Transform:   transform.From(samlProviderTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type SAMLProvider struct {
	Arn                  *string     `min:"20" type:"string"`
	CreateDate           *time.Time  `type:"timestamp"`
	SAMLMetadataDocument *string     `min:"1000" type:"string"`
	Tags                 []types.Tag `type:"list"`
	ValidUntil           *time.Time  `type:"timestamp"`
}

//// LIST FUNCTION

func listIamSamlProviders(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_saml_provider.listIamSamlProviders", "service_creation_error", err)
		return nil, err
	}

	params := &iam.ListSAMLProvidersInput{}

	// List call
	// SDK doesn't have new paginator for ListSAMLProviders action
	result, err := svc.ListSAMLProviders(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_saml_provider.listIamSamlProviders", "api_error", err)
		return nil, err
	}

	for _, row := range result.SAMLProviderList {
		d.StreamListItem(ctx, SAMLProvider{
			Arn:        row.Arn,
			CreateDate: row.CreateDate,
			ValidUntil: row.ValidUntil,
		})

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTION

func getIamSamlProvider(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		entry := h.Item.(SAMLProvider)
		arn = *entry.Arn
	} else {
		arn = d.EqualsQuals["arn"].GetStringValue()
	}

	if arn == "" {
		return nil, nil
	}

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_saml_provider.getIamSamlProvider", "client_error", err)
		return nil, err
	}

	params := &iam.GetSAMLProviderInput{
		SAMLProviderArn: aws.String(arn),
	}

	// List call
	result, err := svc.GetSAMLProvider(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_saml_provider.getIamSamlProvider", "api_error", err)
		return nil, err
	}

	provider := SAMLProvider{
		Arn:                  aws.String(arn),
		CreateDate:           result.CreateDate,
		ValidUntil:           result.ValidUntil,
		SAMLMetadataDocument: result.SAMLMetadataDocument,
		Tags:                 result.Tags,
	}

	return provider, nil
}

//// TRANSFORM FUNCTION

func samlProviderTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	provider := d.HydrateItem.(SAMLProvider)
	if len(provider.Tags) == 0 {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range provider.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}
	return turbotTagsMap, nil
}
