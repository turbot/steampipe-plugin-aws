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

resource "aws_licensemanager_license_configuration" "named_test_resource" {
  name                     = var.resource_name
  description              = "Testing"
  license_count            = 10
  license_count_hard_limit = true
  license_counting_type    = "Socket"

  license_rules = [
    "#minimumSockets=2",
  ]
}

resource "aws_iam_instance_profile" "instance_profile" {
  name = var.resource_name
  role = aws_iam_role.role.name
}

resource "aws_iam_role" "role" {
  name = var.resource_name
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
               "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

resource "aws_launch_template" "named_test_resource" {
  name = var.resource_name

  block_device_mappings {
    device_name = "/dev/sda1"

    ebs {
      volume_size = 20
    }
  }

  capacity_reservation_specification {
    capacity_reservation_preference = "open"
  }

  cpu_options {
    core_count       = 4
    threads_per_core = 2
  }

  iam_instance_profile {
    name = aws_iam_instance_profile.instance_profile.name
  }

  credit_specification {
    cpu_credits = "standard"
  }

  disable_api_stop        = true
  disable_api_termination = true

  ebs_optimized = true

  elastic_gpu_specifications {
    type = "test"
  }

  elastic_inference_accelerator {
    type = "eia1.medium"
  }

  instance_initiated_shutdown_behavior = "terminate"

  instance_market_options {
    market_type = "spot"
  }

  instance_type = "t2.micro"

  license_specification {
    license_configuration_arn = aws_licensemanager_license_configuration.named_test_resource.arn
  }

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
    instance_metadata_tags      = "enabled"
  }

  monitoring {
    enabled = true
  }

  network_interfaces {
    associate_public_ip_address = true
  }

  placement {
    availability_zone = "us-west-2a"
  }

  tag_specifications {
    resource_type = "instance"

    tags = {
      Name = var.resource_name
    }
  }

  tags = {
    Name = var.resource_name
  }

}

output "resource_aka" {
  value = aws_launch_template.named_test_resource.arn
}

output "resource_id" {
  value = aws_launch_template.named_test_resource.id
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
