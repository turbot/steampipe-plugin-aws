select title, akas, region, account_id
from aws.aws_elasticache_serverless_cache
where serverless_cache_name = '{{ output.resource_id.value }}-dummy'; 