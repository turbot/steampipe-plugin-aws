# Table: aws_docdb_cluster_snapshot

A storage volume snapshot of the cluster, backing up the entire cluster

## Examples

### List of cluster snapshots which are not encrypted

```sql
select
  db_cluster_snapshot_identifier,
  type,
  storage_encrypted,
  split_part(kms_key_id, '/', 1) kms_key_id
from
  aws_docdb_cluster_snapshot
where
  not storage_encrypted;
```


### Cluster info of each snapshot

```sql
select
  db_cluster_snapshot_identifier,
  cluster_create_time,
  engine,
  engine_version
from
  aws_docdb_cluster_snapshot;
```


### Cluster snapshot count per cluster

```sql
select
  db_cluster_identifier,
  count(db_cluster_snapshot_identifier) snapshot_count
from
  aws_docdb_cluster_snapshot
group by
  db_cluster_identifier;
```


### List of manual cluster snapshot

```sql
select
  db_cluster_snapshot_identifier,
  engine,
  type
from
  aws_docdb_cluster_snapshot
where
  type = 'manual';
```
