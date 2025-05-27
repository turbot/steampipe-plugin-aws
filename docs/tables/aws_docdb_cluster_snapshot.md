---
title: "Steampipe Table: aws_docdb_cluster_snapshot - Query Amazon DocumentDB Cluster Snapshot using SQL"
description: "Allows users to query Amazon DocumentDB Cluster Snapshots for detailed information about their configuration, status, and associated metadata."
folder: "DocumentDB"
---

# Table: aws_docdb_cluster_snapshot - Query Amazon DocumentDB Cluster Snapshots using SQL

The `aws_docdb_cluster_snapshot` table provides detailed information about snapshots of Amazon DocumentDB clusters. These snapshots are storage volume snapshots that back up the entire cluster, enabling data recovery and historical analysis.

## Table Usage Guide

This table allows DevOps engineers, database administrators, and other technical professionals to query detailed information about Amazon DocumentDB cluster snapshots. Utilize this table to analyze snapshot configurations, encryption statuses, and other metadata. The schema includes attributes of the DocumentDB cluster snapshots, such as identifiers, creation times, and the associated cluster details.

## Examples

### List of cluster snapshots that are not encrypted
Identify unencrypted cluster snapshots to assess and improve your security posture.

```sql+postgres
select
  db_cluster_snapshot_identifier,
  snapshot_type,
  not storage_encrypted as storage_not_encrypted,
  split_part(kms_key_id, '/', 1) as kms_key_id
from
  aws_docdb_cluster_snapshot
where
  not storage_encrypted;
```

```sql+sqlite
select
  db_cluster_snapshot_identifier,
  snapshot_type,
  not storage_encrypted as storage_not_encrypted,
  substr(kms_key_id, 1, instr(kms_key_id, '/') - 1) as kms_key_id
from
  aws_docdb_cluster_snapshot
where
  not storage_encrypted;
```

### Cluster information of each snapshot
Retrieve basic information about each cluster snapshot, including its creation time and the engine details.

```sql+postgres
select
  db_cluster_snapshot_identifier,
  cluster_create_time,
  engine,
  engine_version
from
  aws_docdb_cluster_snapshot;
```

```sql+sqlite
select
  db_cluster_snapshot_identifier,
  cluster_create_time,
  engine,
  engine_version
from
  aws_docdb_cluster_snapshot;
```

### Cluster snapshot count per cluster
Determine the number of snapshots taken for each cluster to help manage snapshot policies and storage.

```sql+postgres
select
  db_cluster_identifier,
  count(db_cluster_snapshot_identifier) as snapshot_count
from
  aws_docdb_cluster_snapshot
group by
  db_cluster_identifier;
```

```sql+sqlite
select
  db_cluster_identifier,
  count(db_cluster_snapshot_identifier) as snapshot_count
from
  aws_docdb_cluster_snapshot
group by
  db_cluster_identifier;
```

### List of manual cluster snapshots
Filter for manually created cluster snapshots to distinguish them from automatic backups.

```sql+postgres
select
  db_cluster_snapshot_identifier,
  engine,
  snapshot_type
from
  aws_docdb_cluster_snapshot
where
  snapshot_type = 'manual';
```

```sql+sqlite
select
  db_cluster_snapshot_identifier,
  engine,
  snapshot_type
from
  aws_docdb_cluster_snapshot
where
  snapshot_type = 'manual';
```
