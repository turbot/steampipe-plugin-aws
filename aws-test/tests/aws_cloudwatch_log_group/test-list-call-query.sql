select arn, name, metric_filter_count, retention_in_days, stored_bytes
from aws.aws_cloudwatch_log_group
where arn = '{{ output.resource_aka.value }}:*'