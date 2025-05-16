---
title: "Steampipe Table: aws_s3tables_table - Query AWS S3 Tables Tables using SQL"
description: "Allows users to query AWS S3 Tables Tables, providing information about the configuration, settings, and properties of your S3 tables."
folder: "S3Tables"
---

# aws_s3tables_table

Amazon S3 Tables provide storage that's optimized for analytics workloads, with built-in Apache Iceberg support and features designed to continuously improve query performance and reduce storage costs for tables. Tables are structured datasets that define how to interpret data stored in S3.

The `aws_s3tables_table` table provides insights into tables within your AWS account. This is useful for monitoring table configurations, analyzing storage properties, and ensuring proper resource ownership and access control.

## Table Usage Guide

The `aws_s3tables_table` table provides information about tables in Amazon S3 Tables within your AWS account. As a data engineer or cloud administrator, this table helps you manage and monitor your table resources. You can use this table to identify tables, analyze their properties, verify owner information, and understand relationships with namespaces and table buckets.

The table uses a parent/child hydration pattern, listing tables for each S3 table bucket, which means queries will be optimized when filtering by table bucket information.

## Examples

### Basic info
Retrieves fundamental information about all S3 Tables tables in your AWS account, including their names, ARNs, creation dates, and namespace information.

```sql+postgresql
select
  name,
  arn,
  namespace,
  created_at,
  table_bucket_arn
from
  aws_s3tables_table;
```

```sql+sqlite
select
  name,
  arn,
  namespace,
  created_at,
  table_bucket_arn
from
  aws_s3tables_table;
```

### List tables with creation and modification details

View tables along with who created them and when they were last modified.

```sql+postgresql
select
  name,
  created_at,
  created_by,
  modified_at,
  modified_by,
  namespace
from
  aws_s3tables_table
order by
  modified_at desc;
```

```sql+sqlite
select
  name,
  created_at,
  created_by,
  modified_at,
  modified_by,
  namespace
from
  aws_s3tables_table
order by
  modified_at desc;
```

### Find tables in a specific table bucket
Get all tables within a specific table bucket to understand their relationships and structure.

```sql+postgresql
select
  name,
  namespace,
  format,
  created_at,
  type
from
  aws_s3tables_table
where
  table_bucket_arn = 'arn:aws:s3tables:us-east-1:123456789012:tablebucket/my-table-bucket';
```

```sql+sqlite
select
  name,
  namespace,
  format,
  created_at,
  type
from
  aws_s3tables_table
where
  table_bucket_arn = 'arn:aws:s3tables:us-east-1:123456789012:tablebucket/my-table-bucket';
```

### Find recently modified tables
Identify tables that have been modified recently, which may indicate active development or data updates.

```sql+postgresql
select
  name,
  created_at,
  modified_at,
  modified_by,
  table_bucket_arn
from
  aws_s3tables_table
where
  modified_at > (current_date - interval '7 days')
order by
  modified_at desc;
```

```sql+sqlite
select
  name,
  created_at,
  modified_at,
  modified_by,
  table_bucket_arn
from
  aws_s3tables_table
where
  modified_at > datetime('now', '-7 days')
order by
  modified_at desc;
```

### Get table details including format and metadata location
Examine the details of tables including their format and metadata location to better understand how data is stored and accessed.

```sql+postgresql
select
  name,
  namespace,
  format,
  metadata_location,
  warehouse_location,
  version_token
from
  aws_s3tables_table;
```

```sql+sqlite
select
  name,
  namespace,
  format,
  metadata_location,
  warehouse_location,
  version_token
from
  aws_s3tables_table;
```

### Filter tables by namespace name
Find all tables within a specific namespace to analyze related data collections.

```sql+postgresql
select
  name,
  created_at,
  format,
  metadata_location
from
  aws_s3tables_table
where
  namespace = 'my-namespace';
```

```sql+sqlite
select
  name,
  created_at,
  format,
  metadata_location
from
  aws_s3tables_table
where
  namespace = 'my-namespace';
```

### Get namespace info for the tables
Get a complete view of tables with their associated namespaces and table buckets to understand the full hierarchy.

```sql+postgresql
select
  t.name as table_name,
  t.format as table_format,
  t.type as table_type,
  n.namespace as namespace_name,
  b.name as bucket_name,
  t.created_at,
  t.modified_at
from
  aws_s3tables_table t
  join aws_s3tables_namespace n on t.namespace_id = n.namespace_id
  join aws_s3tables_table_bucket b on n.table_bucket_id = b.table_bucket_id
order by
  b.name, n.namespace, t.name;
```

```sql+sqlite
select
  t.name as table_name,
  t.format as table_format,
  t.type as table_type,
  n.namespace as namespace_name,
  b.name as bucket_name,
  t.created_at,
  t.modified_at
from
  aws_s3tables_table t
  join aws_s3tables_namespace n on t.namespace_id = n.namespace_id
  join aws_s3tables_table_bucket b on n.table_bucket_id = b.table_bucket_id
order by
  b.name, n.namespace, t.name;
```
