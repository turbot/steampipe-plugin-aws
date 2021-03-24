select akas, title
from aws.aws_cloudwatch_log_stream
where name = '{{ resourceName }}'
