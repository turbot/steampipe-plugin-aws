select arn, tags_src, title
from aws.aws_config_aggregate_authorization
where akas::text = 'dummy-["{{ output.resource_aka.value }}"]';