---
title: "Table: aws_elasticache_cluster - Query Amazon ElastiCache Cluster using SQL"
description: "Allows users to query Amazon ElastiCache Cluster data, providing information about each ElastiCache Cluster within the AWS account."
---

# Table: aws_elasticache_cluster - Query Amazon ElastiCache Cluster using SQL

The `aws_elasticache_cluster` table in Steampipe provides information about each ElastiCache Cluster within the AWS account. This table allows DevOps engineers, database administrators, and other IT professionals to query cluster-specific details, including configuration, status, and associated metadata. Users can utilize this table to gather insights on clusters, such as their availability zones, cache node types, engine versions, and more. The schema outlines the various attributes of the ElastiCache Cluster, including the cluster ID, creation date, current status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_cluster` table, you can use the `.inspect aws_elasticache_cluster` command in Steampipe.

Key columns:

- `cache_cluster_id`: The identifier for the cache cluster. This can be used to join this table with other tables to gather more detailed information about the cache cluster.
- `engine`: The name of the cache engine (`memcached` or `redis`) used by the cluster. This information is crucial when joining with tables related to specific engine configurations or statistics.
- `arn`: The Amazon Resource Name (ARN) of the cache cluster. This unique identifier is important for joining with other AWS resource tables that reference the cache cluster.

## Examples

### List clusters that are not encrypted at rest

```sql
select
  cache_cluster_id,
  cache_node_type,
  at_rest_encryption_enabled
from
  aws_elasticache_cluster
where
  not at_rest_encryption_enabled;
```

### List clusters whose availability zone count is less than 2

```sql
select
  cache_cluster_id,
  preferred_availability_zone
from
  aws_elasticache_cluster
where
  preferred_availability_zone <> 'Multiple';
```

### List clusters that do not enforce encryption in transit

```sql
select
  cache_cluster_id,
  cache_node_type,
  transit_encryption_enabled
from
  aws_elasticache_cluster
where
  not transit_encryption_enabled;
```

### List clusters provisioned with undesired (for example, cache.m5.large and cache.m4.4xlarge are desired) node types

```sql
select
  cache_node_type,
  count(*) as count
from
  aws_elasticache_cluster
where
  cache_node_type not in ('cache.m5.large', 'cache.m4.4xlarge')
group by
  cache_node_type;
```

### List clusters with inactive notification configuration topics

```sql
select
  cache_cluster_id,
  cache_cluster_status,
  notification_configuration ->> 'TopicArn' as topic_arn,
  notification_configuration ->> 'TopicStatus' as topic_status
from
  aws_elasticache_cluster
where
  notification_configuration ->> 'TopicStatus' = 'inactive';
```

### Get security group details for each cluster

```sql
select
  cache_cluster_id,
  sg ->> 'SecurityGroupId' as security_group_id,
  sg ->> 'Status' as status
from
  aws_elasticache_cluster,
  jsonb_array_elements(security_groups) as sg;
```

### List clusters with automatic backup disabled

```sql
select
  cache_cluster_id,
  cache_node_type,
  cache_cluster_status,
  snapshot_retention_limit
from
  aws_elasticache_cluster
where
  snapshot_retention_limit is null;
```
