select stream_arn
from aws_kinesis_consumer
where consumer_name = '{{ resourceName }}';