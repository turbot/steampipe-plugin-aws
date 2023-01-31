
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
  default     = "us-east-2"
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

resource "tls_private_key" "example" {
  algorithm = "RSA"
}

resource "tls_self_signed_cert" "example" {
  # key_algorithm   = "RSA"
  private_key_pem = tls_private_key.example.private_key_pem

  subject {
    common_name  = "turbot.com"
    organization = "Turbot HQ Pvt. Ltd."
  }

  validity_period_hours = 12

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

resource "aws_iam_server_certificate" "named_test_resource" {
  name             = var.resource_name
  certificate_body = tls_self_signed_cert.example.cert_pem
  private_key      = tls_private_key.example.private_key_pem
}

output "certificate_body" {
  value = replace(tls_self_signed_cert.example.cert_pem, "\n", "\\n")
}

output "resource_aka" {
  value = aws_iam_server_certificate.named_test_resource.arn
}

output "resource_id" {
  value = aws_iam_server_certificate.named_test_resource.id
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "resource_name" {
  value = var.resource_name
}
