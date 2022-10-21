# Table: aws_dynamodb_table

Amazon DynamoDB is a key-value and document database that delivers single-digit millisecond performance at any scale. It's a fully managed, multi-region, multi-active, durable database with built-in security, backup and restore, and in-memory caching for internet-scale applications.

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
