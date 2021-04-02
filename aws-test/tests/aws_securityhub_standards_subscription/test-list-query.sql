select standards_arn, name, enabled_by_default, standards_status, standards_subscription_arn
from aws.aws_securityhub_standards_subscription
where name = "CIS AWS Foundations Benchmark v1.2.0";