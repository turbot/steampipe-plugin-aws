# Table: aws_cloudtrail_trail_event

AWS CloudTrail is an AWS service that helps you enable governance, compliance, and operational and risk auditing of your AWS account. Actions taken by a user, role, or an AWS service are recorded as events in CloudTrail. These events can be sent to a CloudWatch log group to allow for easy monitoring.

This table reads CloudTrail event data from a CloudWatch log group that is configured to log events from a trail.

**Important notes:**

- You **_must_** specify `log_group_name` in a `where` clause in order to use this table.
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
