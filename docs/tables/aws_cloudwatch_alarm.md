# Table: aws_cloudwatch_alarm

A metric alarm watches a single CloudWatch metric or the result of a math expression based on CloudWatch metrics.

## Examples

### Basic info

```sql
select
  name,
  state_value,
  metric_name,
  actions_enabled
from
  aws_cloudwatch_alarm;
```

### List of cloudwatch alarms whose state is in alarm

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

### List of cloudwatch alarms whose action enabled is on

```sql
select
  alarm_arn,
  actions_enabled
from
  aws_cloudwatch_alarm
where
  actions_enabled;
```


### Metric attached to cloudwatch alarm based on a single metric

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


### List of metrics attached to cloudwatch alarm based on a metric math expression

```sql
select
  name,
  metric ->> 'Id' as metric_id,
  metric -> 'MetricStat' -> 'Metric' ->> 'MetricName' as metric_name,
  metric -> 'MetricStat' -> 'Metric' ->> 'Namespace' as namespace,
  metric -> 'MetricStat' ->> 'Period' as period,
  metric ->> 'ReturnData' as metric_return_data
from
  aws_cloudwatch_alarm,
  jsonb_array_elements(metrics) as metric;
```