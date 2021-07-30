# Table: aws_cloudtrail_trail_event

AWS CloudTrail is an AWS service that helps you enable governance, compliance, and operational and risk auditing of your AWS account. Actions taken by a user, role, or an AWS service are recorded as events in CloudTrail. These events can be sent to a CloudWatch log group to allow for easy monitoring.

This table reads CloudTrail event data from a CloudWatch log group that is configured to log events from a trail.

**Important notes:**

- You **_must_** specify `log_group_name` in a `where` clause in order to use this table.
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

### List all action events, i.e., not ReadOnly

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
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  not read_only
order by
  event_time asc;
```

### List events for a specific service (IAM)

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
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  event_source = 'iam.amazonaws.com'
order by
  event_time asc;
```

### List events for an IAM user (steampipe)

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
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  username = 'steampipe'
order by
  event_time asc;
```

### List events performed by IAM users

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
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  user_type = 'IAMUser'
order by
  event_time asc;
```

### List events performed with an assumed role

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
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  user_type = 'AssumedRole'
order by
  event_time asc;
```

### List events that were not successfully executed

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
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  error_code is not null
order by
  event_time asc;
```

### Use `filter` qual to search for specific events

Please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html) for filter examples.

#### Filter events originating from a specific IP address range

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
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  filter = '{ $.sourceIPAddress = 203.189.* }'
order by
  event_time asc;
```
