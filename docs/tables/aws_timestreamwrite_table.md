---
title: "Steampipe Table: aws_timestreamwrite_table - Query AWS Timestream Tables using SQL"
description: "Allows users to query AWS Timestream tables, providing detailed information on table configurations, statuses, and retention properties."
folder: "Timestream"
---

# Table: aws_timestreamwrite_table - Query AWS Timestream Tables using SQL

AWS Timestream is a fast, scalable, and fully managed time-series database service for IoT and operational applications. It is designed to store and analyze trillions of events per day at a fraction of the cost of relational databases. The `aws_timestreamwrite_table` table in Steampipe allows you to query information about Timestream tables in your AWS environment. This includes details like table status, creation time, retention properties, and more.

## Table Usage Guide

The `aws_timestreamwrite_table` table enables DevOps engineers, cloud administrators, and data analysts to gather detailed insights on their Timestream tables. You can query various aspects of the table, such as its schema, retention policies, and status. This table is particularly useful for monitoring table health, ensuring data retention compliance, and managing table configurations.

## Examples

### Basic info
Retrieve basic information about your AWS Timestream tables, including their name, ARN, status, and creation time. This can be useful for getting an overview of the tables deployed in your AWS account.

```sql+postgres
select
  table_name,
  arn,
  table_status,
  creation_time,
  last_updated_time,
  region
from
  aws_timestreamwrite_table;
```

```sql+sqlite
select
  table_name,
  arn,
  table_status,
  creation_time,
  last_updated_time,
  region
from
  aws_timestreamwrite_table;
```

### List active tables
Fetch a list of tables that are currently active. This can help in identifying which tables are in use and available for writing and querying data.

```sql+postgres
select
  table_name,
  arn,
  table_status
from
  aws_timestreamwrite_table
where
  table_status = 'ACTIVE';
```

```sql+sqlite
select
  table_name,
  arn,
  table_status
from
  aws_timestreamwrite_table
where
  table_status = 'ACTIVE';
```

### List tables with specific retention settings
Query tables that have specific retention settings in the memory store or magnetic store. This is useful for ensuring that your data retention policies are being enforced properly.

```sql+postgres
select
  table_name,
  arn,
  retention_properties
from
  aws_timestreamwrite_table
where
  retention_properties ->> 'MemoryStoreRetentionPeriodInHours' = '24'
  and retention_properties ->> 'MagneticStoreRetentionPeriodInDays' = '7';
```

```sql+sqlite
select
  table_name,
  arn,
  retention_properties
from
  aws_timestreamwrite_table
where
  json_extract(retention_properties, '$.MemoryStoreRetentionPeriodInHours') = '24'
  and json_extract(retention_properties, '$.MagneticStoreRetentionPeriodInDays') = '7';
```

### List tables with magnetic store writes enabled
Identify tables where magnetic store writes are enabled. This can help in understanding which tables are set up for long-term storage and potentially lower-cost storage.

```sql+postgres
select
  table_name,
  arn,
  magnetic_store_write_properties
from
  aws_timestreamwrite_table
where
  magnetic_store_write_properties ->> 'EnableMagneticStoreWrites' = 'true';
```

```sql+sqlite
select
  table_name,
  arn,
  magnetic_store_write_properties
from
  aws_timestreamwrite_table
where
  json_extract(magnetic_store_write_properties, '$.EnableMagneticStoreWrites') = 'true';
```

### List tables by creation date
Retrieve tables ordered by their creation date, which can be useful for auditing purposes or understanding the lifecycle of your Timestream tables.

```sql+postgres
select
  table_name,
  arn,
  creation_time
from
  aws_timestreamwrite_table
order by
  creation_time desc;
```

```sql+sqlite
select
  table_name,
  arn,
  creation_time
from
  aws_timestreamwrite_table
order by
  creation_time desc;
```

### Get table schema details
Query the schema of your Timestream tables to understand the structure and types of data that are being stored.

```sql+postgres
select
  table_name,
  arn,
  schema
from
  aws_timestreamwrite_table;
```

```sql+sqlite
select
  table_name,
  arn,
  schema
from
  aws_timestreamwrite_table;
```