---
title: "Table: aws_rds_db_snapshot - Query Amazon RDS DB Snapshots using SQL"
description: "Allows users to query Amazon RDS DB Snapshots for information regarding manual and automatic snapshots of an Amazon RDS DB instance."
---

# Table: aws_rds_db_snapshot - Query Amazon RDS DB Snapshots using SQL

The `aws_rds_db_snapshot` table in Steampipe provides information about manual and automatic snapshots of an Amazon RDS DB instance. This table allows DevOps engineers to query snapshot-specific details, such as snapshot type, creation time, allocated storage, and associated metadata. Users can utilize this table to gather insights into snapshot details, including whether a snapshot is shared, public, or encrypted, its engine version, and more. The schema outlines the various attributes of the DB snapshot, including the snapshot ARN, DB instance identifier, snapshot status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_snapshot` table, you can use the `.inspect aws_rds_db_snapshot` command in Steampipe.

**Key columns**:

- `db_snapshot_identifier`: The identifier for the DB snapshot. This column can be used to join this table with other tables as it uniquely identifies each snapshot.
- `db_instance_identifier`: The identifier for the DB instance that was used to create the snapshot. This is useful for joining with the `aws_rds_db_instance` table to get more information about the DB instance.
- `snapshot_create_time`: The time when the snapshot was taken. This is useful for tracking the age of snapshots and managing snapshot lifecycle.

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
