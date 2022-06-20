
variable "resource_name" {
  type        = string
  default     = "turbot-test-create-update"
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

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.my_vpc.id
}

resource "aws_subnet" "my_subnet1" {
  vpc_id            = aws_vpc.my_vpc.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = "${var.aws_region}a"
  depends_on        = [aws_internet_gateway.igw]
}

resource "aws_subnet" "my_subnet2" {
  vpc_id            = aws_vpc.my_vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "${var.aws_region}b"
  depends_on        = [aws_internet_gateway.igw]
}

resource "aws_lb" "my_lb" {
  name               = var.resource_name
  internal           = false
  load_balancer_type = "application"
  subnets            = ["${aws_subnet.my_subnet1.id}", "${aws_subnet.my_subnet2.id}"]

  enable_deletion_protection = false
}

resource "aws_lb_target_group" "my_targetGroup" {
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.my_vpc.id
}

resource "aws_lb_listener" "named_test_resource" {
  load_balancer_arn = aws_lb.my_lb.arn
  port              = "443"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.my_targetGroup.arn
  }
}

output "resource_aka" {
  value = aws_lb_listener.named_test_resource.arn
}

output "resource_id" {
  value = aws_lb_listener.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}

output "load_balancer_arn" {
  value = aws_lb.my_lb.arn
}

output "target_group_arn" {
  value = aws_lb_target_group.my_targetGroup.arn
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
