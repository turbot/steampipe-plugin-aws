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

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_security_group" "security_group" {
  description = "Allow TLS inbound traffic"
  vpc_id      = aws_vpc.my_vpc.id
  ingress {
    from_port = 443
    to_port   = 443
    protocol  = "tcp"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.my_vpc.id
}

resource "aws_subnet" "my_subnet" {
  vpc_id            = aws_vpc.my_vpc.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = "${var.aws_region}a"
  depends_on        = [aws_internet_gateway.igw]
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

resource "aws_instance" "instance" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t2.micro"
  subnet_id                   = aws_subnet.my_subnet.id
  associate_public_ip_address = false
  tags = {
    name = var.resource_name
  }
}

resource "aws_ssm_document" "test_resource" {
  name          = var.resource_name
  document_type = "Command"
  tags = {
    name = var.resource_name
  }
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {
    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "named_test_resource" {
  name             = aws_ssm_document.test_resource.id
  association_name = var.resource_name
  targets {
    key = "InstanceIds"
    values = [
      aws_instance.instance.id
    ]
  }
}

output "instance_id" {
  value = aws_instance.instance.id
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
  value = "arn:${data.aws_partition.current.partition}:ssm:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:association/${aws_ssm_association.named_test_resource.association_id}"
}

output "resource_id" {
  value = aws_ssm_association.named_test_resource.association_id
}

output "resource_name" {
  value = var.resource_name
}
