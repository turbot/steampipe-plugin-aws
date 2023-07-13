
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

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    name = var.resource_name
  }
}

resource "aws_subnet" "named_test_resource" {
  vpc_id     = aws_vpc.main.id
  availability_zone = "us-east-1b"
  cidr_block = "10.0.1.0/24"
  tags = {
    Name = var.resource_name
  }
}

resource "aws_appstream_fleet" "named_test_resource" {
  name = var.resource_name

  compute_capacity {
    desired_instances = 1
  }

  description                        = "test fleet"
  # idle_disconnect_timeout_in_seconds = 60
  display_name                       = var.resource_name
  enable_default_internet_access     = false
  fleet_type                         = "ON_DEMAND"
  image_name                         = "AppStream-WinServer2019-06-12-2023"
  instance_type                      = "stream.standard.small"
  # max_user_duration_in_seconds       = 600

  vpc_config {
    subnet_ids = [aws_subnet.named_test_resource.id]
  }

  tags = {
    TagName = "tag-value"
  }
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "resource_aka" {
  value = aws_appstream_fleet.named_test_resource.arn
}

output "resource_state" {
  value = aws_appstream_fleet.named_test_resource.state
}

output "resource_id" {
  value = aws_appstream_fleet.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}
