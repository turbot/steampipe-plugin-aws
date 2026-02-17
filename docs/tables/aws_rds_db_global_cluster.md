---
title: "Steampipe Table: aws_rds_db_global_cluster - Query AWS RDS Global Clusters using SQL"
description: "Allows users to query AWS RDS Global Clusters and retrieve details about global database topology, status, and member clusters."
folder: "RDS"
---

# Table: aws_rds_db_global_cluster - Query AWS RDS Global Clusters using SQL

An AWS RDS Global Cluster (Aurora Global Database) is an overarching layer that groups regional Aurora DB clusters into a single globally distributed database.

## Table Usage Guide

The `aws_rds_db_global_cluster` table in Steampipe provides information about RDS global clusters in your account. You can use this table to audit global database configurations, check failover readiness, and inspect member clusters and encryption settings.

## Examples

### List RDS global clusters

```sql+postgres
select
  global_cluster_identifier,
  status,
  endpoint,
  engine,
  engine_version
from
  aws_rds_db_global_cluster;
```

```sql+sqlite
select
  global_cluster_identifier,
  status,
  endpoint,
  engine,
  engine_version
from
  aws_rds_db_global_cluster;
```

### Show global cluster members

```sql+postgres
select
  global_cluster_identifier,
  member ->> 'DBClusterArn' as db_cluster_arn,
  member ->> 'IsWriter' as is_writer,
  member ->> 'Readers' as readers
from
  aws_rds_db_global_cluster
  cross join jsonb_array_elements(global_cluster_members) as member;
```

```sql+sqlite
select
  global_cluster_identifier,
  json_extract(member.value, '$.DBClusterArn') as db_cluster_arn,
  json_extract(member.value, '$.IsWriter') as is_writer,
  json_extract(member.value, '$.Readers') as readers
from
  aws_rds_db_global_cluster,
  json_each(global_cluster_members) as member;
```

### Find unencrypted global clusters

```sql+postgres
select
  global_cluster_identifier,
  storage_encrypted
from
  aws_rds_db_global_cluster
where
  not storage_encrypted;
```

```sql+sqlite
select
  global_cluster_identifier,
  storage_encrypted
from
  aws_rds_db_global_cluster
where
  storage_encrypted = 0;
```
