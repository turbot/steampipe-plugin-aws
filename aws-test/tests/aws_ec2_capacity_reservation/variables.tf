variable "resource_name" {
  type        = string
  default     = "turbot-test-20210820-create-update"
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

resource "aws_ec2_capacity_reservation" "named_test_resource" {
  instance_type     = "t2.micro"
  instance_platform = "Linux/UNIX"
  availability_zone = "us-east-2a"
  instance_count    = 1
}

output "resource_aka" {
  value = aws_ec2_capacity_reservation.named_test_resource.arn
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

output "resource_id" {
  value = aws_ec2_capacity_reservation.named_test_resource.id
}

output "instance_type" {
  value = aws_ec2_capacity_reservation.named_test_resource.instance_type
}

output "availability_zone" {
  value = aws_ec2_capacity_reservation.named_test_resource.availability_zone
}

output "instance_count" {
  value = aws_ec2_capacity_reservation.named_test_resource.instance_count
}

output "instance_platform" {
  value = aws_ec2_capacity_reservation.named_test_resource.instance_platform
}
