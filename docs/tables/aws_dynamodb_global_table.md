---
title: "Table: aws_dynamodb_global_table - Query AWS DynamoDB Global Tables using SQL"
description: "Allows users to query AWS DynamoDB Global Tables to gather information about the global tables, including the table name, creation time, status, and other related details."
---

# Table: aws_dynamodb_global_table - Query AWS DynamoDB Global Tables using SQL

The `aws_dynamodb_global_table` table in Steampipe provides information about Global Tables within AWS DynamoDB. This table allows DevOps engineers to query global table-specific details, including the table name, creation time, status, and other related details. Users can utilize this table to gather insights on global tables, such as the tables' replication status, their regions, and more. The schema outlines the various attributes of the DynamoDB Global Table, including the table ARN, creation time, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dynamodb_global_table` table, you can use the `.inspect aws_dynamodb_global_table` command in Steampipe.

**Key columns**:

- `name`: This is the name of the DynamoDB Global Table. It is a unique identifier and can be used to join this table with other tables that contain DynamoDB Global Table names.
- `arn`: This is the Amazon Resource Name (ARN) of the DynamoDB Global Table. It can be used to join this table with other tables that contain DynamoDB Global Table ARNs.
- `region`: This is the AWS region of the DynamoDB Global Table. It can be used to join this table with other tables that contain AWS region information.

## Examples

### List of regions where global table replicas are present

```sql
select
  global_table_name,
  rg -> 'RegionName' as region_name
from
  aws_dynamodb_global_table
  cross join jsonb_array_elements(replication_group) as rg;
```


### DynamoDB global table replica info

```sql
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
