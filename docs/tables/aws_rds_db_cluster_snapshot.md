# Table: aws_rds_db_cluster_snapshot

Amazon RDS creates a storage volume snapshot of DB cluster, backing up the entire DB cluster

## Examples

### List of cluster snapshots which are not encrypted

```sql
select
  db_cluster_snapshot_identifier,
  type,
  storage_encrypted,
  split_part(kms_key_id, '/', 1) kms_key_id
from
  aws_rds_db_cluster_snapshot
where
  not storage_encrypted;
```


### Db cluster info of each snapshot

```sql
select
  db_cluster_snapshot_identifier,
  cluster_create_time,
  engine,
  engine_version,
  license_model
from
  aws_rds_db_cluster_snapshot;
```


### Db cluster snapshot count per db cluster

```sql
select
  db_cluster_identifier,
  count(db_cluster_snapshot_identifier) snapshot_count
from
  aws_rds_db_cluster_snapshot
group by
  db_cluster_identifier;
```


### List of manual db cluster snapshot

```sql
select
  db_cluster_snapshot_identifier,
  engine,
  type
from
  aws_rds_db_cluster_snapshot
where
  type = 'manual';
```
