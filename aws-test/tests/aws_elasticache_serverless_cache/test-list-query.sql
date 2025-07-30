select akas, serverless_cache_name, engine, title
from aws.aws_elasticache_serverless_cache
where akas::text = '["{{ output.resource_aka.value }}"]'; 