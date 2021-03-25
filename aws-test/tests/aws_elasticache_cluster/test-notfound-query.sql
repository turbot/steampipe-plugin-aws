select title, akas, region, account_id
from aws.aws_elasticache_cluster
where cache_cluster_id = '{{ output.resource_id.value }}-dummy';