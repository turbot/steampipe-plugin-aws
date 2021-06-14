select name, arn, metric_name, comparison_operator, alarm_description
from aws.aws_cloudwatch_alarm
where name = '{{ resourceName }}';