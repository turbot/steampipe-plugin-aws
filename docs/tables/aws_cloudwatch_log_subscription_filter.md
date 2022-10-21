# Table: aws_cloudwatch_log_subscription_filter

A subscription filter defines the filter pattern to use for filtering which log events get delivered to your AWS resource, as well as information about where to send matching log events to.

## Examples

### Basic info

```sql
select
  name,
  log_group_name,
  creation_time,
  filter_pattern,
  destination_arn
from
  aws_cloudwatch_log_subscription_filter;
```

### List the cloudwatch subscription filters that sends error logs to cloudwatch log groups

```sql
select
  name,
  log_group_name,
  filter_pattern
from
  aws_cloudwatch_log_subscription_filter
where
  filter_pattern ilike '%error%';
```

### Number of subscription filters attached to each cloudwatch log group

```sql
select
  log_group_name,
  count(name) as subscription_filter_count
from
  aws_cloudwatch_log_subscription_filter
group by
  log_group_name;
```