---
title: "Steampipe Table: aws_docdb_cluster - Query Amazon DocumentDB Cluster using SQL"
description: "Allows users to query Amazon DocumentDB Clusters for detailed information about their configuration, status, and associated metadata."
folder: "DocumentDB"
---

# Table: aws_docdb_cluster - Query Amazon DocumentDB Cluster using SQL

The Amazon DocumentDB Cluster is a fully managed, MongoDB compatible database service designed for workloads that need high availability, reliability, and scalability. It allows you to store, query, and index JSON data. DocumentDB makes it easy to operate mission critical MongoDB workloads at scale.

## Table Usage Guide

The `aws_docdb_cluster` table in Steampipe provides you with information about Amazon DocumentDB clusters within AWS. This table allows you as a DevOps engineer, database administrator, or other technical professional to query cluster-specific details, including configurations, status, and associated metadata. You can utilize this table to gather insights on clusters, such as their availability, backup and restore settings, encryption status, and more. The schema outlines the various attributes of the DocumentDB cluster for you, including the cluster ARN, creation time, DB subnet group, associated VPC, and backup retention period.

## Examples

## Basic Info

```sql+postgres
select
  arn,
  db_cluster_identifier,
  deletion_protection,
  engine,
  status,
  region
from
  aws_docdb_cluster;
```

```sql+sqlite
select
  arn,
  db_cluster_identifier,
  deletion_protection,
  engine,
  status,
  region
from
  aws_docdb_cluster;
```

### List clusters which are not encrypted
Discover the segments that are not encrypted within your database clusters. This can help enhance your security measures by identifying potential vulnerabilities.

```sql+postgres
select
  db_cluster_identifier,
  status,
  cluster_create_time,
  kms_key_id,
  storage_encrypted
from
  aws_docdb_cluster
where
  not storage_encrypted;
```

```sql+sqlite
select
  db_cluster_identifier,
  status,
  cluster_create_time,
  kms_key_id,
  storage_encrypted
from
  aws_docdb_cluster
where
  storage_encrypted = 0;
```

### List clusters where backup retention period is greater than 7 days
Identify instances where the backup retention period for database clusters exceeds a week. This could be useful in managing data storage and ensuring compliance with data retention policies.

```sql+postgres
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_docdb_cluster
where
  backup_retention_period > 7;
```

```sql+sqlite
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_docdb_cluster
where
  backup_retention_period > 7;
```

### Get avalability zone count for each cluster
Determine the number of availability zones for each database cluster in your AWS DocumentDB service to better manage and distribute your databases across different zones for high availability and fault tolerance.

```sql+postgres
select
  db_cluster_identifier,
  jsonb_array_length(availability_zones) as availability_zones_count
from
  aws_docdb_cluster;
```

```sql+sqlite
select
  db_cluster_identifier,
  json_array_length(availability_zones) as availability_zones_count
from
  aws_docdb_cluster;
```

### List clusters where deletion protection is disabled
Discover the segments that have deletion protection disabled in order to identify potential vulnerabilities and enhance security measures. This is particularly useful in maintaining data integrity by preventing accidental deletions.

```sql+postgres
select
  db_cluster_identifier,
  status,
  cluster_create_time,
  deletion_protection
from
  aws_docdb_cluster
where
  not deletion_protection;
```

```sql+sqlite
select
  db_cluster_identifier,
  status,
  cluster_create_time,
  deletion_protection
from
  aws_docdb_cluster
where
  deletion_protection = 0;
```

### List cluster members details
Identify instances where you can assess the status and roles of members within your AWS DocumentDB clusters. This enables you to understand the configuration of each cluster member, including their promotion tier and whether they have write access.

```sql+postgres
select
  db_cluster_identifier,
  member ->> 'DBClusterParameterGroupStatus' as db_cluster_parameter_group_status,
  member ->> 'DBInstanceIdentifier' as db_instance_identifier,
  member ->> 'IsClusterWriter' as is_cluster_writer,
  member ->> 'PromotionTier' as promotion_tier
from
  aws_docdb_cluster
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
  aws_docdb_cluster,
  json_each(members) as member;
```

### List clusters where deletion protection is disabled
Determine the areas in which deletion protection is disabled for your clusters. This can help in identifying potential vulnerabilities and ensuring your data is secure.

```sql+postgres
select
  db_cluster_identifier,
  status,
  cluster_create_time,
  deletion_protection
from
  aws_docdb_cluster
where
  not deletion_protection;
```

```sql+sqlite
select
  db_cluster_identifier,
  status,
  cluster_create_time,
  deletion_protection
from
  aws_docdb_cluster
where
  not deletion_protection = 0;
```