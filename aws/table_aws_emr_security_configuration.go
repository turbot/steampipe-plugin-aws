package aws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/emr/types"

	emrv1 "github.com/aws/aws-sdk-go/service/emr"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEmrSecurityConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_emr_security_configuration",
		Description: "AWS EMR Security Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidRequestException"}),
			},
			Hydrate: getEmrSecurityConfiguration,
			Tags:    map[string]string{"service": "elasticmapreduce", "action": "DescribeSecurityConfiguration"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEmrSecurityConfigurations,
			Tags:    map[string]string{"service": "elasticmapreduce", "action": "ListSecurityConfigurations"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(emrv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "TThe name of the security configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date_time",
				Description: "The date and time the security configuration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "encryption_configuration",
				Description: "The encryption configuration details for a secutiry configuration.",
				Hydrate:     getEmrSecurityConfiguration,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instance_metadata_service_configuration",
				Description: "The instance metadata service configuration details for a secutiry configuration.",
				Hydrate:     getEmrSecurityConfiguration,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_configuration",
				Description: "The security configuration details in JSON format.",
				Hydrate:     getEmrSecurityConfiguration,
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

type ConfigurationDetails struct {
	Name                                 *string
	CreationDateTime                     *time.Time
	EncryptionConfiguration              interface{}
	InstanceMetadataServiceConfiguration interface{}
	SecurityConfiguration                interface{}
}

//// LIST FUNCTION

func listEmrSecurityConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EMRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_security_configuration.listEmrSecurityConfigurations", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &emr.ListSecurityConfigurationsInput{}

	paginator := emr.NewListSecurityConfigurationsPaginator(svc, input, func(o *emr.ListSecurityConfigurationsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_emr_security_configuration.listEmrSecurityConfigurations", "api_error", err)
			return nil, err
		}

		for _, cfg := range output.SecurityConfigurations {
			d.StreamListItem(ctx, cfg)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEmrSecurityConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		cfg := h.Item.(types.SecurityConfigurationSummary)
		name = *cfg.Name
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	// Create service
	svc, err := EMRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_security_configuration.getEmrSecurityConfiguration", "connection_error", err)
		return nil, err
	}

	params := &emr.DescribeSecurityConfigurationInput{
		Name: aws.String(name),
	}

	op, err := svc.DescribeSecurityConfiguration(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_security_configuration.getEmrSecurityConfiguration", "api_error", err)
		return nil, err
	}

	if op != nil {
		var v map[string]interface{}
		if err := json.Unmarshal([]byte(*op.SecurityConfiguration), &v); err != nil {
			return nil, err
		}

		return ConfigurationDetails{op.Name, op.CreationDateTime, v["EncryptionConfiguration"], v["InstanceMetadataServiceConfiguration"], op.SecurityConfiguration}, nil
	}

	return nil, nil
}


{
   "AtRestEncryptionConfiguration": {
    "LocalDiskEncryptionConfiguration": {
     "AwsKmsKey": "arn:aws:kms:us-east-1:632902152528:alias/my-key-alias_turbottest69288",
     "EnableEbsEncryption": true,
     "EncryptionKeyProviderType": "AwsKms"
    },
    "S3EncryptionConfiguration": {
     "AwsKmsKey": "arn:aws:kms:us-east-1:632902152528:alias/my-key-alias_turbottest69288",
     "EncryptionMode": "SSE-KMS"
    }
   },
   "EnableAtRestEncryption": true,
   "EnableInTransitEncryption": true,
   "InTransitEncryptionConfiguration": {
    "TLSCertificateConfiguration": {
     "CertificateProviderType": "PEM",
     "S3Object": "s3://emr-test-1111/azure_thrifty_ss.zip"
    }
   }
  }