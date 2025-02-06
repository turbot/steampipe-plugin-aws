package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"

	gluev1 "github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueSecurityConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_security_configuration",
		Description: "AWS Glue Security Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlueSecurityConfiguration,
			Tags:    map[string]string{"service": "glue", "action": "GetSecurityConfiguration"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueSecurityConfigurations,
			Tags:    map[string]string{"service": "glue", "action": "GetSecurityConfigurations"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(gluev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the security configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time_stamp",
				Description: "The time at which this security configuration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "cloud_watch_encryption",
				Description: "The encryption configuration for Amazon CloudWatch.",
				Transform:   transform.FromField("EncryptionConfiguration.CloudWatchEncryption"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "job_bookmarks_encryption",
				Description: "The encryption configuration for job bookmarks.",
				Transform:   transform.FromField("EncryptionConfiguration.JobBookmarksEncryption"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "s3_encryption",
				Description: "The encryption configuration for Amazon Simple Storage Service (Amazon S3) data.",
				Transform:   transform.FromField("EncryptionConfiguration.S3Encryption"),
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTagsForGlueSecurityConfiguration,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueSecurityConfigurationArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueSecurityConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_security_configuration.listGlueSecurityConfigurations", "connection_error", err)
		return nil, err
	}
	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(200)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	input := &glue.GetSecurityConfigurationsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// List call
	paginator := glue.NewGetSecurityConfigurationsPaginator(svc, input, func(o *glue.GetSecurityConfigurationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glue_security_configuration.listGlueSecurityConfigurations", "api_error", err)
			return nil, err
		}
		for _, configuration := range output.SecurityConfigurations {
			d.StreamListItem(ctx, configuration)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueSecurityConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_security_configuration.getGlueSecurityConfiguration", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &glue.GetSecurityConfigurationInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.GetSecurityConfiguration(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_security_configuration.getGlueSecurityConfiguration", "api_error", err)
		return nil, err
	}
	return *data.SecurityConfiguration, nil
}

func getTagsForGlueSecurityConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	arn, _ := getGlueSecurityConfigurationArn(ctx, d, h)
	return getTagsForGlueResource(ctx, d, arn.(string))
}

func getGlueSecurityConfigurationArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	data := h.Item.(types.SecurityConfiguration)

	// Get common columns

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_security_configuration.getGlueSecurityConfigurationArn", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// arn:aws:glue:region:account-id:security-configuration/configuration-name
	arn := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":security-configuration/" + *data.Name

	return arn, nil
}
