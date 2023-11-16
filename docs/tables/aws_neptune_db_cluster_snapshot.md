---
title: "Table: aws_neptune_db_cluster_snapshot - Query AWS Neptune DB Cluster Snapshots using SQL"
description: "Allows users to query AWS Neptune DB Cluster Snapshots for comprehensive details about their configurations, status, and associated metadata."
---

# Table: aws_neptune_db_cluster_snapshot - Query AWS Neptune DB Cluster Snapshots using SQL

The `aws_neptune_db_cluster_snapshot` table in Steampipe provides information about DB Cluster Snapshots within Amazon Neptune. This table allows DevOps engineers, database administrators, and other technical professionals to query snapshot-specific details, including snapshot status, creation time, associated database engine, and more. Users can utilize this table to gather insights on snapshots, such as their availability, encryption status, and associated database clusters. The schema outlines the various attributes of the Neptune DB Cluster Snapshot, including the snapshot ARN, creation time, associated tags, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_neptune_db_cluster_snapshot` table, you can use the `.inspect aws_neptune_db_cluster_snapshot` command in Steampipe.

**Key columns**:

- `db_cluster_snapshot_identifier`: The identifier for the DB cluster snapshot. This column can be used to join this table with others that contain information about specific DB cluster snapshots.
- `db_cluster_identifier`: The identifier of the DB cluster that the snapshot was created from. This column can be used to join this table with others that contain information about specific DB clusters.
- `snapshot_type`: The type of DB cluster snapshot. This column can be used to distinguish between manual and automated snapshots.

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

### List of publicly restorable DB cluster snapshots

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
