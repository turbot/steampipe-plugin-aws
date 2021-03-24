select replication_group_id, automatic_failover, cache_node_type, akas, title
from aws.aws_elasticache_replication_group
where replication_group_id = '{{ output.resource_id.value }}';