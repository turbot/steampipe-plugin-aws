select akas, cache_subnet_group_name, cache_subnet_group_description, vpc_id, title
from aws.aws_elasticache_subnet_group
where cache_subnet_group_name = '{{ output.resource_id.value }}'; 