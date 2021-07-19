# Table: aws_cloudtrail_trail_event

AWS CloudTrail is an AWS service that helps you enable governance, compliance, and operational and risk auditing of your AWS account. Actions taken by a user, role, or an AWS service are recorded as events in CloudTrail. Events include actions taken in the AWSManagement Console, AWS Command Line Interface, and AWS SDKs and APIs.

CloudTrail can be configured with CloudWatch Logs to monitor your trail logs.

This table reads cloudtrail events from a cloudwatch Log Group, that is configured to log events from a trail.

**Important Notes:**

- You **_must_** specify `log_group_name` in a `where` clause in order to use this table.
- This table supports optional quals. Queries with optional quals are optimised to used aws cloudwatch filters. Optional quals is supported for below columns:
  - `log_stream_name`
  - `filter`
  - `region`
  - `timestamp`
  - `event_category`
  - `event_id`
  - `aws_region` (i.e region of events, useful in case of multiregion trails)
  - `source_ip_address`
  - `error_code`
  - `event_name`
  - `username`
  - `event_source`
  - `access_key_id`

## Examples

### List all action events (i.e not ReadOnly)

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
  aws_osborn.aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  not read_only
order by
  event_time asc;
```

### List events for a specific service

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
  aws_osborn.aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  event_source = 'iam.amazonaws.com'
order by
  event_time asc;
```

### List events for an IAM user

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
  aws_osborn.aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  username = 'lalit'
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
  aws_osborn.aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  user_type = 'IAMUser'
order by
  event_time asc;
```

### List events performed using assumed role

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
  aws_osborn.aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  user_type = 'AssumedRole'
order by
  event_time asc;
```

### List api request that were not successfully executed.

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
  aws_osborn.aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  error_code is not null
order by
  event_time asc;
```

### Use `filter` qual to search for specific events

Please refer for [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html)

###### Filter events for api requests from a specific ip address

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
  aws_osborn.aws_cloudtrail_trail_event
where
  log_group_name = 'aws-cloudtrail-logs-013122550996-77246e11' and
  filter = '{ $.sourceIPAddress = 203.189.* }'
order by
  event_time asc;
```
