

variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "aws_profile" {
  type        = string
  default     = "integration-tests"
  description = "AWS credentials profile used for the test. Default is to use the default profile."
}

variable "aws_region" {
  type        = string
  default     = "us-east-1"
  description = "AWS region used for the test. Does not work with default region in config, so must be defined here."
}

variable "aws_region_alternate" {
  type        = string
  default     = "us-east-2"
  description = "Alternate AWS region used for tests that require two regions (e.g. DynamoDB global tables)."
}

provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
}

provider "aws" {
  alias   = "alternate"
  profile = var.aws_profile
  region  = var.aws_region_alternate
}

data "aws_partition" "current" {}
data "aws_caller_identity" "current" {}
data "aws_region" "primary" {}
data "aws_region" "alternate" {
  provider = aws.alternate
}

locals {
  path = "${path.cwd}/account_summary.json"
}


resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "aws iam get-account-summary --output json --profile ${var.aws_profile} --region ${data.aws_region.primary.name} > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "access_keys_per_user_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "AccessKeysPerUserQuota", "AccessKeysPerUserQuota")
}

output "account_access_keys_present" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "AccountAccessKeysPresent", "AccountAccessKeysPresent")
}

output "account_mfa_enabled" {
  depends_on = [null_resource.named_test_resource]
  value      = (lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "AccountMFAEnabled", "AccountMFAEnabled")) == 1 ? true : false
}

output "account_signing_certificates_present" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "AccountSigningCertificatesPresent", "AccountSigningCertificatesPresent")
}

output "assume_role_policy_size_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "AssumeRolePolicySizeQuota", "AssumeRolePolicySizeQuota")
}

output "attached_policies_per_group_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "AttachedPoliciesPerGroupQuota", "AttachedPoliciesPerGroupQuota")
}

output "attached_policies_per_role_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "AttachedPoliciesPerRoleQuota", "AttachedPoliciesPerRoleQuota")
}
output "attached_policies_per_user_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "AttachedPoliciesPerUserQuota", "AttachedPoliciesPerUserQuota")
}
output "global_endpoint_token_version" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "GlobalEndpointTokenVersion", "GlobalEndpointTokenVersion")
}
output "group_policy_size_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "GroupPolicySizeQuota", "GroupPolicySizeQuota")
}
output "groups" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "Groups", "Groups")
}
output "groups_per_user_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "GroupsPerUserQuota", "GroupsPerUserQuota")
}
output "groups_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "GroupsQuota", "GroupsQuota")
}
output "instance_profiles" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "InstanceProfiles", "InstanceProfiles")
}
output "instance_profiles_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "InstanceProfilesQuota", "InstanceProfilesQuota")
}
output "mfa_devices" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "MFADevices", "MFADevices")
}
output "mfa_devices_in_use" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "MFADevicesInUse", "MFADevicesInUse")
}
output "policies" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "Policies", "Policies")
}
output "policies_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "PoliciesQuota", "PoliciesQuota")
}
output "policy_size_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "PolicySizeQuota", "PolicySizeQuota")
}
output "policy_versions_in_use" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "PolicyVersionsInUse", "PolicyVersionsInUse")
}
output "policy_versions_in_use_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "PolicyVersionsInUseQuota", "PolicyVersionsInUseQuota")
}
output "providers" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "Providers", "Providers")
}
output "role_policy_size_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "RolePolicySizeQuota", "RolePolicySizeQuota")
}
output "roles" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "Roles", "Roles")
}
output "roles_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "RolesQuota", "RolesQuota")
}
output "server_certificates" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "ServerCertificates", "ServerCertificates")
}
output "server_certificates_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "ServerCertificatesQuota", "ServerCertificatesQuota")
}
output "signing_certificates_per_user_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "SigningCertificatesPerUserQuota", "SigningCertificatesPerUserQuota")
}
output "user_policy_size_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "UserPolicySizeQuota", "UserPolicySizeQuota")
}
output "users" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "Users", "Users")
}
output "users_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "UsersQuota", "UsersQuota")
}
output "versions_per_policy_quota" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "SummaryMap", "SummaryMap"), "VersionsPerPolicyQuota", "VersionsPerPolicyQuota")
}

output "aws_account" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}
