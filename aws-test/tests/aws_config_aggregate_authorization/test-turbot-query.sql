select title, tags, akas
from aws.aws_config_aggregate_authorization
where akas::text = '["{{ output.resource_aka.value }}"]';