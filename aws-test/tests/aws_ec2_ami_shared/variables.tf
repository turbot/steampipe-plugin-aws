variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "trusted_accounts_allow" {
  type    = string
  default = "388460667113"
}

variable "trusted_accounts_deny" {
  type    = string
  default = "013122550996"
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

# Create AWS > EBS > Volume
resource "aws_ebs_volume" "named_test_resource" {
  availability_zone = "us-east-1a"
  size              = 8
}

# Create AWS > EBS > Snapshot
resource "aws_ebs_snapshot" "named_test_resource" {
  volume_id = aws_ebs_volume.named_test_resource.id
}

# Create AWS > EC2 > AMI
resource "aws_ami" "named_test_resource" {
  name                = var.resource_name
  description         = var.resource_name
  virtualization_type = "hvm"
  root_device_name    = "/dev/sda1"
  ebs_block_device {
    device_name = "/dev/sda1"
    snapshot_id = aws_ebs_snapshot.named_test_resource.id
    volume_size = 8
  }
}

resource "aws_ami_launch_permission" "allow_access" {
  image_id   = aws_ami.named_test_resource.id
  account_id = var.trusted_accounts_allow
}

resource "aws_ami_launch_permission" "deny_access" {
  image_id   = aws_ami.named_test_resource.id
  account_id = var.trusted_accounts_deny
}

data "template_file" "resource_aka" {
  template = "arn:$${partition}:ec2:$${region}:$${account_id}:image/${aws_ami.named_test_resource.id}"
  vars = {
    partition  = data.aws_partition.current.partition
    account_id = data.aws_caller_identity.current.account_id
    region     = data.aws_region.primary.name
  }
}

output "resource_aka" {
  depends_on = [aws_ami.named_test_resource]
  value      = data.template_file.resource_aka.rendered
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "region_name" {
  value = data.aws_region.primary.name
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = aws_ami.named_test_resource.id
}

output "snapshot_id" {
  value = aws_ebs_snapshot.named_test_resource.id
}
