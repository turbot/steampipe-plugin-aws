# Table: aws_elasticache_replication_group

A Redis replication group is a collection of cache clusters, where one of the clusters is a primary read-write cluster and the others are read-only replicas.

## Examples

### Basic info

```sql
select
  replication_group_id,
  description,
  cache_node_type,
  cluster_enabled,
  auth_token_enabled,
  automatic_failover
from
  aws_elasticache_replication_group;
```


### List replication groups that are not encrypted at rest

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


### List replication groups whose multi AZ feature is not enabled

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


### List replication groups whose backup retention period is less than 30 days

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


### List replication groups by node type

```sql
select
  cache_node_type,
  count (*)
from
  aws_elasticache_replication_group
group by
  cache_node_type;
```


### List member clusters for each replication group

```sql
select
  replication_group_id,
  jsonb_array_elements(member_clusters) as member_clusters
from
  aws_elasticache_replication_group;
```