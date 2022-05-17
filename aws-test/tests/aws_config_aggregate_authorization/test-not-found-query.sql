select arn, tags_src, title
from aab_aab.aws_config_aggregate_authorization
where akas::text = 'dummy-["{{ output.resource_aka.value }}"]';