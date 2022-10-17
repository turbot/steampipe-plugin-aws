package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSageMakerDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_sagemaker_domain",
		Description: "AWS Sagemaker Domain",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ValidationException", "NotFoundException", "ResourceNotFound"}),
			},
			Hydrate: getAwsSageMakerDomain,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSageMakerDomains,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The domain ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainId"),
			},
			{
				Name:        "name",
				Description: "The domain name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainArn"),
			},
			{
				Name:        "creation_time",
				Description: "A timestamp that indicates when the domain was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			{
				Name:        "app_network_access_type",
				Description: "Specifies the VPC used for non-EFS traffic.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "app_security_group_management",
				Description: "The entity that creates and manages the required security groups for inter-app communication in VPCOnly mode.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "auth_mode",
				Description: "The domain's authentication mode.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "failure_reason",
				Description: "The domain's failure reason.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "home_efs_file_system_id",
				Description: "The ID of the Amazon Elastic File System (EFS) managed by this domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "kms_key_id",
				Description: "The Amazon Web Services KMS customer managed key used to encrypt the EFS volume attached to the domain.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "last_modified_time",
				Description: "The domain's last modified time.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "security_group_id_for_domain_boundary",
				Description: "The ID of the security group that authorizes traffic between the RSessionGateway apps and the RStudioServerPro app.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "single_sign_on_managed_application_instance_id",
				Description: "The SSO managed application instance ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "status",
				Description: "The domain's status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the domain.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerDomainTags,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "url",
				Description: "The domain's URL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_user_settings",
				Description: "Settings which are applied to UserProfiles in this domain if settings are not explicitly specified in a given UserProfile.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "domain_settings",
				Description: "A collection of domain settings.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerDomain,
			},
			{
				Name:        "subnet_ids",
				Description: "The VPC subnets that studio uses for communication.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSageMakerDomain,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsSageMakerDomainTags,
				Transform:   transform.FromValue().Transform(sageMakerTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DomainArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSageMakerDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_domain.listAwsSageMakerDomains", "connection_error", err)
		return nil, err
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

	input := &sagemaker.ListDomainsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := sagemaker.NewListDomainsPaginator(svc, input, func(o *sagemaker.ListDomainsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "no such host") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_sagemaker_domain.listAwsSageMakerDomains", "api_error", err)
			return nil, err
		}

		for _, items := range output.Domains {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSageMakerDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		id = sageMakerDomainId(h.Item)
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}
	if id == "" {
		return nil, nil
	}

	// Create service
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_domain.getAwsSageMakerDomain", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &sagemaker.DescribeDomainInput{
		DomainId: aws.String(id),
	}

	// Get call
	data, err := svc.DescribeDomain(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_sagemaker_domain.getAwsSageMakerDomain", "api_error", err)
		return nil, err
	}
	return data, nil
}

func listAwsSageMakerDomainTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var domainArn string
	if h.Item != nil {
		domainArn = sageMakerDomainArn(h.Item)
	}

	// Create Session
	svc, err := SageMakerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_sagemaker_domain.listAwsSageMakerDomainTags", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &sagemaker.ListTagsInput{
		ResourceArn: aws.String(domainArn),
	}

	pagesLeft := true
	tags := []types.Tag{}
	for pagesLeft {
		keyTags, err := svc.ListTags(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_sagemaker_domain.listAwsSageMakerDomainTags", "api_error", err)
			return nil, err
		}
		tags = append(tags, keyTags.Tags...)

		if keyTags.NextToken != nil {
			params.NextToken = keyTags.NextToken
		} else {
			pagesLeft = false
		}
	}

	return tags, nil
}

//// TRANSFORM FUNCTION

func sageMakerDomainId(item interface{}) string {
	switch item := item.(type) {
	case types.DomainDetails:
		return *item.DomainId
	case *sagemaker.DescribeDomainOutput:
		return *item.DomainId
	}
	return ""
}

func sageMakerDomainArn(item interface{}) string {
	switch item := item.(type) {
	case types.DomainDetails:
		return *item.DomainArn
	case *sagemaker.DescribeDomainOutput:
		return *item.DomainArn
	}
	return ""
}
