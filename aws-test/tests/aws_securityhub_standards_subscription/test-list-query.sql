select standards_arn, name, enabled_by_default, standards_status, standards_subscription_arn
from aws.aws_securityhub_standards_subscription
where standards_arn = 'arn:aws:securityhub:::ruleset/cis-aws-foundations-benchmark/v/1.2.0' and region = '{{ output.aws_region.value }}';