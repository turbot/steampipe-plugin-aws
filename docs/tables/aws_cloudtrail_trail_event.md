---
title: "Table: aws_cloudtrail_trail_event - Query AWS CloudTrail Events using SQL"
description: "Allows users to query AWS CloudTrail Events, providing information about each trail event within AWS CloudTrail. The table can be used to retrieve details such as the event time, event name, resources involved, and much more."
---

# Table: aws_cloudtrail_trail_event - Query AWS CloudTrail Events using SQL

The `aws_cloudtrail_trail_event` table in Steampipe provides information about each trail event within AWS CloudTrail. This table allows DevOps engineers to query event-specific details, including event time, event name, resources involved, and more. Users can utilize this table to gather insights on trail events, such as event source, user identity, and request parameters. The schema outlines the various attributes of the CloudTrail event, including the event ID, event version, read only, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudtrail_trail_event` table, you can use the `.inspect aws_cloudtrail_trail_event` command in Steampipe.

**Key columns**:

- `event_time`: This is the timestamp of the event. It can be used to filter events based on the time of occurrence.
- `event_name`: This is the name of the event. It can be used to filter events based on specific actions.
- `event_id`: This is the unique identifier for the event. It can be used to join this table with other tables or to retrieve specific events.

## Examples

### List events that occurred over the last five minutes

```sql
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

### List ordered events that occurred between five to ten minutes ago

```sql
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

### List all action events, i.e., not ReadOnly that occurred over the last hour

```sql
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

### List events for a specific service (IAM) that occurred over the last hour

```sql
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

### List events for an IAM user (steampipe) that occurred over the last hour

```sql
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

### List events performed by IAM users that occurred over the last hour

```sql
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

### List events performed with an assumed role that occurred over the last hour

```sql
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

### List events that were not successfully executed that occurred over the last hour

```sql
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

## Filter examples

For more information on CloudWatch log filters, please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

### List events originating from a specific IP address range that occurred over the last hour

```sql
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
