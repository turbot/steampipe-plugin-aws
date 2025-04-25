---
title: "Steampipe Table: aws_rds_db_cluster_snapshot - Query AWS RDS DB Cluster Snapshots using SQL"
description: "Allows users to query AWS RDS DB Cluster Snapshots for detailed information on each snapshot, such as the snapshot identifier, creation time, status, and more."
folder: "RDS"
---

# Table: aws_rds_db_cluster_snapshot - Query AWS RDS DB Cluster Snapshots using SQL

The AWS RDS DB Cluster Snapshot is a feature of Amazon RDS that enables you to create a point-in-time snapshot of your database cluster. These snapshots are user-initiated backups of your entire DB Instance, capturing data at a particular moment in time. They can be used for backups, database replication, or for troubleshooting purposes.

## Table Usage Guide

The `aws_rds_db_cluster_snapshot` table in Steampipe provides you with information about DB cluster snapshots within Amazon Relational Database Service (RDS). This table allows you, as a DevOps engineer or database administrator, to query snapshot-specific details, including snapshot status, creation time, engine version, and associated metadata. You can utilize this table to gather insights on snapshots, such as snapshot availability, storage used, and source DB cluster identifier. The schema outlines the various attributes of the DB cluster snapshot for you, including the snapshot ARN, snapshot type, VPC ID, and associated tags.

## Examples

### List of cluster snapshots which are not encrypted
Identify instances where your cluster snapshots are not encrypted. This is crucial to uncover potential security risks and ensure data protection compliance within your AWS RDS clusters.

```sql+postgres
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

```sql+sqlite
select
  db_cluster_snapshot_identifier,
  type,
  storage_encrypted,
  substr(kms_key_id, 1, instr(kms_key_id, '/') - 1) as kms_key_id
from
  aws_rds_db_cluster_snapshot
where
  not storage_encrypted;
```


### Db cluster info of each snapshot
Discover the specifics of each database cluster snapshot, such as its creation time, engine type, version, and licensing model. This can be useful in understanding the historical configuration and performance of your database clusters.

```sql+postgres
select
  db_cluster_snapshot_identifier,
  cluster_create_time,
  engine,
  engine_version,
  license_model
from
  aws_rds_db_cluster_snapshot;
```

```sql+sqlite
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
Explore the distribution of snapshots across different database clusters. This can be useful for understanding backup habits and ensuring that data is being adequately protected across all clusters.

```sql+postgres
select
  db_cluster_identifier,
  count(db_cluster_snapshot_identifier) snapshot_count
from
  aws_rds_db_cluster_snapshot
group by
  db_cluster_identifier;
```

```sql+sqlite
select
  db_cluster_identifier,
  count(db_cluster_snapshot_identifier) as snapshot_count
from
  aws_rds_db_cluster_snapshot
group by
  db_cluster_identifier;
```


### List of manual db cluster snapshot
Explore which database cluster snapshots have been manually created within your AWS RDS service. This could be useful to track and manage backup strategies or to validate compliance with internal policies regarding data persistence.

```sql+postgres
select
  db_cluster_snapshot_identifier,
  engine,
  type
from
  aws_rds_db_cluster_snapshot
where
  type = 'manual';
```

```sql+sqlite
select
  db_cluster_snapshot_identifier,
  engine,
  type
from
  aws_rds_db_cluster_snapshot
where
  type = 'manual';
```