---
title: "Steampipe Table: aws_keyspaces_table - Query AWS Keyspaces Tables using SQL"
description: "Allows users to query AWS Keyspaces tables, providing detailed information on table configurations, throughput capacity, encryption, and more."
folder: "Keyspaces"
---

# Table: aws_keyspaces_table - Query AWS Keyspaces Tables using SQL

Amazon Keyspaces (for Apache Cassandra) is a scalable, highly available, and managed Apache Cassandra-compatible database service. It enables you to run Cassandra workloads on AWS without managing the underlying infrastructure. The `aws_keyspaces_table` table in Steampipe allows you to query information about your Keyspaces tables in AWS, including their capacity specifications, encryption settings, and schema definitions.

## Table Usage Guide

The `aws_keyspaces_table` table enables cloud administrators and DevOps engineers to gather detailed insights into their Keyspaces tables. You can query various aspects of the table, such as its creation timestamp, throughput capacity, encryption, and status. This table is particularly useful for monitoring table performance, ensuring security compliance, and managing table configurations.

## Examples

### Basic table information
Retrieve basic information about your AWS Keyspaces tables, including their name, ARN, status, and region. This can be useful for getting an overview of the tables deployed in your AWS account.

```sql+postgres
select
  table_name,
  arn,
  status,
  creation_timestamp,
  region
from
  aws_keyspaces_table;
```

```sql+sqlite
select
  table_name,
  arn,
  status,
  creation_timestamp,
  region
from
  aws_keyspaces_table;
```

### List active tables
Fetch a list of tables that are currently active. This can help in identifying which tables are operational and available for use.

```sql+postgres
select
  table_name,
  arn,
  status
from
  aws_keyspaces_table
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  table_name,
  arn,
  status
from
  aws_keyspaces_table
where
  status = 'ACTIVE';
```

### List tables with specific encryption settings
Identify tables that have specific encryption settings, which can help ensure that your data is secured according to your organization’s encryption policies.

```sql+postgres
select
  table_name,
  arn,
  encryption_specification_type,
  kms_key_identifier
from
  aws_keyspaces_table
where
  encryption_specification_type = 'AWS_OWNED_KMS_KEY';
```

```sql+sqlite
select
  table_name,
  arn,
  encryption_specification_type,
  kms_key_identifier
from
  aws_keyspaces_table
where
  encryption_specification_type = 'AWS_OWNED_KMS_KEY';
```

### List tables by creation date
Retrieve tables ordered by their creation date, which can be useful for auditing purposes or understanding the lifecycle of your Keyspaces tables.

```sql+postgres
select
  table_name,
  arn,
  creation_timestamp
from
  aws_keyspaces_table
order by
  creation_timestamp desc;
```

```sql+sqlite
select
  table_name,
  arn,
  creation_timestamp
from
  aws_keyspaces_table
order by
  creation_timestamp desc;
```

### List tables with default Time to Live (TTL) settings
Identify tables that have a specific default Time to Live (TTL) setting, which is useful for managing data retention and ensuring compliance with data lifecycle policies.

```sql+postgres
select
  table_name,
  arn,
  default_time_to_live,
  ttl_status
from
  aws_keyspaces_table
where
  default_time_to_live is not null;
```

```sql+sqlite
select
  table_name,
  arn,
  default_time_to_live,
  ttl_status
from
  aws_keyspaces_table
where
  default_time_to_live is not null;
```

### Get table schema definitions
Retrieve detailed schema definitions for your Keyspaces tables, which can help in understanding the structure and types of data stored in the tables.

```sql+postgres
select
  table_name,
  arn,
  schema_definition
from
  aws_keyspaces_table
where
  keyspace_name = 'myKey';
```

```sql+sqlite
select
  table_name,
  arn,
  schema_definition
from
  aws_keyspaces_table
where
  keyspace_name = 'myKey';
```

### List tables with Point-in-Time Recovery (PITR) enabled
Identify tables where Point-in-Time Recovery (PITR) is enabled, which is essential for ensuring that you can restore your table to a previous state within a specified timeframe.

```sql+postgres
select
  table_name,
  arn,
  point_in_time_recovery
from
  aws_keyspaces_table
where
  (point_in_time_recovery ->> 'status') = 'ENABLED';
```

```sql+sqlite
select
  table_name,
  arn,
  point_in_time_recovery
from
  aws_keyspaces_table
where
  json_extract(point_in_time_recovery, '$.status') = 'ENABLED';
```