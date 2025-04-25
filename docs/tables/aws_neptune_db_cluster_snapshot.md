---
title: "Steampipe Table: aws_neptune_db_cluster_snapshot - Query AWS Neptune DB Cluster Snapshots using SQL"
description: "Allows users to query AWS Neptune DB Cluster Snapshots for comprehensive details about their configurations, status, and associated metadata."
folder: "Neptune"
---

# Table: aws_neptune_db_cluster_snapshot - Query AWS Neptune DB Cluster Snapshots using SQL

AWS Neptune DB Cluster Snapshots are a point-in-time copy of data from an Amazon Neptune DB cluster. These snapshots can be used to restore a cluster to the specific time the snapshot was taken, which is useful for disaster recovery or data analysis purposes. They can be automatically or manually created and deleted within the Amazon Neptune service.

## Table Usage Guide

The `aws_neptune_db_cluster_snapshot` table in Steampipe provides you with information about DB Cluster Snapshots within Amazon Neptune. This table allows you, as a DevOps engineer, database administrator, or other technical professional, to query snapshot-specific details, including snapshot status, creation time, associated database engine, and more. You can utilize this table to gather insights on snapshots, such as their availability, encryption status, and associated database clusters. The schema outlines the various attributes of the Neptune DB Cluster Snapshot for you, including the snapshot ARN, creation time, associated tags, and more.

## Examples

### List of DB cluster snapshots which are not encrypted
Uncover the details of database cluster snapshots that lack encryption. This information is crucial for identifying potential security risks within your system.

```sql+postgres
select
  db_cluster_snapshot_identifier,
  snapshot_type,
  storage_encrypted
from
  aws_neptune_db_cluster_snapshot
where
  not storage_encrypted;
```

```sql+sqlite
select
  db_cluster_snapshot_identifier,
  snapshot_type,
  storage_encrypted
from
  aws_neptune_db_cluster_snapshot
where
  storage_encrypted = 0;
```

### DB cluster info of each snapshot
Explore the creation times, engines used, and licensing models of different database clusters. This is beneficial for understanding the configuration and setup of each database cluster in your AWS Neptune service.

```sql+postgres
select
  db_cluster_snapshot_identifier,
  cluster_create_time,
  engine,
  engine_version,
  license_model
from
  aws_neptune_db_cluster_snapshot;
```

```sql+sqlite
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
Explore the number of snapshots created for each database cluster to assess the frequency of data backup and to ensure data recovery readiness in case of a failure. This is crucial for maintaining data integrity and minimizing potential data loss.

```sql+postgres
select
  db_cluster_identifier,
  count(db_cluster_snapshot_identifier) snapshot_count
from
  aws_neptune_db_cluster_snapshot
group by
  db_cluster_identifier;
```

```sql+sqlite
select
  db_cluster_identifier,
  count(db_cluster_snapshot_identifier) snapshot_count
from
  aws_neptune_db_cluster_snapshot
group by
  db_cluster_identifier;
```

### List of publicly restorable DB cluster snapshots
Discover the segments that include snapshots of your database clusters that can be restored by anyone. This is useful for identifying potential security risks or for planning data recovery strategies.

```sql+postgres
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

```sql+sqlite
select
  db_cluster_snapshot_identifier,
  engine,
  snapshot_type
from
  aws_neptune_db_cluster_snapshot
where
  json_extract(db_cluster_snapshot_attributes, '$.AttributeValues') = '["all"]';
```