select tags, akas, title, partition, region, account_id
from aws.aws_cloudwatch_log_group
where name = '{{ resourceName }}'