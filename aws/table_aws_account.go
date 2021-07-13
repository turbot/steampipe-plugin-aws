package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/organizations"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_account",
		Description: "AWS Account",
		List: &plugin.ListConfig{
			Hydrate: listAccountAlias,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "account_aliases",
				Description: "A list of aliases associated with the account, if applicable.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Aliases"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(accountARN),
			},
			{
				Name:        "organization_id",
				Description: "The unique identifier (ID) of an organization, if applicable.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrganizationDetails,
				Transform:   transform.FromField("Organization.Id"),
			},
			{
				Name:        "organization_arn",
				Description: "The Amazon Resource Name (ARN) of an organization.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrganizationDetails,
				Transform:   transform.FromField("Organization.Arn"),
			},
			{
				Name:        "organization_feature_set",
				Description: "Specifies the functionality that currently is available to the organization. If set to \"ALL\", then all features are enabled and policies can be applied to accounts in the organization. If set to \"CONSOLIDATED_BILLING\", then only consolidated billing functionality is available.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrganizationDetails,
				Transform:   transform.FromField("Organization.FeatureSet"),
			},
			{
				Name:        "organization_master_account_arn",
				Description: "The Amazon Resource Name (ARN) of the account that is designated as the management account for the organization",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrganizationDetails,
				Transform:   transform.FromField("Organization.MasterAccountArn"),
			},
			{
				Name:        "organization_master_account_email",
				Description: "The email address that is associated with the AWS account that is designated as the management account for the organization",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrganizationDetails,
				Transform:   transform.FromField("Organization.MasterAccountEmail"),
			},
			{
				Name:        "organization_master_account_id",
				Description: "The unique identifier (ID) of the management account of an organization",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getOrganizationDetails,
				Transform:   transform.FromField("Organization.MasterAccountId"),
			},
			{
				Name:        "organization_available_policy_types",
				Description: "The Region opt-in status. The possible values are opt-in-not-required, opted-in, and not-opted-in",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOrganizationDetails,
				Transform:   transform.FromField("Organization.AvailablePolicyTypes"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(accountDataToTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(accountARN).Transform(transform.EnsureStringArray),
			},
		}),
	}
}

type accountData struct {
	commonColumnData awsCommonColumnData
	Aliases          []*string
}

//// LIST FUNCTION

func listAccountAlias(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// execute list call
	op, err := svc.ListAccountAliases(&iam.ListAccountAliasesInput{})
	if err != nil {
		return nil, err
	}

	if op.AccountAliases != nil {
		d.StreamListItem(ctx, &accountData{
			commonColumnData: *commonColumnData,
			Aliases:          op.AccountAliases,
		})
		return nil, nil
	}

	d.StreamListItem(ctx, &accountData{
		commonColumnData: *commonColumnData,
	})

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOrganizationDetails(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getOrganizationDetails")

	// Create Session
	svc, err := OrganizationService(ctx, d)
	if err != nil {
		return nil, err
	}

	op, err := svc.DescribeOrganization(&organizations.DescribeOrganizationInput{})
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "AWSOrganizationsNotInUseException" {
				return nil, nil
			}
		}
		return nil, err
	}

	return op, nil
}

//// Transform Functions

func accountDataToTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsAccountAkas")
	accountInfo := d.HydrateItem.(*accountData)

	if accountInfo.Aliases != nil && len(accountInfo.Aliases) > 0 {
		return accountInfo.Aliases[0], nil
	}

	return accountInfo.commonColumnData.AccountId, nil
}

func accountARN(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("accountARN")
	accountInfo := d.HydrateItem.(*accountData)

	arn := "arn:" + accountInfo.commonColumnData.Partition + ":::" + accountInfo.commonColumnData.AccountId

	return arn, nil
}
