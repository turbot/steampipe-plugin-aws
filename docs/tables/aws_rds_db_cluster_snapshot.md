---
title: "Table: aws_rds_db_cluster_snapshot - Query AWS RDS DB Cluster Snapshots using SQL"
description: "Allows users to query AWS RDS DB Cluster Snapshots for detailed information on each snapshot, such as the snapshot identifier, creation time, status, and more."
---

# Table: aws_rds_db_cluster_snapshot - Query AWS RDS DB Cluster Snapshots using SQL

The `aws_rds_db_cluster_snapshot` table in Steampipe provides information about DB cluster snapshots within Amazon Relational Database Service (RDS). This table allows DevOps engineers and database administrators to query snapshot-specific details, including snapshot status, creation time, engine version, and associated metadata. Users can utilize this table to gather insights on snapshots, such as snapshot availability, storage used, and source DB cluster identifier. The schema outlines the various attributes of the DB cluster snapshot, including the snapshot ARN, snapshot type, VPC ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_cluster_snapshot` table, you can use the `.inspect aws_rds_db_cluster_snapshot` command in Steampipe.

### Key columns:

- `db_cluster_snapshot_identifier`: The identifier for the DB cluster snapshot. This column is important as it is the unique identifier for each snapshot and can be used to join this table with other tables that need snapshot-specific information.
- `db_cluster_identifier`: The DB cluster identifier. This column is useful for joining with tables that provide information at the DB cluster level.
- `snapshot_create_time`: The time when the snapshot was taken. This column is useful for tracking snapshot history and understanding the lifecycle of your DB clusters.


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
