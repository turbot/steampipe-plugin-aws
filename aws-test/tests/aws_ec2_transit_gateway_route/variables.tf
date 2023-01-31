
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

# Create AWS > EC2 > Transit Gateway
resource "aws_ec2_transit_gateway" "named_test_resource" {
  description = "Test transit gateway"
}

# Create AWS > EC2 > Transit Gateway Route Table
resource "aws_ec2_transit_gateway_route_table" "named_test_resource" {
  transit_gateway_id = aws_ec2_transit_gateway.named_test_resource.id
  tags = {
    Name = var.resource_name
  }
}

resource "aws_ec2_transit_gateway_route" "named_test_resource" {
  destination_cidr_block         = "0.0.0.0/0"
  blackhole                      = true
  transit_gateway_route_table_id = aws_ec2_transit_gateway.named_test_resource.association_default_route_table_id
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  value = "arn:${data.aws_partition.current.partition}:ec2:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:transit-gateway-route-table/${replace(aws_ec2_transit_gateway_route.named_test_resource.id, "_", ":")}"
}

output "destination_cidr_block" {
  value = split("_", aws_ec2_transit_gateway_route.named_test_resource.id)[1]
}

output "transit_gateway_rtb_id" {
  value = split("_", aws_ec2_transit_gateway_route.named_test_resource.id)[0]
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
