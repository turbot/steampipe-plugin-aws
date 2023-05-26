package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSTSCallerIdentity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sts_caller_identity",
		Description: "AWS STS Caller Identity",
		List: &plugin.ListConfig{
			Hydrate: getStsCallerIdentity,
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Web Services ARN associated with the calling entity.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_id",
				Description: "The unique identifier of the calling entity. The exact value depends on the type of entity that is making the call.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserId"),
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

//// GET FUNCTION

func getStsCallerIdentity(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	callerIdentity, err := getCallerIdentity(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sts_caller_identity.getStsCallerIdentity", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, callerIdentity)

	return nil, nil
}
