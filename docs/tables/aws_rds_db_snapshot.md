# Table: aws_rds_db_snapshot

Amazon RDS creates a storage volume snapshot of your DB instance, backing up the entire DB instance and not just individual databases.

## Examples

### DB snapshot basic info

```sql
select
  db_snapshot_identifier,
  encrypted
from
  aws_rds_db_snapshot
where
  not encrypted;
```


### List of all manual DB snapshots

```sql
select
  db_snapshot_identifier,
  type
from
  aws_rds_db_snapshot
where
  type = 'manual';
```


### List of snapshots which are not encrypted

```sql
select
  db_snapshot_identifier,
  encrypted
from
  aws_rds_db_snapshot
where
  not encrypted;
```


### DB instance info of each db snapshot

```sql
select
  db_snapshot_identifier,
  db_instance_identifier,
  engine,
  engine_version,
  allocated_storage,
  storage_type
from
  aws_rds_db_snapshot;
```
