select title, tags, akas
from aab_aab.aws_config_aggregate_authorization
where akas::text = '["{{ output.resource_aka.value }}"]';