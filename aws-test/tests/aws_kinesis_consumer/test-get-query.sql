select consumer_name, consumer_arn
from aws_kinesis_consumer
where consumer_arn = '{{ output.resource_aka.value }}';