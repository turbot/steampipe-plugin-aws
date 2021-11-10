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
  timestamp asc;
```

## Filter Examples

For more information on CloudWatch log filters, please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

### List events with specific filter pattern based on **eventName**

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
  log_group_name = 'cis/test-log-grp'
  and filter = '{$.eventName="DescribeVpcs"}'
```

### List events with filter pattern based on **errorCode**

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
  log_group_name = 'cis/test-log-grp'
  and filter = '{ ($.errorCode = "*UnauthorizedOperation") || ($.errorCode = "AccessDenied*") }'
```

### List events with specific time frame in hour(s)

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
  log_group_name = 'cis/test-log-grp'
  and filter = '{$.eventName = "AuthorizeSecurityGroupIngress"}'
  and region = 'us-east-1'
  and timestamp >= now() - interval '3 hours'
```

### List events containing either of the **eventName** in the pattern

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
  log_group_name = 'cis/test-log-grp'
  and filter = '{($.eventName = AuthorizeSecurityGroupIngress) || ($.eventName = AuthorizeSecurityGroupEgress) || ($.eventName = RevokeSecurityGroupIngress) || ($.eventName = RevokeSecurityGroupEgress) || ($.eventName = CreateSecurityGroup) || ($.eventName = DeleteSecurityGroup)}'
  and region = 'us-east-1'
  and timestamp >= now() - interval '1 hour'

```

### List events with specific attributes

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
  log_group_name = 'cis/test-log-grp'
  and filter = '{$.userIdentity.sessionContext.sessionIssuer.userName="turbot_superuser"}'
  and timestamp >= now() - interval '1 day'
```