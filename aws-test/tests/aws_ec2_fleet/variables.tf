variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update-fleet"
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

data "aws_ami" "ubuntu" {
  most_recent = true
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  owners = ["099720109477"]
}

resource "aws_launch_template" "named_test_resource" {
  name = var.resource_name
  image_id = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"

  tags = {
    Name = var.resource_name
  }
}

resource "aws_ec2_fleet" "named_test_resource" {
  launch_template_config {
    launch_template_specification {
      launch_template_id = aws_launch_template.named_test_resource.id
      version            = aws_launch_template.named_test_resource.latest_version
    }
  }

  target_capacity_specification {
    default_target_capacity_type = "spot"
    total_target_capacity        = 1
  }

  type = "maintain"

  excess_capacity_termination_policy = "termination"

  replace_unhealthy_instances = false

  tags = {
      Name = var.resource_name
  }
}

output "resource_aka" {
  value = aws_ec2_fleet.named_test_resource.arn
}

output "resource_id" {
  value = aws_ec2_fleet.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "launch_template_id" {
  value = aws_launch_template.named_test_resource.id
}
