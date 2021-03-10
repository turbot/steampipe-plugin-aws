select akas, cache_cluster_id, cache_node_type, engine, title
from aws.aws_elasticache_cluster
where akas::text = '["{{ output.resource_aka.value }}"]';
