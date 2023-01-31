
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

# Create AWS > EC2 > KeyPair
resource "aws_key_pair" "named_test_resource" {
  key_name   = var.resource_name
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC6B60dQX19dOO0wmNnPf27NJBXQhuVAsFcshY7YqTmNA1dfl0lXvJFSF0nVTvaxtE2854gmvygm/yafopG33zJeH+k2ZrblivkEuH0NNajKMp2mBumkk8RiTVvU1ET4m6Z1jZaK0dTSeV2k6gUk7Wmvo4+6RIu+wdAtwGGGcSq4HJ4M0G1CGHruBywMdOs5CklqRn7BRMh2aexowDachlabUKRFwTf7OCjDcoCWHo+o4kHV7NDEdhFzrE2fxXttt054kd5CMZOEOB9p01pAqSPrxFv3S/BpiXuhmjwsUsIdfHIgFrFhNY/M4+KhTa3pI69NqsigReotTEzk2QwHyHSjsRwfzeJNRgI8iqGIdYr3+t6zCKswyzgV0hqalZhXtpVUEBbumJAKWy+VzTnOXUQ5zEVvabTxi3dsW9zMGofjwhiC75x97RE7U/ZEanRM+0Avr8e3c3ob4NLf9JOwbvB/jYvS6j/nuwJ5H6igT3Oj0oDiknh2WtOgH/BKKaCq5Nyl15Df65i0PNgKRcLy+x0fubYUtsCOMPpQTfztSYGjbvAhAOrH2LJF0i2dTD3PcVw6R6uDpJSw9QkzdD+2ofxnqjGPUeKR3mraGOXt2hZe9F9bvvdZScDfdAoyBUYTO2/sxsQFobRtdiEk2up+mNIx9QXMXO2dOgtNUYZKc1DdQ== clara@turbot.com"
  tags = {
    name = var.resource_name
  }
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

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = aws_key_pair.named_test_resource.key_pair_id
}

output "key_fingerprint" {
  value = aws_key_pair.named_test_resource.fingerprint
}
