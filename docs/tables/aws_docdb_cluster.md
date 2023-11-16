---
title: "Table: aws_docdb_cluster - Query Amazon DocumentDB Cluster using SQL"
description: "Allows users to query Amazon DocumentDB Clusters for detailed information about their configuration, status, and associated metadata."
---

# Table: aws_docdb_cluster - Query Amazon DocumentDB Cluster using SQL

The `aws_docdb_cluster` table in Steampipe provides information about Amazon DocumentDB clusters within AWS. This table allows DevOps engineers, database administrators, and other technical professionals to query cluster-specific details, including configurations, status, and associated metadata. Users can utilize this table to gather insights on clusters, such as their availability, backup and restore settings, encryption status, and more. The schema outlines the various attributes of the DocumentDB cluster, including the cluster ARN, creation time, DB subnet group, associated VPC, and backup retention period.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_docdb_cluster` table, you can use the `.inspect aws_docdb_cluster` command in Steampipe.

**Key columns**:

- `cluster_identifier`: This is the unique identifier for the Amazon DocumentDB cluster. It is crucial for joining this table with others as it uniquely identifies each cluster.
- `vpc_security_groups`: This column provides information about the VPC security groups associated with the cluster. It is useful for querying security configurations and joining with other security-related tables.
- `db_subnet_group`: This column contains information about the DB subnet group associated with the cluster. It is important for understanding the network configuration of the cluster and can be used to join with network-related tables.

## Examples

## Basic Info

```sql
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

```sql
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

### List clusters where backup retention period is greater than 7 days

```sql
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_docdb_cluster
where
  backup_retention_period > 7;
```

### Get avalability zone count for each cluster

```sql
select
  db_cluster_identifier,
  jsonb_array_length(availability_zones) as availability_zones_count
from
  aws_docdb_cluster;
```

### List cluster members details

```sql
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

### List clusters where deletion protection is disabled

```sql
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