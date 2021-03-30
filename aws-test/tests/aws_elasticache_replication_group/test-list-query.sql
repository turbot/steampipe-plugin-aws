select akas, replication_group_id, automatic_failover, description, cache_node_type, title
from aws.aws_elasticache_replication_group
where akas::text = '["{{ output.resource_aka.value }}"]';