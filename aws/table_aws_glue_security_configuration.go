package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueSecurityConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_security_configuration",
		Description: "AWS Glue Security Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlueSecurityConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueSecurityConfigurations,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := GlueService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_security_configuration.listGlueSecurityConfigurations", "service_creation_error", err)
		return nil, err
	}

	input := &glue.GetSecurityConfigurationsInput{
		MaxResults: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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

	// List call
	err = svc.GetSecurityConfigurationsPages(
		input,
		func(page *glue.GetSecurityConfigurationsOutput, isLast bool) bool {
			for _, configuration := range page.SecurityConfigurations {
				d.StreamListItem(ctx, configuration)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_security_configuration.listGlueSecurityConfigurations", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueSecurityConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	// Create Session
	svc, err := GlueService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_security_configuration.getGlueSecurityConfiguration", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &glue.GetSecurityConfigurationInput{
		Name: aws.String(name),
	}

	// Get call
	data, err := svc.GetSecurityConfiguration(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_security_configuration.getGlueSecurityConfiguration", "api_error", err)
		return nil, err
	}
	return data.SecurityConfiguration, nil
}

func getGlueSecurityConfigurationArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*glue.SecurityConfiguration)

	// Get common columns
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// arn:aws:glue:region:account-id:security-configuration/configuration-name
	arn := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":security-configuration/" + *data.Name

	return arn, nil
}
