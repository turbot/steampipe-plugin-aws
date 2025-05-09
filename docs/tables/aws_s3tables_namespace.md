---
title: "Steampipe Table: aws_s3tables_namespace - Query AWS S3 Tables Namespaces using SQL"
description: "Allows users to query AWS S3 Tables Namespaces, providing information about the configuration, settings, and properties of your S3 namespaces."
folder: "S3Tables"
---

# aws_s3tables_namespace

Amazon S3 Tables provide storage that's optimized for analytics workloads, with built-in Apache Iceberg support and features designed to continuously improve query performance and reduce storage costs for tables. Namespaces in S3 Tables help organize and categorize tables, providing a logical grouping mechanism.

The `aws_s3tables_namespace` table provides insights into namespaces within your AWS account. This is useful for monitoring namespace configurations, analyzing organizational structures, and ensuring proper resource ownership and access control.

## Table Usage Guide

The `aws_s3tables_namespace` table provides information about namespaces in Amazon S3 Tables within your AWS account. As a data engineer or cloud administrator, this table helps you manage and monitor your namespace resources. You can use this table to identify namespaces, analyze their properties, verify owner information, and understand how tables are organized within your S3 Tables environment.

## Examples

### Basic info

Retrieves fundamental information about all S3 Tables namespaces in your AWS account, including their names, IDs, creation dates, and owner account IDs.

```sql+postgresql
select
  namespace,
  namespace_id,
  created_at,
  owner_account_id,
  table_bucket_arn
from
  aws_s3tables_namespace;
```

```sql+sqlite
select
  namespace,
  namespace_id,
  created_at,
  owner_account_id,
  table_bucket_arn
from
  aws_s3tables_namespace;
```

### Find recently created namespaces

Identify namespaces that have been created recently, which may indicate new data organization initiatives or projects.

```sql+postgresql
select
  namespace,
  namespace_id,
  created_at,
  table_bucket_arn
from
  aws_s3tables_namespace
where
  created_at > (current_date - interval '30 days')
order by
  created_at desc;
```

```sql+sqlite
select
  namespace,
  namespace_id,
  created_at,
  table_bucket_arn
from
  aws_s3tables_namespace
where
  created_at > datetime('now', '-30 days')
order by
  created_at desc;
```

### Count tables in each namespace

Analyze the distribution of tables across namespaces to understand your data organization.

```sql+postgresql
select
  n.namespace,
  count(t.name) as table_count
from
  aws_s3tables_namespace n
left join
  aws_s3tables_table t on t.namespace_id = n.namespace_id
group by
  n.namespace
order by
  table_count desc;
```

```sql+sqlite
select
  n.namespace,
  count(t.name) as table_count
from
  aws_s3tables_namespace n
left join
  aws_s3tables_table t on t.namespace_id = n.namespace_id
group by
  n.namespace
order by
  table_count desc;
```

### Get namespaces and their associated table buckets

Examine the relationships between namespaces and table buckets to understand your S3 Tables resource organization.

```sql+postgresql
select
  n.namespace,
  n.namespace_id,
  n.table_bucket_arn
from
  aws_s3tables_namespace n
order by
  namespace;
```

```sql+sqlite
select
  n.namespace,
  n.namespace_id,
  n.table_bucket_arn
from
  aws_s3tables_namespace n
order by
  namespace;
```

### Find namespaces in a specific table bucket

Identify all namespaces within a specific table bucket to understand the namespace organization.

```sql+postgresql
select
  namespace,
  namespace_id,
  created_at,
  owner_account_id,
  table_bucket_arn
from
  aws_s3tables_namespace
where
  table_bucket_arn = 'arn:aws:s3tables:us-east-1:123456789012:tablebucket/my-table-bucket';
```

```sql+sqlite
select
  namespace,
  namespace_id,
  created_at,
  owner_account_id,
  table_bucket_arn
from
  aws_s3tables_namespace
where
  table_bucket_arn = 'arn:aws:s3tables:us-east-1:123456789012:tablebucket/my-table-bucket';
```

### Join namespaces with table buckets and tables for a complete hierarchy

Provides a comprehensive view of your S3 Tables organization, showing the relationships between table buckets, namespaces, and tables.

```sql+postgresql
select
  b.name as bucket_name,
  n.namespace as namespace_name,
  t.name as table_name,
  t.type as table_type,
  t.format as table_format,
  n.created_at as namespace_created_at,
  t.created_at as table_created_at
from
  aws_s3tables_namespace n
  join aws_s3tables_table_bucket b on n.table_bucket_id = b.table_bucket_id
  left join aws_s3tables_table t on n.namespace_id = t.namespace_id
order by
  b.name, n.namespace, t.name;
```

```sql+sqlite
select
  b.name as bucket_name,
  n.namespace as namespace_name,
  t.name as table_name,
  t.type as table_type,
  t.format as table_format,
  n.created_at as namespace_created_at,
  t.created_at as table_created_at
from
  aws_s3tables_namespace n
  join aws_s3tables_table_bucket b on n.table_bucket_id = b.table_bucket_id
  left join aws_s3tables_table t on n.namespace_id = t.namespace_id
order by
  b.name, n.namespace, t.name;
```
