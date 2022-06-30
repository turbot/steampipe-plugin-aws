
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
  default     = "ap-northeast-3"
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

resource "aws_default_vpc" "default" {}

resource "aws_default_subnet" "named_test_resource" {
  depends_on = [ aws_default_vpc.default]
  availability_zone = "${var.aws_region}c"
}

data "aws_ami" "linux" {
  most_recent = true
  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-2.0.20210427.0-x86_64-gp2"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  owners = ["137112412989"]
}

resource "aws_iam_role_policy_attachment" "test-attach" {
  role       = aws_iam_role.named_test_resource.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role" "named_test_resource" {
  name               = var.resource_name
  assume_role_policy = "{\n  \"Version\": \"2012-10-17\",\n  \"Statement\": [\n    {\n      \"Action\": \"sts:AssumeRole\",\n      \"Principal\": {\n        \"Service\": \"ec2.amazonaws.com\"\n      },\n      \"Effect\": \"Allow\",\n      \"Sid\": \"test\"\n    }\n  ]\n}\n"
  description        = "Test Role"
  tags = {
    name = var.resource_name
  }
}

resource "aws_iam_instance_profile" "named_test_resource" {
  depends_on = [
    aws_iam_role_policy_attachment.test-attach
  ]
  name = var.resource_name
  role = aws_iam_role.named_test_resource.name
}

resource "aws_instance" "named_test_resource" {
  ami                         = data.aws_ami.linux.id
  instance_type               = "t3.micro"
  subnet_id                   = aws_default_subnet.named_test_resource.id
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.named_test_resource.name
  tags = {
    Name = var.resource_name
  }
}

resource "null_resource" "delay" {
  depends_on = [
    aws_instance.named_test_resource
  ]
  provisioner "local-exec" {
    command = "sleep 180"
  }
}

output "arn" {
  value = "arn:${data.aws_partition.current.partition}:ssm:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:managed-instance/${aws_instance.named_test_resource.id}"
}

output "resource_id" {
  value = aws_instance.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}

output "instance_private_dns" {
  value = aws_instance.named_test_resource.private_dns
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
