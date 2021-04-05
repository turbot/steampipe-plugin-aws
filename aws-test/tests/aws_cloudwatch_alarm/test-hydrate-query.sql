select name, alarm_arn, metric_name, comparison_operator, alarm_description, threshold, tags, akas
from aws.aws_cloudwatch_alarm
where name = '{{ resourceName }}';