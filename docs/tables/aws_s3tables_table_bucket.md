---
title: "Steampipe Table: aws_s3tables_table_bucket - Query AWS S3 Tables Table Buckets using SQL"
description: "Allows users to query AWS S3 Tables Table Buckets, providing information about the configuration, settings, and properties of your S3 table buckets."
folder: "S3Tables"
---

# aws_s3tables_table_bucket

Amazon S3 Tables provide storage that's optimized for analytics workloads, with built-in Apache Iceberg support and features designed to continuously improve query performance and reduce storage costs for tables. S3 Tables store tabular data in table buckets, which are purpose-built for tables and provide higher transactions per second (TPS) and better query throughput compared to self-managed tables in S3 general purpose buckets.

The `aws_s3tables_table_bucket` table provides insights into table buckets within your AWS account. This is useful for monitoring table bucket configurations, analyzing storage properties, and ensuring proper resource ownership and access control.

## Table Usage Guide

The `aws_s3tables_table_bucket` table provides information about S3 Tables table buckets in your AWS account. As a data engineer or cloud administrator, this table helps you manage and monitor your table bucket resources. You can use this table to identify table buckets, analyze their properties, and verify owner information.

## Examples

### Basic info
Retrieves fundamental information about all S3 Tables table buckets in your AWS account, including their names, ARNs, creation dates, and owner account IDs. This provides a quick overview of all table buckets you have access to.

```sql+postgresql
select
  name,
  arn,
  created_at,
  owner_account_id
from
  aws_s3tables_table_bucket;
```

```sql+sqlite
select
  name,
  arn,
  created_at,
  owner_account_id
from
  aws_s3tables_table_bucket;
```

### Get details for a specific table bucket
Show detailed information about a specific table bucket by name. This is useful when you need to examine the properties of a particular table bucket.

```sql+postgresql
select
  name,
  arn,
  created_at,
  table_bucket_id,
  owner_account_id
from
  aws_s3tables_table_bucket
where
  name = 'my-table-bucket';
```

```sql+sqlite
select
  name,
  arn,
  created_at,
  table_bucket_id,
  owner_account_id
from
  aws_s3tables_table_bucket
where
  name = 'my-table-bucket';
```

### Count table buckets by region
Aggregates and counts table buckets by AWS region, ordering the results by count in descending order. This helps you understand the distribution of your table buckets across different regions.

```sql+postgresql
select
  region,
  count(*) as bucket_count
from
  aws_s3tables_table_bucket
group by
  region
order by
  bucket_count desc;
```

```sql+sqlite
select
  region,
  count(*) as bucket_count
from
  aws_s3tables_table_bucket
group by
  region
order by
  bucket_count desc;
```

### Find table buckets created in the last 30 days
Identifies table buckets that were created within the last 30 days, ordered by creation date. This is helpful for tracking recent bucket creations and monitoring new resources.

```sql+postgresql
select
  name,
  arn,
  created_at,
  region
from
  aws_s3tables_table_bucket
where
  created_at > (current_date - interval '30' day)
order by
  created_at desc;
```

```sql+sqlite
select
  name,
  arn,
  created_at,
  region
from
  aws_s3tables_table_bucket
where
  created_at > datetime('now', '-30 day')
order by
  created_at desc;
```

### List table buckets belonging to a specific account
Table buckets to show only those owned by a specific AWS account. This is useful for multi-account environments where you need to identify resources associated with a particular account.

```sql+postgresql
select
  name,
  arn,
  region,
  created_at
from
  aws_s3tables_table_bucket
where
  owner_account_id = '123456789012';
```

```sql+sqlite
select
  name,
  arn,
  region,
  created_at
from
  aws_s3tables_table_bucket
where
  owner_account_id = '123456789012';
```

### Join with tables and namespaces to view complete hierarchy
Shows the hierarchical relationship between table buckets, namespaces, and tables, providing a complete view of your S3 Tables resources.

```sql+postgresql
select
  b.name as bucket_name,
  n.namespace as namespace_name,
  t.name as table_name,
  t.type as table_type,
  t.format as table_format,
  t.created_at as table_created_at
from
  aws_s3tables_table_bucket b
  left join aws_s3tables_namespace n on b.table_bucket_id = n.table_bucket_id
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
  t.created_at as table_created_at
from
  aws_s3tables_table_bucket b
  left join aws_s3tables_namespace n on b.table_bucket_id = n.table_bucket_id
  left join aws_s3tables_table t on n.namespace_id = t.namespace_id
order by
  b.name, n.namespace, t.name;
```
