select name
from aws.aws_cloudwatch_event_rule
where name = '{{ resourceName }}'
