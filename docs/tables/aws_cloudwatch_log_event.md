---
title: "Steampipe Table: aws_cloudwatch_log_event - Query AWS CloudWatch Log Events using SQL"
description: "Allows users to query AWS CloudWatch Log Events to retrieve information about log events from a specified log group. Users can utilize this table to monitor and troubleshoot systems and applications using their existing log data."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_event - Query AWS CloudWatch Log Events using SQL

The AWS CloudWatch Log Events is a feature of Amazon CloudWatch that enables you to monitor, store, and access your log files from Amazon Elastic Compute Cloud (EC2) instances, AWS CloudTrail, and other sources. It allows you to centralize the logs from all your systems, applications, and AWS services that you use, in a single, highly scalable service. With CloudWatch Log Events, you can quickly search and filter your log data for specific error codes or patterns, and set alarms for specific phrases, values or patterns that appear in your log data.

## Table Usage Guide

The `aws_cloudwatch_log_event` table in Steampipe provides you with information about Log Events within AWS CloudWatch. This table allows you, as a DevOps engineer, system administrator, or developer, to query event-specific details, including the event message, event timestamp, and associated metadata. You can utilize this table to gather insights on log events, such as event patterns, event frequency, event sources, and more. The schema outlines the various attributes of the Log Event for you, including the event ID, log group name, log stream name, and ingestion time.

**Important Notes**
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
Explore recent activity within your system by identifying events that have occurred in the past five minutes. This is particularly useful for real-time monitoring and immediate issue detection.

```sql+postgres
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

```sql+sqlite
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
  and timestamp >= datetime('now', '-5 minutes');
```

### List ordered events that occurred between five to ten minutes ago
Determine the sequence of events that transpired within a specific timeframe in your AWS CloudWatch logs. This is useful for tracking activity and identifying potential issues that occurred between five to ten minutes ago.

```sql+postgres
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

```sql+sqlite
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
  and timestamp between (strftime('%s','now') - 600) and (strftime('%s','now') - 300)
order by
  timestamp asc;
```

## Filter examples

For more information on CloudWatch log filters, please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

### List events that match the filter pattern term **eventName** to a single value that occurred over the last hour
Determine the occurrences of a specific event within the last hour in your AWS CloudWatch logs. This is particularly useful for tracking and analyzing specific activities or changes over a short period of time.

```sql+postgres
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

```sql+sqlite
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
  and json_extract(filter, '$.eventName') = "DescribeVpcs"
  and timestamp >= datetime('now', '-1 hour');
```

### List events that match the filter pattern term **errorCode** to a single value that occurred over the last hour
The query is designed to monitor and identify instances of unauthorized access or access denial within the last hour. This is particularly useful for maintaining security and troubleshooting access issues in real-time.

```sql+postgres
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

```sql+sqlite
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
  and (json_extract(filter, '$.errorCode') = "*UnauthorizedOperation" or json_extract(filter, '$.errorCode') = "AccessDenied*")
  and timestamp >= datetime('now', '-1 hours');
```

### List events that match the filter pattern term **eventName** to multiple values that occurred over the last hour
Explore the specific security-related events in your AWS CloudWatch logs from the past hour to gain insights into potential security changes or threats. This helps in maintaining a secure and compliant environment by tracking changes in security groups and identifying suspicious activities.

```sql+postgres
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

```sql+sqlite
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
  and json_extract(filter, '$.eventName') in ('AuthorizeSecurityGroupIngress', 'AuthorizeSecurityGroupEgress', 'RevokeSecurityGroupIngress', 'RevokeSecurityGroupEgress', 'CreateSecurityGroup', 'DeleteSecurityGroup')
  and region = 'us-east-1'
  and timestamp >= datetime('now', '-1 hour');
```

### List events which match a specific field in a JSON object that occurred over the past day
This query is useful for monitoring user activity within a specific time frame. Specifically, it helps identify actions taken by a 'superuser' within the last day, providing insights into their behavior and potential security implications.

```sql+postgres
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

```sql+sqlite
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
  and json_extract(filter, '$.userIdentity.sessionContext.sessionIssuer.userName')="turbot_superuser"
  and timestamp >= datetime('now', '-1 day');
```