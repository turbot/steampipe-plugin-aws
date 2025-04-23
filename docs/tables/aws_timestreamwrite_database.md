---
title: "Steampipe Table: aws_timestreamwrite_database - Query AWS Timestream Databases using SQL"
description: "Allows users to query AWS Timestream databases, providing detailed information on database configurations, statuses, and associated tables."
folder: "Timestream"
---

# Table: aws_timestreamwrite_database - Query AWS Timestream Databases using SQL

AWS Timestream is a fast, scalable, and fully managed time-series database service for IoT and operational applications. It is designed to store and analyze trillions of events per day at a fraction of the cost of relational databases. The `aws_timestreamwrite_database` table in Steampipe allows you to query information about Timestream databases in your AWS environment. This includes details like database creation time, encryption settings, table count, and more.

## Table Usage Guide

The `aws_timestreamwrite_database` table enables DevOps engineers, cloud administrators, and data analysts to gather detailed insights on their Timestream databases. You can query various aspects of the database, such as its KMS encryption key, number of tables, and creation time. This table is particularly useful for monitoring database health, ensuring security compliance, and managing database configurations.

## Examples

### Basic database information
Retrieve basic information about your AWS Timestream databases, including their name, ARN, creation time, and region. This can be useful for getting an overview of the databases deployed in your AWS account.

```sql+postgres
select
  database_name,
  arn,
  creation_time,
  region,
  kms_key_id
from
  aws_timestreamwrite_database;
```

```sql+sqlite
select
  database_name,
  arn,
  creation_time,
  region,
  kms_key_id
from
  aws_timestreamwrite_database;
```

### List databases with a specific KMS key
Identify databases that are encrypted with a specific KMS key. This can help in ensuring that your data is secured according to your organization’s encryption policies.

```sql+postgres
select
  database_name,
  arn,
  kms_key_id
from
  aws_timestreamwrite_database
where
  kms_key_id = 'your-kms-key-id';
```

```sql+sqlite
select
  database_name,
  arn,
  kms_key_id
from
  aws_timestreamwrite_database
where
  kms_key_id = 'your-kms-key-id';
```

### List databases by creation date
Retrieve databases ordered by their creation date, which can be useful for auditing purposes or understanding the lifecycle of your Timestream databases.

```sql+postgres
select
  database_name,
  arn,
  creation_time
from
  aws_timestreamwrite_database
order by
  creation_time desc;
```

```sql+sqlite
select
  database_name,
  arn,
  creation_time
from
  aws_timestreamwrite_database
order by
  creation_time desc;
```

### List databases with the most tables
Identify the databases that contain the most tables, which can help in understanding data distribution and load across your Timestream environment.

```sql+postgres
select
  database_name,
  arn,
  table_count
from
  aws_timestreamwrite_database
order by
  table_count desc;
```

```sql+sqlite
select
  database_name,
  arn,
  table_count
from
  aws_timestreamwrite_database
order by
  table_count desc;
```

### Get database details with last updated time
Retrieve detailed information about your databases, including when they were last updated, to monitor changes and updates over time.

```sql+postgres
select
  database_name,
  arn,
  last_updated_time,
  region
from
  aws_timestreamwrite_database;
```

```sql+sqlite
select
  database_name,
  arn,
  last_updated_time,
  region
from
  aws_timestreamwrite_database;
```