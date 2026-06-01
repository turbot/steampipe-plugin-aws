---
title: "Steampipe Table: aws_rds_global_cluster - Query AWS RDS Global Clusters using SQL"
description: "Allows users to query AWS RDS Global Clusters and retrieve details about global database topology, status, and member clusters."
folder: "RDS"
---

# Table: aws_rds_global_cluster - Query AWS RDS Global Clusters using SQL

An AWS RDS Global Cluster (Aurora Global Database) is a single globally distributed database that spans multiple AWS Regions, grouping a primary Aurora DB cluster together with one or more read-only secondary clusters in other Regions. It is designed for low-latency global reads and fast cross-Region disaster recovery. A global cluster is an account-level, global resource, so it is returned once per account rather than per Region.

## Table Usage Guide

The `aws_rds_global_cluster` table in Steampipe provides you with information about RDS global clusters within Amazon Relational Database Service (RDS). This table allows you, as a DevOps engineer, to query global cluster-specific details, including topology, status, member clusters, and encryption settings. You can utilize this table to gather insights on global clusters, such as auditing global database configurations, checking failover readiness, and verifying storage encryption. The schema outlines the various attributes of the global cluster for you, including the global cluster identifier, ARN, engine, member clusters, and associated tags. Because global clusters are a global resource, the `region` column is always `global`.

## Examples

### List RDS global clusters
Get an overview of the global database clusters in your account to track their status and engine versions. This helps you understand your global database footprint at a glance.

```sql+postgres
select
  global_cluster_identifier,
  status,
  engine,
  engine_version
from
  aws_rds_global_cluster;
```

```sql+sqlite
select
  global_cluster_identifier,
  status,
  engine,
  engine_version
from
  aws_rds_global_cluster;
```

### Show global cluster members
Inspect the primary and secondary clusters that make up each global cluster to understand its cross-Region topology and identify which member is the writer. This is useful when planning failover or reviewing read replica placement.

```sql+postgres
select
  global_cluster_identifier,
  member ->> 'DBClusterArn' as db_cluster_arn,
  member ->> 'IsWriter' as is_writer,
  member -> 'Readers' as readers
from
  aws_rds_global_cluster
  cross join jsonb_array_elements(global_cluster_members) as member;
```

```sql+sqlite
select
  global_cluster_identifier,
  json_extract(member.value, '$.DBClusterArn') as db_cluster_arn,
  json_extract(member.value, '$.IsWriter') as is_writer,
  json_extract(member.value, '$.Readers') as readers
from
  aws_rds_global_cluster,
  json_each(global_cluster_members) as member;
```

### Find unencrypted global clusters
Identify global clusters that do not have storage encryption enabled so you can address potential data-at-rest security gaps. This helps you enforce encryption compliance across your global databases.

```sql+postgres
select
  global_cluster_identifier,
  storage_encrypted
from
  aws_rds_global_cluster
where
  not storage_encrypted;
```

```sql+sqlite
select
  global_cluster_identifier,
  storage_encrypted
from
  aws_rds_global_cluster
where
  storage_encrypted = 0;
```
