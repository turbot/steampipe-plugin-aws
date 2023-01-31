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

resource "aws_kinesis_video_stream" "named_test_resource" {
  name                    = var.resource_name
  data_retention_in_hours = 24
  device_name             = "kinesis-video-device-name"
  media_type              = "video/h264"
  tags = {
    Name = var.resource_name,
    Foo  = "Bar"
  }
}

output "resource_aka" {
  value = aws_kinesis_video_stream.named_test_resource.arn
}

output "resource_id" {
  value = aws_kinesis_video_stream.named_test_resource.id
}

output "resource_creation_time" {
  value = aws_kinesis_video_stream.named_test_resource.creation_time
}

output "resource_version" {
  value = aws_kinesis_video_stream.named_test_resource.version
}

output "resource_name" {
  value = var.resource_name
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "aws_region" {
  value = var.aws_region
}
