select title, akas, region, account_id
from aws.aws_kinesis_consumer
where consumer_name = '{{ resourceName }}';
