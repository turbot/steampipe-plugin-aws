select name, arn, log_group_name
from aws.aws_cloudwatch_log_stream
where name = '{{ resourceName }}' and log_group_name = '{{ resourceName }}'
