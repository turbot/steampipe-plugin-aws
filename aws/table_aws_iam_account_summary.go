package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type awsIamAccountSummary struct {
	AccessKeysPerUserQuota            int32
	AccountAccessKeysPresent          int32
	AccountMFAEnabled                 bool
	AccountPasswordPresent            int32
	AccountSigningCertificatesPresent int32
	AssumeRolePolicySizeQuota         int32
	AttachedPoliciesPerGroupQuota     int32
	AttachedPoliciesPerRoleQuota      int32
	AttachedPoliciesPerUserQuota      int32
	GlobalEndpointTokenVersion        int32
	GroupPolicySizeQuota              int32
	Groups                            int32
	GroupsPerUserQuota                int32
	GroupsQuota                       int32
	InstanceProfiles                  int32
	InstanceProfilesQuota             int32
	MFADevices                        int32
	MFADevicesInUse                   int32
	Policies                          int32
	PoliciesQuota                     int32
	PolicySizeQuota                   int32
	PolicyVersionsInUse               int32
	PolicyVersionsInUseQuota          int32
	Providers                         int32
	RolePolicySizeQuota               int32
	Roles                             int32
	RolesQuota                        int32
	ServerCertificates                int32
	ServerCertificatesQuota           int32
	SigningCertificatesPerUserQuota   int32
	UserPolicySizeQuota               int32
	Users                             int32
	UsersQuota                        int32
	VersionsPerPolicyQuota            int32
}

func tableAwsIamAccountSummary(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "aws_iam_account_summary",
		Description:      "AWS IAM Account Summary",
		DefaultTransform: transform.FromGo(),
		List: &plugin.ListConfig{
			Hydrate: listAccountSummary,
			Tags:    map[string]string{"service": "iam", "action": "GetAccountSummary"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "access_keys_per_user_quota",
				Description: "Specifies the allowed quota of access keys per user.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "account_access_keys_present",
				Description: "Specifies the number of account level access keys present.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "account_mfa_enabled",
				Description: "Specifies whether MFA is enabled for the account.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AccountMFAEnabled"),
			},
			{
				Name:        "account_password_present",
				Description: "Specifies the number of account passwords present.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "account_signing_certificates_present",
				Description: "Specifies the number of account signing certificates present.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "assume_role_policy_size_quota",
				Description: "Specifies the allowed assume role policy size.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "attached_policies_per_group_quota",
				Description: "Specifies the allowed attached policies per group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "attached_policies_per_role_quota",
				Description: "Specifies the allowed attached policies per role.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "attached_policies_per_user_quota",
				Description: "Specifies the allowed attached policies per user.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "global_endpoint_token_version",
				Description: "Specifies the token version of the global endpoint.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "group_policy_size_quota",
				Description: "Specifies the allowed group policy size.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "groups",
				Description: "Specifies the number of groups.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "groups_per_user_quota",
				Description: "Specifies the allowed number of groups.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "groups_quota",
				Description: "Specifies the allowed number of groups.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "instance_profiles",
				Description: "Specifies the number of groups.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "instance_profiles_quota",
				Description: "Specifies the allowed number of groups.",
				Type:        proto.ColumnType_INT,
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
			},
			{
				Name:        "policies_quota",
				Description: "Specifies the allowed number of policies.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "policy_size_quota",
				Description: "Specifies the allowed size of policies.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "policy_versions_in_use",
				Description: "Specifies the number of policy versions in use.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "policy_versions_in_use_quota",
				Description: "Specifies the allowed number of policy versions.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "providers",
				Description: "Specifies the number of providers.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "role_policy_size_quota",
				Description: "Specifies the allowed role policy size.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "roles",
				Description: "Specifies the number of roles.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "roles_quota",
				Description: "Specifies the allowed number of roles.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "server_certificates",
				Description: "Specifies the number of server certificates.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "server_certificates_quota",
				Description: "Specifies the allowed number of server certificates.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "signing_certificates_per_user_quota",
				Description: "Specifies the allowed number of signing certificates per user.",
				Type:        proto.ColumnType_INT,
			}, {
				Name:        "user_policy_size_quota",
				Description: "Specifies the allowed user policy size.",
				Type:        proto.ColumnType_INT,
			}, {
				Name:        "users",
				Description: "Specifies the number of users.",
				Type:        proto.ColumnType_INT,
			}, {
				Name:        "users_quota",
				Description: "Specifies the allowed number of users.",
				Type:        proto.ColumnType_INT,
			}, {
				Name:        "versions_per_policy_quota",
				Description: "Specifies the allowed number of versions per policy.",
				Type:        proto.ColumnType_INT,
			},
		}),
	}
}

//// LIST FUNCTION

func listAccountSummary(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_account_summary.listAccountSummary", "client_error", err)
		return nil, err
	}

	resp, err := svc.GetAccountSummary(ctx, &iam.GetAccountSummaryInput{})
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_account_summary.listAccountSummary", "api_error", err)
		return nil, err
	}

	summaryMap := resp.SummaryMap
	accountSummary := &awsIamAccountSummary{
		AccessKeysPerUserQuota:            summaryMap["AccessKeysPerUserQuota"],
		AccountAccessKeysPresent:          summaryMap["AccountAccessKeysPresent"],
		AccountMFAEnabled:                 summaryMap["AccountMFAEnabled"] == int32(1),
		AccountPasswordPresent:            summaryMap["AccountPasswordPresent"],
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
