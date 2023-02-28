# Table: aws_cloudwatch_log_event

A log group is a group of log streams that share the same retention, monitoring, and access control settings. Log streams contain sequences of log events that share the same source.

This table lists events from a CloudWatch log group.

**Important notes:**

- You **_must_** specify `log_group_name` in a `where` clause in order to use this table.
- For improved performance, it is advised that you use the optional qual `timestamp` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to use CloudWatch filters. Optional quals are supported for the following columns:
  - `filter`
  - `log_stream_name`
  - `region`
  - `timestamp`

The following tables also retrieve data from CloudWatch log groups, but have columns specific to the log type for easier querying:

- [aws_cloudtrail_trail_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_cloudtrail_trail_event)
- [aws_vpc_flow_log_event](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_vpc_flow_log_event)

## Examples

### List events that occurred over the last five minutes

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
  log_group_name = 'cloudwatch-log-event-group-name'
  and timestamp >= now() - interval '5 minutes';
```

### List ordered events that occurred between five to ten minutes ago

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
  log_group_name = 'cloudwatch-log-event-group-name'
  and timestamp between (now() - interval '10 minutes') and (now() - interval '5 minutes')
order by
  timestamp asc;
```

## Filter examples

For more information on CloudWatch log filters, please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

### List events that match the filter pattern term **eventName** to a single value that occurred over the last hour

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
  log_group_name = 'cloudwatch-log-event-group-name'
  and filter = '{$.eventName="DescribeVpcs"}'
  and timestamp >= now() - interval '1 hour';
```

### List events that match the filter pattern term **errorCode** to a single value that occurred over the last hour

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
  log_group_name = 'cloudwatch-log-event-group-name'
  and filter = '{ ($.errorCode = "*UnauthorizedOperation") || ($.errorCode = "AccessDenied*") }'
  and timestamp >= now() - interval '1 hour';
```

### List events that match the filter pattern term **eventName** to multiple values that occurred over the last hour

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
  log_group_name = 'cloudwatch-log-event-group-name'
  and filter = '{($.eventName = AuthorizeSecurityGroupIngress) || ($.eventName = AuthorizeSecurityGroupEgress) || ($.eventName = RevokeSecurityGroupIngress) || ($.eventName = RevokeSecurityGroupEgress) || ($.eventName = CreateSecurityGroup) || ($.eventName = DeleteSecurityGroup)}'
  and region = 'us-east-1'
  and timestamp >= now() - interval '1 hour';
```

### List events which match a specific field in a JSON object that occurred over the past day

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
  log_group_name = 'cloudwatch-log-event-group-name'
  and filter = '{$.userIdentity.sessionContext.sessionIssuer.userName="turbot_superuser"}'
  and timestamp >= now() - interval '1 day';
```
