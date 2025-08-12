select serverless_cache_name, engine, tags_src, cache_usage_limits
from aws.aws_elasticache_serverless_cache
where serverless_cache_name = '{{ output.resource_id.value }}'; 