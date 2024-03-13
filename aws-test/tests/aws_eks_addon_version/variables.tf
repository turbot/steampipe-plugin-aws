
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

resource "aws_eks_addon" "named_test_resource" {
  cluster_name = aws_eks_cluster.named_test_resource.name
  addon_name   = "vpc-cni"
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
  vpc_id            = aws_vpc.named_test_resource.id
  cidr_block        = "172.31.0.0/20"
  availability_zone = "${var.aws_region}b"
}

resource "aws_subnet" "named_test_resource2" {
  vpc_id            = aws_vpc.named_test_resource.id
  cidr_block        = "172.31.32.0/20"
  availability_zone = "${var.aws_region}d"
}

locals {
  path = "${path.cwd}/output.json"
}

resource "null_resource" "named_test_resource" {
  depends_on = [aws_eks_addon.named_test_resource]
  provisioner "local-exec" {
    command = "aws eks describe-addon-versions --addon-name vpc-cni > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

data "template_file" "resource_aka" {
  depends_on = [null_resource.named_test_resource]
  template   = "arn:$${partition}:eks:$${region}:$${account_id}:addonversion/$${addon_name}/$${version_name}"
  vars = {
    version_name     = jsondecode(data.local_file.input.content).addons[0].addonVersions[0].addonVersion
    addon_name       = jsondecode(data.local_file.input.content).addons[0].addonName
    partition        = data.aws_partition.current.partition
    account_id       = data.aws_caller_identity.current.account_id
    region           = data.aws_region.primary.name
    alternate_region = data.aws_region.alternate.name
  }
}

output "resource_aka" {
  depends_on = [null_resource.named_test_resource]
  value      = data.template_file.resource_aka.rendered
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "region_id" {
  value = var.aws_region
}

output "addon_version" {
  value = jsondecode(data.local_file.input.content).addons[0].addonVersions[0].addonVersion
}

output "addon_name" {
  value = jsondecode(data.local_file.input.content).addons[0].addonName
}
