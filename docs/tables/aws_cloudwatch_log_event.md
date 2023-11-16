---
title: "Table: aws_cloudwatch_log_event - Query AWS CloudWatch Log Events using SQL"
description: "Allows users to query AWS CloudWatch Log Events to retrieve information about log events from a specified log group. Users can utilize this table to monitor and troubleshoot systems and applications using their existing log data."
---

# Table: aws_cloudwatch_log_event - Query AWS CloudWatch Log Events using SQL

The `aws_cloudwatch_log_event` table in Steampipe provides information about Log Events within AWS CloudWatch. This table allows DevOps engineers, system administrators, and developers to query event-specific details, including the event message, event timestamp, and associated metadata. Users can utilize this table to gather insights on log events, such as event patterns, event frequency, event sources, and more. The schema outlines the various attributes of the Log Event, including the event ID, log group name, log stream name, and ingestion time.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudwatch_log_event` table, you can use the `.inspect aws_cloudwatch_log_event` command in Steampipe.

Key columns:

- `log_group_name`: This is the name of the log group where the log event was created. This can be used to join with other tables that contain log group information.
- `log_stream_name`: This is the name of the log stream where the log event was created. This can be used to join with other tables that contain log stream information.
- `event_id`: This is the unique identifier of the log event. This can be used to join with other tables that contain event-specific information.

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
