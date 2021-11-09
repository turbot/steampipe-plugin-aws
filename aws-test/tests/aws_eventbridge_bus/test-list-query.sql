select arn, name, partition, region, tags, title
from aws.aws_eventbridge_bus
where akas::text = '["{{ output.resource_arn.value }}"]';
