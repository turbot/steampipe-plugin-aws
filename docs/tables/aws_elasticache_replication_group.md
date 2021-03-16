# Table: aws_elasticache_replication_group

Provides an ElastiCache Replication Group resource.

## Examples

### List of unencrypted elasticache replication groups

```sql
select
  replication_group_id,
  cache_node_type,
  at_rest_encryption_enabled
from
  aws_elasticache_replication_group
where
  not at_rest_encryption_enabled;
```


### List of elasticache replication groups whose multi AZ feature is not enabled

```sql
select
  replication_group_id,
  cache_node_type,
  multi_az
from
  aws_elasticache_replication_group
where
  multi_az = 'disabled';
```


### List of elasticache replication groups whose backup retention period is less than 30 days

```sql
select
  replication_group_id,
  snapshot_retention_limit,
  snapshot_window,
  snapshotting_cluster_id
from
  aws_elasticache_replication_group
where
  snapshot_retention_limit < 30;
```


### List of elasticache replication groups which are not enabled for automatic failover

```sql
select
  replication_group_id,
  cache_node_type,
  multi_az
from
  aws_elasticache_replication_group
where
  automatic_failover = 'disabled';
```


### Count of replication group by node type

```sql
select
  cache_node_type,
  count (*)
from
  aws_elasticache_replication_group
group by
  cache_node_type;
```


### List of member clusters for each replication group

```sql
select
  replication_group_id,
  jsonb_array_elements(member_clusters) as member_clusters
from
  aws_elasticache_replication_group;
```

