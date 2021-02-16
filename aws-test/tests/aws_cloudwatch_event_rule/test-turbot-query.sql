select akas, title, partition, region, account_id
from aws.aws_cloudwatch_event_rule
where name = '{{ resourceName }}'
