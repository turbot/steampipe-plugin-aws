select stream_arn
from aws.aws_kinesis_consumer
where consumer_name = '{{ resourceName }}';