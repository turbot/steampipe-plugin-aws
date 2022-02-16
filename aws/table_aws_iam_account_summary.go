package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

type awsIamAccountSummary struct {
	AccessKeysPerUserQuota            *int64
	AccountAccessKeysPresent          *int64
	AccountMFAEnabled                 bool
	AccountSigningCertificatesPresent *int64
	AssumeRolePolicySizeQuota         *int64
	AttachedPoliciesPerGroupQuota     *int64
	AttachedPoliciesPerRoleQuota      *int64
	AttachedPoliciesPerUserQuota      *int64
	GlobalEndpointTokenVersion        *int64
	GroupPolicySizeQuota              *int64
	Groups                            *int64
	GroupsPerUserQuota                *int64
	GroupsQuota                       *int64
	InstanceProfiles                  *int64
	InstanceProfilesQuota             *int64
	MFADevices                        *int64
	MFADevicesInUse                   *int64
	Policies                          *int64
	PoliciesQuota                     *int64
	PolicySizeQuota                   *int64
	PolicyVersionsInUse               *int64
	PolicyVersionsInUseQuota          *int64
	Providers                         *int64
	RolePolicySizeQuota               *int64
	Roles                             *int64
	RolesQuota                        *int64
	ServerCertificates                *int64
	ServerCertificatesQuota           *int64
	SigningCertificatesPerUserQuota   *int64
	UserPolicySizeQuota               *int64
	Users                             *int64
	UsersQuota                        *int64
	VersionsPerPolicyQuota            *int64
}

func tableAwsIamAccountSummary(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_account_summary",
		Description: "AWS IAM Account Summary",
		List: &plugin.ListConfig{
			Hydrate: listAccountSummary,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "access_keys_per_user_quota",
				Description: "Specifies the allowed quota of access keys per user.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "account_access_keys_present",
				Description: "Specifies the number of account level access keys present.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "account_mfa_enabled",
				Description: "Specifies whether MFA is enabled for the account.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AccountMFAEnabled"),
			},
			{
				Name:        "account_signing_certificates_present",
				Description: "Specifies the number of account signing certificates present.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "assume_role_policy_size_quota",
				Description: "Specifies the allowed assume role policy size.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "attached_policies_per_group_quota",
				Description: "Specifies the allowed attached policies per group.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "attached_policies_per_role_quota",
				Description: "Specifies the allowed attached policies per role.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "attached_policies_per_user_quota",
				Description: "Specifies the allowed attached policies per user.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "global_endpoint_token_version",
				Description: "Specifies the token version of the global endpoint.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "group_policy_size_quota",
				Description: "Specifies the allowed group policy size.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "groups",
				Description: "Specifies the number of groups.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "groups_per_user_quota",
				Description: "Specifies the allowed number of groups.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "groups_quota",
				Description: "Specifies the allowed number of groups.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "instance_profiles",
				Description: "Specifies the number of groups.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "instance_profiles_quota",
				Description: "Specifies the allowed number of groups.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "mfa_devices",
				Description: "Specifies the number of MFA devices.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MFADevices"),
			},
			{
				Name:        "mfa_devices_in_use",
				Description: "Specifies the number of MFA devices in use.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MFADevicesInUse"),
			},
			{
				Name:        "policies",
				Description: "Specifies the number of policies.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "policies_quota",
				Description: "Specifies the allowed number of policies.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "policy_size_quota",
				Description: "Specifies the allowed size of policies.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "policy_versions_in_use",
				Description: "Specifies the number of policy versions in use.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "policy_versions_in_use_quota",
				Description: "Specifies the allowed number of policy versions.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "providers",
				Description: "Specifies the number of providers.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "role_policy_size_quota",
				Description: "Specifies the allowed role policy size.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "roles",
				Description: "Specifies the number of roles.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "roles_quota",
				Description: "Specifies the allowed number of roles.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "server_certificates",
				Description: "Specifies the number of server certificates.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "server_certificates_quota",
				Description: "Specifies the allowed number of server certificates.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "signing_certificates_per_user_quota",
				Description: "Specifies the allowed number of signing certificates per user.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			}, {
				Name:        "user_policy_size_quota",
				Description: "Specifies the allowed user policy size.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			}, {
				Name:        "users",
				Description: "Specifies the number of users.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			}, {
				Name:        "users_quota",
				Description: "Specifies the allowed number of users.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			}, {
				Name:        "versions_per_policy_quota",
				Description: "Specifies the allowed number of versions per policy.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAccountSummary(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAccountSummary")

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	resp, err := svc.GetAccountSummary(&iam.GetAccountSummaryInput{})
	if err != nil {
		return nil, err
	}

	summaryMap := resp.SummaryMap
	accountSummary := &awsIamAccountSummary{
		AccessKeysPerUserQuota:            summaryMap["AccessKeysPerUserQuota"],
		AccountAccessKeysPresent:          summaryMap["AccountAccessKeysPresent"],
		AccountMFAEnabled:                 *summaryMap["AccountMFAEnabled"] == int64(1),
		AccountSigningCertificatesPresent: summaryMap["AccountSigningCertificatesPresent"],
		AssumeRolePolicySizeQuota:         summaryMap["AssumeRolePolicySizeQuota"],
		AttachedPoliciesPerGroupQuota:     summaryMap["AttachedPoliciesPerGroupQuota"],
		AttachedPoliciesPerRoleQuota:      summaryMap["AttachedPoliciesPerRoleQuota"],
		AttachedPoliciesPerUserQuota:      summaryMap["AttachedPoliciesPerUserQuota"],
		GlobalEndpointTokenVersion:        summaryMap["GlobalEndpointTokenVersion"],
		GroupPolicySizeQuota:              summaryMap["GroupPolicySizeQuota"],
		Groups:                            summaryMap["Groups"],
		GroupsPerUserQuota:                summaryMap["GroupsPerUserQuota"],
		GroupsQuota:                       summaryMap["GroupsQuota"],
		InstanceProfiles:                  summaryMap["InstanceProfiles"],
		InstanceProfilesQuota:             summaryMap["InstanceProfilesQuota"],
		MFADevices:                        summaryMap["MFADevices"],
		MFADevicesInUse:                   summaryMap["MFADevicesInUse"],
		Policies:                          summaryMap["Policies"],
		PoliciesQuota:                     summaryMap["PoliciesQuota"],
		PolicySizeQuota:                   summaryMap["PolicySizeQuota"],
		PolicyVersionsInUse:               summaryMap["PolicyVersionsInUse"],
		PolicyVersionsInUseQuota:          summaryMap["PolicyVersionsInUseQuota"],
		Providers:                         summaryMap["Providers"],
		RolePolicySizeQuota:               summaryMap["RolePolicySizeQuota"],
		Roles:                             summaryMap["Roles"],
		RolesQuota:                        summaryMap["RolesQuota"],
		ServerCertificates:                summaryMap["ServerCertificates"],
		ServerCertificatesQuota:           summaryMap["ServerCertificatesQuota"],
		SigningCertificatesPerUserQuota:   summaryMap["SigningCertificatesPerUserQuota"],
		UserPolicySizeQuota:               summaryMap["UserPolicySizeQuota"],
		Users:                             summaryMap["Users"],
		UsersQuota:                        summaryMap["UsersQuota"],
		VersionsPerPolicyQuota:            summaryMap["VersionsPerPolicyQuota"],
	}

	d.StreamListItem(ctx, accountSummary)
	return nil, nil
}
