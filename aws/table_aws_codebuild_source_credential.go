package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCodeBuildSourceCredential(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_codebuild_source_credential",
		Description: "AWS CodeBuild Source Credential",
		List: &plugin.ListConfig{
			Hydrate: listCodeBuildSourceCredentials,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the token.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auth_type",
				Description: "The type of authentication used by the credentials. Possible values are: OAUTH, BASIC_AUTH, or PERSONAL_ACCESS_TOKEN.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server_type",
				Description: "The type of source provider. Possible values are: GITHUB, GITHUB_ENTERPRISE, or BITBUCKET.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(codebuildSourceCredentialTitle),
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

//// LIST FUNCTION

func listCodeBuildSourceCredentials(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CodeBuildService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	resp, err := svc.ListSourceCredentials(&codebuild.ListSourceCredentialsInput{})
	if err != nil {
		return nil, err
	}
	for _, cred := range resp.SourceCredentialsInfos {
		d.StreamListItem(ctx, cred)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func codebuildSourceCredentialTitle(_ context.Context, d *transform.TransformData) (interface{},
	error) {
	data := d.HydrateItem.(*codebuild.SourceCredentialsInfo)

	splitPart := strings.Split(*data.Arn, ":")
	title := *data.ServerType + " - " + splitPart[3]

	return title, nil
}
