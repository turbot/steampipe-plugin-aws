---
title: "Steampipe Table: aws_dynamodb_global_table - Query AWS DynamoDB Global Tables using SQL"
description: "Allows users to query AWS DynamoDB Global Tables to gather information about the global tables, including the table name, creation time, status, and other related details."
folder: "DynamoDB"
---

# Table: aws_dynamodb_global_table - Query AWS DynamoDB Global Tables using SQL

The AWS DynamoDB Global Table is a fully managed, multi-region, and multi-active database that provides fast, reliable and secure in-memory data storage and retrieval with seamless scalability. It allows for replication of your Amazon DynamoDB tables in one or more AWS regions, enabling you to access your data from any of these regions and to recover from region-wide failures. This service is suitable for all applications that need to run with low latency and high availability.

## Table Usage Guide

The `aws_dynamodb_global_table` table in Steampipe provides you with information about Global Tables within AWS DynamoDB. This table allows you, as a DevOps engineer, to query global table-specific details, including the table name, creation time, status, and other related details. You can utilize this table to gather insights on global tables, such as the tables' replication status, their regions, and more. The schema outlines for you the various attributes of the DynamoDB Global Table, including the table ARN, creation time, status, and associated tags.

## Examples

### List of regions where global table replicas are present
Discover the segments that have global table replicas in different regions. This is useful for understanding the geographical distribution of your DynamoDB global tables.

```sql+postgres
select
  global_table_name,
  rg -> 'RegionName' as region_name
from
  aws_dynamodb_global_table
  cross join jsonb_array_elements(replication_group) as rg;
```

```sql+sqlite
select
  policy_id,
  arn,
  date_created,
  policy_type,
  json_extract(json_extract(s.value, '$.RetainRule'), '$.Count') as retain_count
from
  aws_dlm_lifecycle_policy,
  json_each(json_extract(policy_details, '$.Schedules')) as s
where 
  json_extract(s.value, '$.RetainRule') is not null;
```


### DynamoDB global table replica info
Explore the status and progress of global replicas in your DynamoDB service. This can help in identifying any inconsistencies or issues in the global data distribution, enabling you to take necessary actions for maintaining data consistency and availability.

```sql+postgres
select
  global_table_name,
  global_table_status,
  rg -> 'GlobalSecondaryIndexes' as global_secondary_indexes,
  rg -> 'RegionName' as region_name,
  rg -> 'ReplicaInaccessibleDateTime' as replica_inaccessible_date_time,
  rg -> 'ReplicaStatus' as replica_status,
  rg -> 'ReplicaStatusDescription' as replica_status_description,
  rg -> 'ReplicaStatusPercentProgress' as replica_status_percent_progress
from
  aws_dynamodb_global_table
  cross join jsonb_array_elements(replication_group) as rg;
```

```sql+sqlite
select
  global_table_name,
  global_table_status,
  json_extract(rg.value, '$.GlobalSecondaryIndexes') as global_secondary_indexes,
  json_extract(rg.value, '$.RegionName') as region_name,
  json_extract(rg.value, '$.ReplicaInaccessibleDateTime') as replica_inaccessible_date_time,
  json_extract(rg.value, '$.ReplicaStatus') as replica_status,
  json_extract(rg.value, '$.ReplicaStatusDescription') as replica_status_description,
  json_extract(rg.value, '$.ReplicaStatusPercentProgress') as replica_status_percent_progress
from
  aws_dynamodb_global_table,
  json_each(replication_group) as rg;
```