select akas, cache_subnet_group_name, cache_subnet_group_description, title
from aws.aws_elasticache_subnet_group
where akas::text = '["{{ output.resource_aka.value }}"]';