select cache_cluster_id, cache_node_type, engine, tags_src
from aws.aws_elasticache_cluster
where cache_cluster_id = '{{ output.resource_id.value }}';
