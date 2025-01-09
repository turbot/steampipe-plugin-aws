package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"

	ssoadminEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSsoAdminInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssoadmin_instance",
		Description: "AWS SSO Instance",
		List: &plugin.ListConfig{
			Hydrate: listSsoAdminInstances,
			Tags:    map[string]string{"service": "sso", "action": "ListInstances"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssoadminEndpoint.SSOServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The ARN of the SSO instance under which the operation will be executed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceArn"),
			},
			{
				Name:        "identity_store_id",
				Description: "The identifier of the identity store that is connected to the SSO instance.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceArn"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InstanceArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSsoAdminInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SSOAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssoadmin_instance.listSsoAdminInstances", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &ssoadmin.ListInstancesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := ssoadmin.NewListInstancesPaginator(svc, input, func(o *ssoadmin.ListInstancesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssoadmin_instance.listSsoAdminInstances", "api_error", err)
			return nil, err
		}

		for _, items := range output.Instances {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}
