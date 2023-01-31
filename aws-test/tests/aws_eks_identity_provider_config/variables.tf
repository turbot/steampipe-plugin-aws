
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

resource "aws_eks_cluster" "named_test_resource" {
  name     = var.resource_name
  role_arn = aws_iam_role.named_test_resource.arn

  vpc_config {
    subnet_ids = [aws_subnet.named_test_resource1.id, aws_subnet.named_test_resource2.id]
  }

  # Ensure that IAM Role permissions are created before and deleted after EKS Cluster handling.
  # Otherwise, EKS will not be able to properly delete EKS managed EC2 infrastructure such as Security Groups.
  depends_on = [
    aws_iam_role_policy_attachment.named_test_resource,
    aws_iam_role_policy_attachment.named_test_resource2,
  ]
  tags = {
    Name = var.resource_name
  }
}


resource "aws_iam_role" "named_test_resource" {
  name = var.resource_name

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "named_test_resource" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.named_test_resource.name
}

# Optionally, enable Security Groups for Pods
# Reference: https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html
resource "aws_iam_role_policy_attachment" "named_test_resource2" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
  role       = aws_iam_role.named_test_resource.name
}

resource "aws_vpc" "named_test_resource" {
  cidr_block = "172.31.0.0/16"
}

resource "aws_subnet" "named_test_resource1" {
  vpc_id     = aws_vpc.named_test_resource.id
  cidr_block = "172.31.0.0/20"
  availability_zone = "${var.aws_region}b"
}

resource "aws_subnet" "named_test_resource2" {
  vpc_id     = aws_vpc.named_test_resource.id
  cidr_block = "172.31.32.0/20"
  availability_zone = "${var.aws_region}d"
}

resource "aws_eks_identity_provider_config" "named_test_resource" {
  cluster_name = aws_eks_cluster.named_test_resource.name

  oidc {
    client_id                     = aws_eks_cluster.named_test_resource.id
    identity_provider_config_name = var.resource_name
    issuer_url                    = "https://${aws_eks_cluster.named_test_resource.id}"
  }

  tags = {
    Name = var.resource_name
  }

  timeouts {
    create = "40m"
    delete = "40m"
  }
}

output "role_arn" {
  value = aws_iam_role.named_test_resource.arn
}

output "resource_id" {
  value = aws_eks_identity_provider_config.named_test_resource.id
}

output "resource_aka" {
  value = aws_eks_identity_provider_config.named_test_resource.arn
}

output "cluster_name" {
  value = aws_eks_cluster.named_test_resource.name
}

output "identity_provider_config_type" {
  value = "oidc"
}

output "status" {
  value = aws_eks_identity_provider_config.named_test_resource.status
}

output "identity" {
  value = aws_eks_cluster.named_test_resource.identity
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "resource_name" {
  value = var.resource_name
}