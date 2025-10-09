
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

# Create AWS > EKS > Cluster
resource "aws_eks_cluster" "named_test_resource" {
  name     = var.resource_name
  role_arn = aws_iam_role.cluster_role.arn

  vpc_config {
    subnet_ids = [aws_subnet.named_test_resource1.id, aws_subnet.named_test_resource2.id]
  }

  depends_on = [
    aws_iam_role_policy_attachment.cluster_policy,
    aws_iam_role_policy_attachment.cluster_vpc_policy,
  ]

  tags = {
    Name = var.resource_name
  }
}

# Create IAM role for EKS cluster
resource "aws_iam_role" "cluster_role" {
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

resource "aws_iam_role_policy_attachment" "cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.cluster_role.name
}

resource "aws_iam_role_policy_attachment" "cluster_vpc_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
  role       = aws_iam_role.cluster_role.name
}

# Create IAM role for access entry
resource "aws_iam_role" "access_entry_role" {
  name = "${var.resource_name}-access-entry"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY

  tags = {
    Name = "${var.resource_name}-access-entry"
  }
}

# Create VPC and subnets
resource "aws_vpc" "named_test_resource" {
  cidr_block = "172.31.0.0/16"
}

resource "aws_subnet" "named_test_resource1" {
  vpc_id            = aws_vpc.named_test_resource.id
  cidr_block        = "172.31.0.0/20"
  availability_zone = "${var.aws_region}b"
}

resource "aws_subnet" "named_test_resource2" {
  vpc_id            = aws_vpc.named_test_resource.id
  cidr_block        = "172.31.32.0/20"
  availability_zone = "${var.aws_region}d"
}

# Create EKS Access Entry
resource "aws_eks_access_entry" "named_test_resource" {
  cluster_name  = aws_eks_cluster.named_test_resource.name
  principal_arn = aws_iam_role.access_entry_role.arn
  type          = "STANDARD"
}

# Create EKS Access Policy Association
resource "aws_eks_access_policy_association" "named_test_resource" {
  cluster_name  = aws_eks_cluster.named_test_resource.name
  principal_arn = aws_iam_role.access_entry_role.arn
  policy_arn    = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"

  access_scope {
    type = "cluster"
  }

  depends_on = [aws_eks_access_entry.named_test_resource]
}

output "cluster_name" {
  value = aws_eks_cluster.named_test_resource.name
}

output "principal_arn" {
  value = aws_iam_role.access_entry_role.arn
}

output "policy_arn" {
  value = aws_eks_access_policy_association.named_test_resource.policy_arn
}

output "access_scope_type" {
  value = aws_eks_access_policy_association.named_test_resource.access_scope[0].type
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "resource_name" {
  value = var.resource_name
}


