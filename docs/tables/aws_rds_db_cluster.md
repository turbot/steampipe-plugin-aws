---
title: "Steampipe Table: aws_rds_db_cluster - Query AWS RDS DB Clusters using SQL"
description: "Allows users to query AWS RDS DB Clusters and retrieve valuable information about the status, configuration, and security settings of each DB cluster."
folder: "RDS"
---

# Table: aws_rds_db_cluster - Query AWS RDS DB Clusters using SQL

The AWS RDS DB Cluster is a component of Amazon Relational Database Service (RDS). It is a virtual database where multiple DB instances are associated under a single endpoint. This allows for efficient scaling and management of databases, providing high availability and failover support for DB instances.

## Table Usage Guide

The `aws_rds_db_cluster` table in Steampipe provides you with information about DB clusters within Amazon Relational Database Service (RDS). This table allows you, as a DevOps engineer, to query DB cluster-specific details, including configuration, status, and security settings. You can utilize this table to gather insights on DB clusters, such as their availability, backup settings, encryption status, and more. The schema outlines the various attributes of the DB cluster for you, including the DB cluster identifier, creation time, DB cluster members, and associated tags.

## Examples

### List of DB clusters which are not encrypted
Discover the segments of your database clusters that lack encryption. This is crucial for identifying potential security vulnerabilities within your AWS RDS database clusters.

```sql+postgres
select
  db_cluster_identifier,
  allocated_storage,
  kms_key_id
from
  aws_rds_db_cluster
where
  kms_key_id is null;
```

```sql+sqlite
select
  db_cluster_identifier,
  allocated_storage,
  kms_key_id
from
  aws_rds_db_cluster
where
  kms_key_id is null;
```

### List of DB clusters where backup retention period is greater than 7 days
Explore which database clusters have a backup retention period set for more than a week. This can be useful for identifying databases that have longer data retention policies, potentially indicating important or sensitive data.

```sql+postgres
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_rds_db_cluster
where
  backup_retention_period > 7;
```

```sql+sqlite
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_rds_db_cluster
where
  backup_retention_period > 7;
```

### Avalability zone count for each db instance
Determine the areas in which each database cluster is available by counting the availability zones. This can be useful for understanding the spread and redundancy of your databases across different geographical zones.

```sql+postgres
select
  db_cluster_identifier,
  jsonb_array_length(availability_zones) availability_zones_count
from
  aws_rds_db_cluster;
```

```sql+sqlite
select
  db_cluster_identifier,
  json_array_length(json(availability_zones)) as availability_zones_count
from
  aws_rds_db_cluster;
```

### DB cluster Members info
Explore the configuration of your database clusters to understand the status of each member, their roles, and their promotion tiers. This can help optimize the performance and reliability of your cloud databases.

```sql+postgres
select
  db_cluster_identifier,
  member ->> 'DBClusterParameterGroupStatus' as db_cluster_parameter_group_status,
  member ->> 'DBInstanceIdentifier' as db_instance_identifier,
  member ->> 'IsClusterWriter' as is_cluster_writer,
  member ->> 'PromotionTier' as promotion_tier
from
  aws_rds_db_cluster
  cross join jsonb_array_elements(members) as member;
```

```sql+sqlite
select
  db_cluster_identifier,
  json_extract(member.value, '$.DBClusterParameterGroupStatus') as db_cluster_parameter_group_status,
  json_extract(member.value, '$.DBInstanceIdentifier') as db_instance_identifier,
  json_extract(member.value, '$.IsClusterWriter') as is_cluster_writer,
  json_extract(member.value, '$.PromotionTier') as promotion_tier
from
  aws_rds_db_cluster,
  json_each(members) as member;
```

### List DB cluster pending maintenance actions
List DB clusters pending maintenance actions to plan and prioritize maintenance schedules effectively.

```sql+postgres
select
  a.db_cluster_identifier,
  b.action,
  a.status,
  b.opt_in_status,
  b.forced_apply_date,
  b.current_apply_date,
  b.auto_applied_after_date
from 
  aws_rds_db_cluster as a
  join aws_rds_pending_maintenance_action as b on b.resource_identifier = a.arn;
```

```sql+sqlite
select
  a.db_cluster_identifier,
  b.action,
  a.status,
  b.opt_in_status,
  b.forced_apply_date,
  b.current_apply_date,
  b.auto_applied_after_date
from 
  aws_rds_db_cluster as a
  join aws_rds_pending_maintenance_action as b on b.resource_identifier = a.arn;
```