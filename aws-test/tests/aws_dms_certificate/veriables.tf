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

# Create a new certificate
resource "aws_dms_certificate" "named_test_resource" {
  certificate_id  = var.resource_name
  certificate_wallet = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tDQpNSUlDbmpDQ0FZWUNDUURQOUtxc0s5Unk3akFOQmdrcWhraUc5dzBCQVFzRkFEQVJNUTh3RFFZRFZRUUREQVoxDQpiblZ6WldRd0hoY05Nak13TlRFM01EZ3pOekE1V2hjTk1qUXdOVEUyTURnek56QTVXakFSTVE4d0RRWURWUVFEDQpEQVoxYm5WelpXUXdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFDMHNPeldDYTgwDQpGVDFLamg0U3dpQ1FoaTJza25SbHhQb1FrWGMrcVh3WjhlWUw1Z2FsWGxxRUsyRHNJZUVDYXpaR2kreXFFVUprDQpEd3VsUE1hOTBTZEVLNDdkcGJDMGRENEdjYUI2M3BkeXpSY1BaTnQzZldCNloxdlVaRUxUNE1GaEpvQ04rNGZBDQpYSVpQV295OUJxQ0JtUmpzcXdrOGRVWE0weXlKVnVIYUZnTjRQUHh2QzJVWW4yZ3FMMDB1RCtiYWFmWVArcmNoDQpRSldZNHFUVGFyUTRoend0Ym5qZTdTK0xudGJjbU5RUVpFQmRCS1duYzkrc0FTemZCZmtTbE1wOVBOR3BTTkxuDQpUYU5UMFdNUUdaVHkrd1E1dkJMOUo3WWZxNTJMQ1VyWHBkVmNGd2tuOVMwVVQxMEwwRGxaSHljdDVSUCtoL1JKDQpTTk80RVI1a1FpN3pBZ01CQUFFd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFGMjBRckR0TnNYNHZQNmxlWmxPDQo1ZkRyQ25sR0VxK1hhY0ZreUxEV3VkbDF2L3dHT0dvVWI5V1c5VFBubXJWSkwyVVVtZFdnRGs1VGhvbEV0MHYwDQo0ZDNSWHFCR0xlQUN3UlN1NllXVzczNUlFTmdWT21lOFJDYVN1OEltL1I3NU5GZU9ib3kydWo5bTdPZGVQSlI2DQpQWkJIeVRKWEdiNmJBUFJyeEhodERrTStYZVIwcGJKV2ZDME5ER25RL1dsMmx6UURtbTdDaUNlWk9SbE5uR3l3DQppN2k4V0FkZTlXZVBGNmxxaEdJdUlaaDhERHdxY2htanNqSGxBOC9XZUxudTJpVXRsd2FKbHBOelo2ai8vSUEzDQo1YVkxMVNHeFlFbVFQekdTd1NGbkJYQmFKZ2RrY0Z1dW04d3d1YWlHT1pNMHdneE4zcmR2eXNGSlZva29sdmwwDQpHdk09DQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t"

  tags = {
    Name = "test"
  }

}

output "resource_aka" {
  value = aws_dms_certificate.named_test_resource.certificate_arn
}

output "resource_name" {
  value = var.resource_name
}
