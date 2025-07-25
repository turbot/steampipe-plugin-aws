select serverless_cache_name, akas, tags, title
from aws.aws_elasticache_serverless_cache
where serverless_cache_name = '{{ output.resource_id.value }}'; 