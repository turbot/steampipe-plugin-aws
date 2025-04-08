---
title: "Steampipe Table: aws_cloudtrail_trail_event - Query AWS CloudTrail Events using SQL"
description: "Allows users to query AWS CloudTrail Events, providing information about each trail event within AWS CloudTrail. The table can be used to retrieve details such as the event time, event name, resources involved, and much more."
folder: "CloudTrail"
---

# Table: aws_cloudtrail_trail_event - Query AWS CloudTrail Events using SQL

AWS CloudTrail Events are records of activity within your AWS environment. This service captures all API calls for your account, including calls made via the AWS Management Console, SDKs, and CLI. It provides a history of AWS API calls for your account, including API calls made via the AWS Management Console, AWS SDKs, command line tools, and higher-level AWS services.

## Table Usage Guide

The `aws_cloudtrail_trail_event` table in Steampipe provides you with information about each trail event within AWS CloudTrail. This table allows you, as a DevOps engineer, to query event-specific details, including event time, event name, resources involved, and more. You can utilize this table to gather insights on trail events, such as event source, user identity, and request parameters. The schema outlines the various attributes of the CloudTrail event for you, including the event ID, event version, read only, and associated tags.

**Important Notes**
- You must specify `log_group_name` in a `where` clause in order to use this table.
- For improved performance, it is advised that you use the optional qual `timestamp` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to use CloudWatch filters. Optional quals are supported for the following columns:
  - `access_key_id`
  - `aws_region` (region of the event, useful in case of multi-region trails)
  - `error_code`
  - `event_category`
  - `event_id`
  - `event_name`
  - `event_source`
  - `filter`
  - `log_stream_name`
  - `region`
  - `source_ip_address`
  - `timestamp`
  - `username`

## Examples

### List events that occurred over the last five minutes
This query is useful for gaining insights into recent activity within your AWS environment. It provides a quick overview of the events that have taken place in the last five minutes, which can be particularly useful for immediate incident response or real-time monitoring.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  jsonb_pretty(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and timestamp >= now() - interval '5 minutes';
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  json(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and timestamp >= datetime('now', '-5 minutes');
```

### List ordered events that occurred between five to ten minutes ago
Explore the sequence of events that occurred within a specific time frame in the recent past. This can be useful for auditing activities, identifying anomalies, or tracking user behaviour within a given period.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  jsonb_pretty(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and timestamp between (now() - interval '10 minutes') and (now() - interval '5 minutes')
order by
  event_time asc;
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  json(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and timestamp between (datetime('now', '-10 minutes')) and (datetime('now', '-5 minutes'))
order by
  event_time asc;
```

### List all action events, i.e., not ReadOnly that occurred over the last hour
Explore which action events have occurred in the last hour on AWS Cloudtrail. This is useful for identifying recent activities that have potentially altered your system.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  jsonb_pretty(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and not read_only
  and timestamp >= now() - interval '1 hour'
order by
  event_time asc;
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  json(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and not read_only
  and timestamp >= datetime('now', '-1 hours')
order by
  event_time asc;
```

### List events for a specific service (IAM) that occurred over the last hour
This query allows users to monitor recent activity for a specific service, in this case, AWS's Identity and Access Management (IAM). It is particularly useful for security audits, as it provides a chronological overview of events, including who initiated them and what actions were taken, over the last hour.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  user_type,
  user_identifier,
  jsonb_pretty(request_parameters) as request_parameters,
  jsonb_pretty(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and event_source = 'iam.amazonaws.com'
  and timestamp >= now() - interval '1 hour'
order by
  event_time asc;
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  user_type,
  user_identifier,
  json(request_parameters) as request_parameters,
  json(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and event_source = 'iam.amazonaws.com'
  and timestamp >= datetime('now', '-1 hour')
order by
  event_time asc;
```

### List events for an IAM user (steampipe) that occurred over the last hour
Explore which events have occurred on your system over the past hour that are associated with a specific IAM user. This can help in monitoring user activity and identifying potential security concerns.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  jsonb_pretty(request_parameters) as request_parameters,
  jsonb_pretty(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and username = 'steampipe'
  and timestamp >= now() - interval '1 hour'
order by
  event_time asc;
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  request_parameters,
  response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and username = 'steampipe'
  and timestamp >= datetime('now', '-1 hour')
order by
  event_time asc;
```

### List events performed by IAM users that occurred over the last hour
Determine the activities undertaken by IAM users within the past hour in your AWS environment. This can help in understanding user behaviors, monitoring security, and auditing purposes.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  jsonb_pretty(request_parameters) as request_parameters,
  jsonb_pretty(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and user_type = 'IAMUser'
  and timestamp >= now() - interval '1 hour'
order by
  event_time asc;
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  request_parameters,
  response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and user_type = 'IAMUser'
  and timestamp >= datetime('now','-1 hours')
order by
  event_time asc;
```

### List events performed with an assumed role that occurred over the last hour
Explore which actions were carried out using an assumed role in the past hour. This is useful in monitoring and auditing for any unusual or unauthorized activities.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  jsonb_pretty(request_parameters) as request_parameters,
  jsonb_pretty(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and user_type = 'AssumedRole'
  and timestamp >= now() - interval '1 hour'
order by
  event_time asc;
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  user_type,
  username,
  user_identifier,
  request_parameters,
  response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and user_type = 'AssumedRole'
  and timestamp >= datetime('now', '-1 hours')
order by
  event_time asc;
```

### List events that were not successfully executed that occurred over the last hour
Identify instances where events were not executed successfully in the past hour. This is useful for monitoring system performance and quickly addressing any operational issues.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  error_code,
  error_message,
  user_type,
  username,
  user_identifier,
  jsonb_pretty(request_parameters) as request_parameters,
  jsonb_pretty(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and error_code is not null
  and timestamp >= now() - interval '1 hour'
order by
  event_time asc;
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  error_code,
  error_message,
  user_type,
  username,
  user_identifier,
  request_parameters,
  response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and error_code is not null
  and timestamp >= datetime('now','-1 hours')
order by
  event_time asc;
```

## Filter examples

For more information on CloudWatch log filters, please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

### List events originating from a specific IP address range that occurred over the last hour
Explore which events have originated from a specific IP address range in the last hour. This is useful for understanding and monitoring recent activity and potential security incidents related to that IP range.

```sql+postgres
select
  event_name,
  event_source,
  event_time,
  error_code,
  error_message,
  user_type,
  username,
  user_identifier,
  jsonb_pretty(request_parameters) as request_parameters,
  jsonb_pretty(response_elements) as response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and filter = '{ $.sourceIPAddress = 203.189.* }'
  and timestamp >= now() - interval '1 hour'
order by
  event_time asc;
```

```sql+sqlite
select
  event_name,
  event_source,
  event_time,
  error_code,
  error_message,
  user_type,
  username,
  user_identifier,
  request_parameters,
  response_elements
from
  aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-log-group-name'
  and json_extract(filter, '$.sourceIPAddress') like '203.189.%'
  and timestamp >= datetime('now', '-1 hour')
order by
  event_time asc;
```

