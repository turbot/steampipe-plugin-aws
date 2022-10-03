

variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "aws_profile" {
  type        = string
  default     = "default"
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
  path = "${path.cwd}/output.json"
}


resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "aws iam get-account-password-policy --output json --profile ${var.aws_profile} --region ${data.aws_region.primary.name} > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "expire_passwords" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "PasswordPolicy", "what_1"), "ExpirePasswords", "what_2")
}

output "allow_users_to_change_password" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "PasswordPolicy", "what_1"), "AllowUsersToChangePassword", "what_2")
}

output "minimum_password_length" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "PasswordPolicy", "what_1"), "MinimumPasswordLength", "what_2")
}

output "require_numbers" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "PasswordPolicy", "what_1"), "RequireNumbers", "what_2")
}

output "require_symbols" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "PasswordPolicy", "what_1"), "RequireSymbols", "what_2")
}

output "require_uppercase_characters" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "PasswordPolicy", "what_1"), "RequireUppercaseCharacters", "what_2")
}

output "require_lowercase_characters" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(lookup(jsondecode(data.local_file.input.content), "PasswordPolicy", "what_1"), "RequireLowercaseCharacters", "what_2")
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
