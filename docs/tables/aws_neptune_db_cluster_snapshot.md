# Table: aws_neptune_db_cluster_snapshot

Amazon Neptune creates a storage volume snapshot of your DB cluster, backing up the entire DB cluster and not just individual databases.

## Examples

### List of DB cluster snapshots which are not encrypted

```sql
select
  db_cluster_snapshot_identifier,
  snapshot_type,
  storage_encrypted
from
  aws_neptune_db_cluster_snapshot
where
  not storage_encrypted;
```


### DB cluster info of each snapshot

```sql
select
  db_cluster_snapshot_identifier,
  cluster_create_time,
  engine,
  engine_version,
  license_model
from
  aws_neptune_db_cluster_snapshot;
```


### DB cluster snapshot count per DB cluster

```sql
select
  db_cluster_identifier,
  count(db_cluster_snapshot_identifier) snapshot_count
from
  aws_neptune_db_cluster_snapshot
group by
  db_cluster_identifier;
```


### List of publicly restorable DB cluster snapshot

```sql
select
  db_cluster_snapshot_identifier,
  engine,
  snapshot_type
from
  aws_neptune_db_cluster_snapshot,
  jsonb_array_elements(db_cluster_snapshot_attributes) as cluster_snapshot
where
  cluster_snapshot -> 'AttributeValues' = '["all"]';
```
