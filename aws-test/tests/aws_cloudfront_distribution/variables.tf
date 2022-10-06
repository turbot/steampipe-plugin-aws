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

resource "aws_s3_bucket" "test_bucket" {
  bucket        = var.resource_name
  force_destroy = true
}

resource "aws_cloudfront_distribution" "named_test_resource" {
  origin {
    domain_name = aws_s3_bucket.test_bucket.bucket_regional_domain_name
    origin_id   = "s3-${aws_s3_bucket.test_bucket.bucket}"
  }
  enabled             = false
  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "s3-${aws_s3_bucket.test_bucket.bucket}"
    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  ordered_cache_behavior {
    path_pattern     = "/content/*"
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "s3-${aws_s3_bucket.test_bucket.bucket}"
    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
    compress               = true
    viewer_protocol_policy = "redirect-to-https"
  }
  price_class = "PriceClass_200"
  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["US", "CA", "GB", "DE"]
    }
  }
  tags = {
    name = var.resource_name
  }
  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

output "resource_aka" {
  value = aws_cloudfront_distribution.named_test_resource.arn
}

output "domain_name" {
  value = aws_cloudfront_distribution.named_test_resource.domain_name
}

output "resource_id" {
  value = aws_cloudfront_distribution.named_test_resource.id
}

output "is_ipv6_enabled" {
  value = aws_cloudfront_distribution.named_test_resource.is_ipv6_enabled
}

output "etag" {
  value = aws_cloudfront_distribution.named_test_resource.etag
}

output "resource_name" {
  value = var.resource_name
}
