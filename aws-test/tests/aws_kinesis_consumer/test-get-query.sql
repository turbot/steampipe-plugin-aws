select consumer_name, consumer_arn, consumer_status
from aws.aws_kinesis_consumer
where consumer_arn = '{{ output.resource_aka.value }}';