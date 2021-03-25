select arn, name, log_group_name
from aws.aws_cloudwatch_log_stream
where name = '{{ resourceName }}'
