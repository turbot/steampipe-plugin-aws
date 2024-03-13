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

resource "aws_cloudfront_response_headers_policy" "named_test_resource" {
  name    = var.resource_name
  comment = "test comment"

  security_headers_config {
    content_type_options {
      override = true
    }
    frame_options {
      frame_option = "SAMEORIGIN"
      override     = false
    }
    referrer_policy {
      override        = false
      referrer_policy = "strict-origin-when-cross-origin"
    }

    strict_transport_security {
      access_control_max_age_sec = 31536000
      override                   = false
    }
    xss_protection {
      mode_block = true
      override   = false
      protection = true
    }
  }

  custom_headers_config {
    items {
      header   = "X-Permitted-Cross-Domain-Policies"
      override = true
      value    = "none"
    }

    items {
      header   = "X-Test"
      override = true
      value    = "none"
    }
  }

  cors_config {
    access_control_allow_credentials = true

    access_control_allow_headers {
      items = ["test"]
    }

    access_control_allow_methods {
      items = ["GET"]
    }

    access_control_allow_origins {
      items = ["test.example.comtest"]
    }

    origin_override = true
  }
}

locals {
  resource_aka = "arn:aws:cloudfront::${data.aws_caller_identity.current.account_id}:response-headers-policy/${aws_cloudfront_response_headers_policy.named_test_resource.id}"
}

output "resource_aka" {
  value = local.resource_aka
}

output "resource_id" {
  value = aws_cloudfront_response_headers_policy.named_test_resource.id
}

output "etag" {
  value = aws_cloudfront_response_headers_policy.named_test_resource.etag
}

output "resource_name" {
  value = var.resource_name
}

output "region_name" {
  value = data.aws_region.primary.name
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}
