
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

# Create AWS > EC2 > Volume
resource "aws_ebs_volume" "my_volume" {
  availability_zone = "us-east-2a"
  size              = 8
  tags = {
    Name = "turbot-volume-test"
  }
}

# Create AWS > EC2 > Snapshot
resource "aws_ebs_snapshot" "my_snapshot" {
  volume_id = aws_ebs_volume.my_volume.id
  tags = {
    Name = "turbot-snapshot-test"
  }
}

# Create AWS > EC2 > AMI
resource "aws_ami" "named_test_resource" {
  name                = var.resource_name
  virtualization_type = "hvm"
  root_device_name    = "/dev/sda1"
  ebs_block_device {
    device_name = "/dev/sda1"
    snapshot_id = aws_ebs_snapshot.my_snapshot.id
    volume_size = 8
  }
  tags = {
    name = var.resource_name
  }
}

# Create AWS > EC2 > Launch Configuration
resource "aws_launch_configuration" "named_test_resource" {
  name          = var.resource_name
  image_id      = aws_ami.named_test_resource.id
  instance_type = "t2.micro"
}

# Create AWS > EC2 > AutoScaling Group
resource "aws_autoscaling_group" "named_test_resource" {
  name                 = var.resource_name
  max_size             = 0
  min_size             = 0
  availability_zones   = ["us-east-2a"]
  force_delete         = true
  launch_configuration = aws_launch_configuration.named_test_resource.name
  tag {
    key                 = "Name"
    value               = var.resource_name
    propagate_at_launch = true
  }
}

resource "aws_autoscaling_policy" "bat" {
  name                   = var.resource_name
  scaling_adjustment     = 4
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = aws_autoscaling_group.named_test_resource.name
}

output "resource_aka" {
  value = aws_autoscaling_group.named_test_resource.arn
}

output "policy_name" {
  value = aws_autoscaling_policy.bat.name
}

output "policy_arn" {
  value = aws_autoscaling_policy.bat.arn
}

output "resource_id" {
  value = aws_autoscaling_group.named_test_resource.id
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

output "region_name" {
  value = data.aws_region.primary.name
}
