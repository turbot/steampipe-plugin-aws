
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

data "null_data_source" "resource" {
  inputs = {
    scope = "arn:${data.aws_partition.current.partition}:::${data.aws_caller_identity.current.account_id}"
  }
}

resource "tls_private_key" "example" {
  algorithm = "RSA"
}

resource "tls_self_signed_cert" "example" {
  private_key_pem = tls_private_key.example.private_key_pem

  subject {
    common_name  = "${var.resource_name}.com"
    organization = "Turbot HQ Pvt. Ltd."
  }

  validity_period_hours = 12

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

resource "aws_acm_certificate" "example" {
  private_key      = tls_private_key.example.private_key_pem
  certificate_body = tls_self_signed_cert.example.cert_pem
}

resource "aws_apigatewayv2_domain_name" "named_test_resource" {
  domain_name = "${var.resource_name}.com"

  tags = {
    name = var.resource_name
  }

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.example.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

output "resource_aka" {
  value = aws_apigatewayv2_domain_name.named_test_resource.arn
}

output "resource_id" {
  value = aws_apigatewayv2_domain_name.named_test_resource.id
}

output "api_mapping_selection_expression" {
  value = aws_apigatewayv2_domain_name.named_test_resource.api_mapping_selection_expression
}

output "resource_name" {
  value = "${var.resource_name}.com"
}

output "tags" {
  value = var.resource_name
}
