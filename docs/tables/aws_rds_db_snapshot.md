---
title: "Steampipe Table: aws_rds_db_snapshot - Query Amazon RDS DB Snapshots using SQL"
description: "Allows users to query Amazon RDS DB Snapshots for information regarding manual and automatic snapshots of an Amazon RDS DB instance."
folder: "RDS"
---

# Table: aws_rds_db_snapshot - Query Amazon RDS DB Snapshots using SQL

The AWS RDS DB Snapshot is a feature of Amazon RDS that enables you to create backups of your database instances. These snapshots are point-in-time copies of your databases that can be used for disaster recovery, database migration, and improving backup compliance. The DB Snapshot captures the entire DB instance and not just individual databases, ensuring a consistent snapshot of all your databases at a specific time.

## Table Usage Guide

The `aws_rds_db_snapshot` table in Steampipe provides you with information about manual and automatic snapshots of an Amazon RDS DB instance. This table allows you as a DevOps engineer to query snapshot-specific details, such as snapshot type, creation time, allocated storage, and associated metadata. You can utilize this table to gather insights into snapshot details, including whether a snapshot is shared, public, or encrypted, its engine version, and more. The schema outlines the various attributes of the DB snapshot for you, including the snapshot ARN, DB instance identifier, snapshot status, and associated tags.

## Examples

### DB snapshot basic info
Explore which database snapshots in your AWS RDS service are not encrypted. This can help you identify potential security risks and ensure compliance with encryption policies.

```sql+postgres
select
  db_snapshot_identifier,
  encrypted
from
  aws_rds_db_snapshot
where
  not encrypted;
```

```sql+sqlite
select
  db_snapshot_identifier,
  encrypted
from
  aws_rds_db_snapshot
where
  encrypted = 0;
```


### List of all manual DB snapshots
Discover the segments that consist of all manually created database snapshots, which can help in tracking and managing your backups effectively. This is particularly useful in scenarios where you want to ensure that all your important data is being manually backed up as per your organization's policies.

```sql+postgres
select
  db_snapshot_identifier,
  type
from
  aws_rds_db_snapshot
where
  type = 'manual';
```

```sql+sqlite
select
  db_snapshot_identifier,
  type
from
  aws_rds_db_snapshot
where
  type = 'manual';
```


### List of snapshots which are not encrypted
Determine the areas in which your AWS RDS database snapshots are lacking encryption. This is useful for identifying potential security vulnerabilities and ensuring compliance with data protection policies.

```sql+postgres
select
  db_snapshot_identifier,
  encrypted
from
  aws_rds_db_snapshot
where
  not encrypted;
```

```sql+sqlite
select
  db_snapshot_identifier,
  encrypted
from
  aws_rds_db_snapshot
where
  encrypted = 0;
```


### DB instance info of each db snapshot
Determine the areas in which specific database snapshots are associated with their respective database instances. This query can be beneficial for understanding the storage and engine details of each snapshot, helping in efficient resource management and optimization.

```sql+postgres
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

```sql+sqlite
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