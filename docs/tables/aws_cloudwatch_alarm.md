# Table: aws_cloudwatch_alarm

A metric alarm watches a single CloudWatch metric or the result of a math expression based on CloudWatch metrics.

## Examples

### Basic info

```sql
select
  name,
  state_value,
  metric_name,
  actions_enabled,
  comparison_operator,
  namespace,
  statistic
from
  aws_cloudwatch_alarm;
```


### List alarms in alarm state

```sql
select
  name,
  alarm_arn,
  state_value,
  state_reason
from
  aws_cloudwatch_alarm
where
 state_value = 'ALARM';
```


### List alarms with alarm actions enabled

```sql
select
  alarm_arn,
  actions_enabled,
  alarm_actions
from
  aws_cloudwatch_alarm
where
  actions_enabled;
```


### Get the metric attached to each alarm based on a single metric

```sql
select
  name,
  metric_name,
  namespace,
  period,
  statistic,
  dimensions
from
  aws_cloudwatch_alarm
where
  metric_name is not null;
```


### Get metrics attached to each alarm based on a metric math expression

```sql
select
  name,
  metric ->> 'Id' as metric_id,
  metric ->> 'Expression' as metric_expression,
  metric -> 'MetricStat' -> 'Metric' ->> 'MetricName' as metric_name,
  metric -> 'MetricStat' -> 'Metric' ->> 'Namespace' as metric_namespace,
  metric -> 'MetricStat' -> 'Metric' ->> 'Dimensions' as metric_dimensions,
  metric ->> 'ReturnData' as metric_return_data
from
  aws_cloudwatch_alarm,
  jsonb_array_elements(metrics) as metric;
```
