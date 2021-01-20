# Table: aws_cloudwatch_log_metric_filter

Metric filters define the terms and patterns to look for in log data as it is sent to CloudWatch Logs.

## Examples

### Basic AWS cloudwatch log metric info

```sql
select
  name,
  log_group_name,
  creation_time,
  filter_pattern,
  metric_transformation_name,
  metric_transformation_namespace,
  metric_transformation_value
from
  aws_cloudwatch_log_metric_filter;
```


### List the cloudwatch metric filters that sends error logs to cloudwatch log groups

```sql
select
  name,
  log_group_name,
  filter_pattern
from
  aws_cloudwatch_log_metric_filter
where
  filter_pattern ilike '%error%';
```


### Number of metric filters attached to each cloudwatch log group

```sql
select
  log_group_name,
  count(name) as metric_filter_count
from
  aws_cloudwatch_log_metric_filter
group by
  log_group_name;
```