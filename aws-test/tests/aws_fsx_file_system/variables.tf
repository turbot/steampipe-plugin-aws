variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125"
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

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
}
resource "aws_subnet" "my_subnet" {
  vpc_id     = aws_vpc.my_vpc.id
  cidr_block = "10.0.1.0/24"
}

resource "aws_fsx_lustre_file_system" "named_test_resource" {
  storage_capacity = 1200
  subnet_ids = [
    "${aws_subnet.my_subnet.id}"
  ]
  tags = {
    foo = var.resource_name
  }
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  depends_on = [aws_fsx_lustre_file_system.named_test_resource]
  value      = "arn:${data.aws_partition.current.partition}:fsx:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:file-system/${aws_fsx_lustre_file_system.named_test_resource.id}"
}

output "resource_id" {
  value = aws_fsx_lustre_file_system.named_test_resource.id
}
