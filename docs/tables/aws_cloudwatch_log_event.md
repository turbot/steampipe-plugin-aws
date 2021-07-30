# Table: aws_cloudwatch_log_event

A log group is a group of log streams that share the same retention, monitoring, and access control settings. Log streams contain sequences of log events that share the same source.

This table lists events from a CloudWatch log group.

**Important notes:**

- You **_must_** specify `log_group_name` in a `where` clause in order to use this table.
- This table supports optional quals. Queries with optional quals are optimised to used CloudWatch filters. Optional quals are supported for the following columns:
  - `filter`
  - `log_stream_name`
  - `region`
  - `timestamp`

The following tables also retrieve data from CloudWatch log groups, but have columns specific to the log type for easier querying:
- [aws_cloudtrail_trail_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_trail_event)
- [aws_vpc_flow_log_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_flow_log_event)

## Examples

### Basic info

```sql
select
  log_group_name,
  log_stream_name,
  event_id,
  timestamp,
  ingestion_time,
  message
from
  aws_cloudwatch_log_event
where
  log_group_name = '/aws/lambda/myfunction';
```

## Filter Examples

For more information on CloudWatch log filters, please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

### List events for the last 1 hour

```sql
select
  log_group_name,
  log_stream_name,
  event_id,
  timestamp,
  ingestion_time,
  message
from
  aws_cloudwatch_log_event
where
  log_group_name = '/aws/lambda/myfunction' and
  timestamp > ('now'::timestamp - interval '1 hr')
order by
  event_time asc;
```
