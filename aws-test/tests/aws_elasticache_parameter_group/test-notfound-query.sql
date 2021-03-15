select title, akas, region, account_id
from aws.aws_elasticache_parameter_group
where cache_parameter_group_name = '{{ output.resource_name.value }}-dummy';