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

resource "aws_efs_file_system" "my_fileSystem" {
  creation_token = "var.resource_name"
}

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "my_subnet" {
  vpc_id            = aws_vpc.my_vpc.id
  availability_zone = "${var.aws_region}a"
  cidr_block        = "10.0.1.0/24"
}

resource "aws_efs_mount_target" "named_test_resource" {
  file_system_id = aws_efs_file_system.my_fileSystem.id
  subnet_id      = aws_subnet.my_subnet.id
}

data "template_file" "resource_aka" {
  template = "arn:$${partition}:elasticfilesystem:$${region}:$${account_id}:file-system/${aws_efs_file_system.my_fileSystem.id}/mount-target/${aws_efs_mount_target.named_test_resource.id}"
  vars = {
    resource_name    = var.resource_name
    partition        = data.aws_partition.current.partition
    account_id       = data.aws_caller_identity.current.account_id
    region           = data.aws_region.primary.name
    alternate_region = data.aws_region.alternate.name
  }
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "resource_aka" {
  depends_on = [aws_efs_mount_target.named_test_resource]
  value      = data.template_file.resource_aka.rendered
}

output "file_system_id" {
  value = aws_efs_file_system.my_fileSystem.id
}

output "mount_target_id" {
  value = aws_efs_mount_target.named_test_resource.id
}
