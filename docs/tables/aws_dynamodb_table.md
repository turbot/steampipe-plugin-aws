---
title: "Steampipe Table: aws_dynamodb_table - Query AWS DynamoDB Tables using SQL"
description: "Allows users to query AWS DynamoDB Tables and retrieve detailed information about their configuration, status, and associated attributes."
folder: "DynamoDB"
---

# Table: aws_dynamodb_table - Query AWS DynamoDB Tables using SQL

The AWS DynamoDB service provides fully managed NoSQL database tables that are designed to provide quick and predictable performance by automatically distributing data across multiple servers. These tables support both key-value and document data models, and enable developers to build web, mobile, and IoT applications without worrying about hardware and setup. DynamoDB tables also offer built-in security, in-memory caching, backup and restore, and in-place update capabilities.

## Table Usage Guide

The `aws_dynamodb_table` table in Steampipe provides you with information about tables within AWS DynamoDB. This table allows you, as a DevOps engineer, to query table-specific details, including provisioned throughput, global secondary indexes, local secondary indexes, and associated metadata. You can utilize this table to gather insights on tables, such as their read/write capacity mode, encryption status, and more. The schema outlines the various attributes of the DynamoDB table for you, including the table name, creation date, item count, and associated tags.

## Examples

### List of Dynamodb tables which are not encrypted with CMK
Identify instances where DynamoDB tables are not encrypted with a Customer Master Key (CMK). This is useful for enhancing security and compliance by ensuring all data is adequately protected.

```sql+postgres
select
  name,
  sse_description
from
  aws_dynamodb_table
where
  sse_description is null;
```

```sql+sqlite
select
  name,
  sse_description
from
  aws_dynamodb_table
where
  sse_description is null;
```


### List of tables where continuous backup is not enabled
Explore which tables have not enabled continuous backup, a critical feature for data loss prevention and recovery in AWS DynamoDB. This can help identify potential vulnerabilities and areas for improvement in your database management practices.

```sql+postgres
select
  name,
  continuous_backups_status
from
  aws_dynamodb_table
where
  continuous_backups_status = 'DISABLED';
```

```sql+sqlite
select
  name,
  continuous_backups_status
from
  aws_dynamodb_table
where
  continuous_backups_status = 'DISABLED';
```


### Point in time recovery info for each table
Determine the areas in which you can restore your AWS DynamoDB tables by identifying the earliest and latest possible recovery times. This is particularly useful in disaster recovery scenarios, where understanding the recovery timeline is crucial.

```sql+postgres
select
  name,
  point_in_time_recovery_description ->> 'EarliestRestorableDateTime' as earliest_restorable_date_time,
  point_in_time_recovery_description ->> 'LatestRestorableDateTime' as latest_restorable_date_time,
  point_in_time_recovery_description ->> 'PointInTimeRecoveryStatus' as point_in_time_recovery_status
from
  aws_dynamodb_table;
```

```sql+sqlite
select
  name,
  json_extract(point_in_time_recovery_description, '$.EarliestRestorableDateTime') as earliest_restorable_date_time,
  json_extract(point_in_time_recovery_description, '$.LatestRestorableDateTime') as latest_restorable_date_time,
  json_extract(point_in_time_recovery_description, '$.PointInTimeRecoveryStatus') as point_in_time_recovery_status
from
  aws_dynamodb_table;
```

### List of tables where streaming is enabled with destination status
Determine the areas in which streaming is enabled and assess the status of these destinations. This is useful for monitoring the health and activity of your streaming destinations.

```sql+postgres
select
  name,
  d ->> 'StreamArn' as kinesis_stream_arn,
  d ->> 'DestinationStatus' as stream_status
from
  aws_dynamodb_table,
  jsonb_array_elements(streaming_destination -> 'KinesisDataStreamDestinations') as d
```

```sql+sqlite
select
  name,
  json_extract(d.value, '$.StreamArn') as kinesis_stream_arn,
  json_extract(d.value, '$.DestinationStatus') as stream_status
from
  aws_dynamodb_table,
  json_each(streaming_destination, 'KinesisDataStreamDestinations') as d
```