select cache_subnet_group_name, akas, title
from aws.aws_elasticache_subnet_group
where cache_subnet_group_name = '{{ output.resource_id.value }}';