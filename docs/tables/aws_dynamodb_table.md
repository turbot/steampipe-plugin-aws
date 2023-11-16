---
title: "Table: aws_dynamodb_table - Query AWS DynamoDB Tables using SQL"
description: "Allows users to query AWS DynamoDB Tables and retrieve detailed information about their configuration, status, and associated attributes."
---

# Table: aws_dynamodb_table - Query AWS DynamoDB Tables using SQL

The `aws_dynamodb_table` table in Steampipe provides information about tables within AWS DynamoDB. This table allows DevOps engineers to query table-specific details, including provisioned throughput, global secondary indexes, local secondary indexes, and associated metadata. Users can utilize this table to gather insights on tables, such as their read/write capacity mode, encryption status, and more. The schema outlines the various attributes of the DynamoDB table, including the table name, creation date, item count, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dynamodb_table` table, you can use the `.inspect aws_dynamodb_table` command in Steampipe.

### Key columns:

- `table_name`: This is the name of the DynamoDB table. It is a critical column as it uniquely identifies each table and can be used to join with other tables that require table-specific information.
- `arn`: This is the Amazon Resource Name (ARN) of the DynamoDB table. It is important because it provides a universally unique identifier for the table across all of AWS.
- `status`: This column indicates the current status of the table (e.g., 'ACTIVE'). This is useful for assessing the availability and health of the table.

## Examples

### List of Dynamodb tables which are not encrypted with CMK

```sql
select
  name,
  sse_description
from
  aws_dynamodb_table
where
  sse_description is null;
```


### List of tables where continuous backup is not enabled

```sql
select
  name,
  continuous_backups_status
from
  aws_dynamodb_table
where
  continuous_backups_status = 'DISABLED';
```


### Point in time recovery info for each table

```sql
select
  name,
  point_in_time_recovery_description ->> 'EarliestRestorableDateTime' as earliest_restorable_date_time,
  point_in_time_recovery_description ->> 'LatestRestorableDateTime' as latest_restorable_date_time,
  point_in_time_recovery_description ->> 'PointInTimeRecoveryStatus' as point_in_time_recovery_status
from
  aws_dynamodb_table;
```

### List of tables where streaming is enabled with destination status

```sql
select
  name,
  d ->> 'StreamArn' as kinesis_stream_arn,
  d ->> 'DestinationStatus' as stream_status
from
  aws_dynamodb_table,
  jsonb_array_elements(streaming_destination -> 'KinesisDataStreamDestinations') as d
```
