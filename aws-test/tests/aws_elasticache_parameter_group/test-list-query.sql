select akas, cache_parameter_group_name, cache_parameter_group_family, title
from aws.aws_elasticache_parameter_group
where akas::text = '["{{ output.resource_aka.value }}"]';