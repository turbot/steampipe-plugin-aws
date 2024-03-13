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

data "null_data_source" "resource" {
  inputs = {
    scope = "arn:${data.aws_partition.current.partition}:::${data.aws_caller_identity.current.account_id}"
  }
}

resource "aws_ebs_volume" "named_test_resource" {
  availability_zone = "${var.aws_region}a"
  size              = 1
  tags = {
    Name = var.resource_name
  }
}

data "aws_resourcegroupstaggingapi_resources" "named_test_resource" {
  depends_on = [aws_ebs_volume.named_test_resource]
  tag_filter {
    key    = "Name"
    values = [var.resource_name]
  }
}

output "resource_aka" {
  value = data.aws_resourcegroupstaggingapi_resources.named_test_resource.resource_tag_mapping_list[0].resource_arn
}

output "resource_id" {
  value = aws_ebs_volume.named_test_resource.id
}

output "resource_tags" {
  value = data.aws_resourcegroupstaggingapi_resources.named_test_resource.resource_tag_mapping_list[0].tags
}

output "resource_name" {
  value = var.resource_name
}

output "resource_title" {
  value = "volume/${aws_ebs_volume.named_test_resource.id}"
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "aws_region" {
  value = data.aws_region.primary.id
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}
