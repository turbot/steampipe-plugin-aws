---
title: "Steampipe Table: aws_elasticache_cluster - Query Amazon ElastiCache Cluster using SQL"
description: "Allows users to query Amazon ElastiCache Cluster data, providing information about each ElastiCache Cluster within the AWS account."
folder: "ElastiCache"
---

# Table: aws_elasticache_cluster - Query Amazon ElastiCache Cluster using SQL

The Amazon ElastiCache Cluster is a part of AWS's ElastiCache service that offers fully managed in-memory data store and cache services. This resource is designed to improve the performance of web applications by allowing you to retrieve information from fast, managed, in-memory caches, instead of relying solely on slower disk-based databases. ElastiCache supports two open-source in-memory caching engines: Memcached and Redis.

## Table Usage Guide

The `aws_elasticache_cluster` table in Steampipe provides you with information about each ElastiCache Cluster within your AWS account. This table enables you, as a DevOps engineer, database administrator, or other IT professional, to query cluster-specific details, including configuration, status, and associated metadata. You can utilize this table to gather insights on clusters, such as their availability zones, cache node types, engine versions, and more. The schema outlines the various attributes of the ElastiCache Cluster for you, including the cluster ID, creation date, current status, and associated tags.

## Examples

### List clusters that are not encrypted at rest
Determine the areas in which data clusters are lacking proper encryption at rest. This is essential for identifying potential security vulnerabilities and ensuring data protection compliance.

```sql+postgres
select
  cache_cluster_id,
  cache_node_type,
  at_rest_encryption_enabled
from
  aws_elasticache_cluster
where
  not at_rest_encryption_enabled;
```

```sql+sqlite
select
  cache_cluster_id,
  cache_node_type,
  at_rest_encryption_enabled
from
  aws_elasticache_cluster
where
  at_rest_encryption_enabled = 0;
```

### List clusters whose availability zone count is less than 2
Determine the areas in which your AWS ElastiCache clusters are potentially vulnerable due to having less than two availability zones. This could be useful for improving disaster recovery strategies and ensuring high availability.

```sql+postgres
select
  cache_cluster_id,
  preferred_availability_zone
from
  aws_elasticache_cluster
where
  preferred_availability_zone <> 'Multiple';
```

```sql+sqlite
select
  cache_cluster_id,
  preferred_availability_zone
from
  aws_elasticache_cluster
where
  preferred_availability_zone <> 'Multiple';
```

### List clusters that do not enforce encryption in transit
Determine the areas in your system where encryption in transit is not enforced. This is useful for identifying potential security risks and ensuring that all data is properly protected during transmission.

```sql+postgres
select
  cache_cluster_id,
  cache_node_type,
  transit_encryption_enabled
from
  aws_elasticache_cluster
where
  not transit_encryption_enabled;
```

```sql+sqlite
select
  cache_cluster_id,
  cache_node_type,
  transit_encryption_enabled
from
  aws_elasticache_cluster
where
  transit_encryption_enabled = 0;
```

### List clusters provisioned with undesired (for example, cache.m5.large and cache.m4.4xlarge are desired) node types
Identify instances where clusters have been provisioned with undesired node types, enabling you to streamline your resources and align with your preferred configurations. This is particularly useful for maintaining consistency and optimizing performance across your infrastructure.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which clusters have inactive notification configurations to assess the elements within your system that may not be receiving important updates or alerts.

```sql+postgres
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

```sql+sqlite
select
  cache_cluster_id,
  cache_cluster_status,
  json_extract(notification_configuration, '$.TopicArn') as topic_arn,
  json_extract(notification_configuration, '$.TopicStatus') as topic_status
from
  aws_elasticache_cluster
where
  json_extract(notification_configuration, '$.TopicStatus') = 'inactive';
```

### Get security group details for each cluster
Determine the security status of each cluster by examining the associated security group details. This can help in evaluating the security posture of your clusters and identifying any potential vulnerabilities.

```sql+postgres
select
  cache_cluster_id,
  sg ->> 'SecurityGroupId' as security_group_id,
  sg ->> 'Status' as status
from
  aws_elasticache_cluster,
  jsonb_array_elements(security_groups) as sg;
```

```sql+sqlite
select
  cache_cluster_id,
  json_extract(sg.value, '$.SecurityGroupId') as security_group_id,
  json_extract(sg.value, '$.Status') as status
from
  aws_elasticache_cluster,
  json_each(security_groups) as sg;
```

### List clusters with automatic backup disabled
Determine the areas in which automatic backups are disabled for your clusters. This is useful for ensuring data safety and minimizing the risk of data loss.

```sql+postgres
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

```sql+sqlite
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

### Find all ElastiCache update actions for clusters
Retrieve a list of all service updates that are associated with ElastiCache clusters.

```sql+postgres
select
  c.cache_cluster_id,
  c.engine,
  c.engine_version,
  c.cache_node_type,
  a.service_update_name,
  a.service_update_severity,
  a.service_update_status,
  a.service_update_type,
  a.service_update_recommended_apply_by_date
from
  aws_elasticache_cluster as c
  right join aws_elasticache_update_action as a on c.cache_cluster_id = a.cache_cluster_id
where
  a.service_update_status = 'available'
order by
  a.service_update_severity,
  a.service_update_recommended_apply_by_date;
```

```sql+sqlite
select
  c.cache_cluster_id,
  c.engine,
  c.engine_version,
  c.cache_node_type,
  a.service_update_name,
  a.service_update_severity,
  a.service_update_status,
  a.service_update_type,
  a.service_update_recommended_apply_by_date
from
  aws_elasticache_cluster as c
  join aws_elasticache_update_action as a on c.cache_cluster_id = a.cache_cluster_id
where
  a.service_update_status = 'available'
order by
  a.service_update_severity,
  a.service_update_recommended_apply_by_date;
```