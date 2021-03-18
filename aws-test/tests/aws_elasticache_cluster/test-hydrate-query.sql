select cache_cluster_id, akas, tags, title
from aws.aws_elasticache_cluster
where cache_cluster_id = '{{ output.resource_id.value }}';
