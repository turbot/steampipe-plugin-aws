---
title: "Steampipe Table: aws_memorydb_cluster - Query AWS MemoryDB Clusters using SQL"
description: "Allows users to query AWS MemoryDB clusters, providing detailed information on cluster configurations and statuses."
folder: "MemoryDB"
---

# Table: aws_memorydb_cluster - Query AWS MemoryDB Clusters using SQL

AWS MemoryDB is a Redis-compatible, fully managed, in-memory database service that delivers ultra-fast performance for modern applications. It is designed for building and running real-time applications that require sub-millisecond latency. The `aws_memorydb_cluster` table in Steampipe allows you to query information about MemoryDB clusters in your AWS environment. This includes details like cluster configuration, status, security, and maintenance settings.

## Table Usage Guide

The `aws_memorydb_cluster` table enables DevOps engineers and cloud administrators to gather detailed insights on their MemoryDB clusters. You can query various aspects of the cluster, such as its endpoint configuration, encryption settings, and shard details. This table is particularly useful for monitoring cluster health, ensuring security compliance, and managing cluster configurations.

## Examples

### Basic cluster information
Retrieve basic information about your AWS MemoryDB clusters, including their name, ARN, status, and node type. This can be useful for getting an overview of the clusters deployed in your AWS account.

```sql+postgres
select
  name,
  arn,
  status,
  node_type,
  engine_version,
  region
from
  aws_memorydb_cluster;
```

```sql+sqlite
select
  name,
  arn,
  status,
  node_type,
  engine_version,
  region
from
  aws_memorydb_cluster;
```

### List clusters with auto-minor version upgrade enabled
Identify clusters that have automatic minor version upgrades enabled, which can be useful for understanding how your clusters are being maintained and ensuring they are kept up to date.

```sql+postgres
select
  name,
  arn,
  auto_minor_version_upgrade
from
  aws_memorydb_cluster
where
  auto_minor_version_upgrade = true;
```

```sql+sqlite
select
  name,
  arn,
  auto_minor_version_upgrade
from
  aws_memorydb_cluster
where
  auto_minor_version_upgrade = 1;
```

### List multi-AZ clusters
Fetch a list of clusters that are configured with Multi-AZ for high availability. This can help in identifying which clusters have enhanced availability features enabled.

```sql+postgres
select
  name,
  arn,
  availability_mode
from
  aws_memorydb_cluster
where
  availability_mode = 'multiaz';
```

```sql+sqlite
select
  name,
  arn,
  availability_mode
from
  aws_memorydb_cluster
where
  availability_mode = 'multiaz';
```

### List clusters with encryption in transit disabled
Find clusters where encryption in transit is not enabled, which may indicate potential security risks that need to be addressed.

```sql+postgres
select
  name,
  arn,
  tls_enabled
from
  aws_memorydb_cluster
where
  tls_enabled = false;
```

```sql+sqlite
select
  name,
  arn,
  tls_enabled
from
  aws_memorydb_cluster
where
  tls_enabled = 0;
```

### Clusters by maintenance window
Retrieve clusters along with their scheduled maintenance windows, which is useful for planning updates and understanding potential downtime.

```sql+postgres
select
  name,
  arn,
  maintenance_window
from
  aws_memorydb_cluster;
```

```sql+sqlite
select
  name,
  arn,
  maintenance_window
from
  aws_memorydb_cluster;
```

### List clusters with specific node type
Query clusters that are using a particular node type, allowing you to evaluate the resource allocations and performance characteristics across your MemoryDB clusters.

```sql+postgres
select
  name,
  arn,
  node_type
from
  aws_memorydb_cluster
where
  node_type = 'db.r6gd.xlarge';
```

```sql+sqlite
select
  name,
  arn,
  node_type
from
  aws_memorydb_cluster
where
  node_type = 'db.r6gd.xlarge';
```

### List clusters with shard details
Retrieve detailed information about the shards within each cluster, including shard configuration and number of shards. This can help in understanding the distribution of data and load across your clusters.

```sql+postgres
select
  name,
  arn,
  number_of_shards,
  shards
from
  aws_memorydb_cluster;
```

```sql+sqlite
select
  name,
  arn,
  number_of_shards,
  shards
from
  aws_memorydb_cluster;
```

### Clusters with pending updates
Identify clusters that have pending updates, which may require attention to ensure the clusters are up-to-date and running optimally.

```sql+postgres
select
  name,
  arn,
  pending_updates
from
  aws_memorydb_cluster
where
  jsonb_array_length(pending_updates) > 0;
```

```sql+sqlite
select
  name,
  arn,
  pending_updates
from
  aws_memorydb_cluster
where
  json_array_length(pending_updates) > 0;
```

### List clusters with snapshot retention details
Gather information on the snapshot retention settings of your clusters to ensure that your data backup policies are being properly enforced.

```sql+postgres
select
  name,
  arn,
  snapshot_retention_limit,
  snapshot_window
from
  aws_memorydb_cluster;
```

```sql+sqlite
select
  name,
  arn,
  snapshot_retention_limit,
  snapshot_window
from
  aws_memorydb_cluster;
```
