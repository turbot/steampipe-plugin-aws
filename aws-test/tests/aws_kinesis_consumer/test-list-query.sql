select consumer_name, consumer_arn
from aws.aws_kinesis_consumer
where akas::text = '["{{ output.resource_aka.value }}"]';