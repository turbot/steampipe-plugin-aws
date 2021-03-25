select cache_parameter_group_name, akas,title
from aws.aws_elasticache_parameter_group
where cache_parameter_group_name = '{{ output.resource_name.value }}';