---
title: "Steampipe Table: aws_cloudwatch_alarm - Query AWS CloudWatch Alarms using SQL"
description: "Allows users to query AWS CloudWatch Alarms, providing detailed information about each alarm, including its configuration, state, and associated actions."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_alarm - Query AWS CloudWatch Alarms using SQL

The AWS CloudWatch Alarms is a feature of Amazon CloudWatch, a monitoring service for AWS resources and applications. CloudWatch Alarms allow you to monitor Amazon Web Services resources and trigger actions when changes in data points meet certain defined thresholds. They help you react quickly to issues that may affect your applications or infrastructure, thereby enhancing your ability to keep applications running smoothly.

## Table Usage Guide

The `aws_cloudwatch_alarm` table in Steampipe provides you with information about alarms within AWS CloudWatch. This table allows you, as a DevOps engineer, to query alarm-specific details, including its current state, configuration, and actions associated with each alarm. You can utilize this table to gather insights on alarms, such as alarms in a particular state, alarms associated with specific AWS resources, and understanding the actions that will be triggered when an alarm state changes. The schema outlines the various attributes of the CloudWatch alarm for you, including the alarm name, alarm description, metric name, comparison operator, and associated tags.

## Examples

### Basic info
Explore the status and configurations of your CloudWatch alarms to understand their current operational state and the conditions that trigger them. This can help you monitor the health and performance of your AWS resources more effectively.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are currently in an alarm state. This is useful to quickly identify and address any issues within your cloud infrastructure.

```sql+postgres
select
  name,
  arn,
  state_value,
  state_reason
from
  aws_cloudwatch_alarm
where
 state_value = 'ALARM';
```

```sql+sqlite
select
  name,
  arn,
  state_value,
  state_reason
from
  aws_cloudwatch_alarm
where
 state_value = 'ALARM';
```


### List alarms with alarm actions enabled
Identify instances where alarms have been activated with specific actions in the AWS CloudWatch service. This can be useful in understanding the active monitoring and alerting mechanisms in place for system events.

```sql+postgres
select
  arn,
  actions_enabled,
  alarm_actions
from
  aws_cloudwatch_alarm
where
  actions_enabled;
```

```sql+sqlite
select
  arn,
  actions_enabled,
  alarm_actions
from
  aws_cloudwatch_alarm
where
  actions_enabled = 1;
```


### Get the metric attached to each alarm based on a single metric
Discover the segments that have alarms set based on specific metrics within the AWS Cloudwatch service. This is particularly useful for monitoring and managing application performance, resource utilization, and operational health.

```sql+postgres
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

```sql+sqlite
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
Identify the metrics associated with each alarm based on mathematical expressions. This can help in understanding the performance of various elements and aid in proactive monitoring and troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(metric, '$.Id') as metric_id,
  json_extract(metric, '$.Expression') as metric_expression,
  json_extract(metric, '$.MetricStat.Metric.MetricName') as metric_name,
  json_extract(metric, '$.MetricStat.Metric.Namespace') as metric_namespace,
  json_extract(metric, '$.MetricStat.Metric.Dimensions') as metric_dimensions,
  json_extract(metric, '$.ReturnData') as metric_return_data
from
  aws_cloudwatch_alarm,
  json_each(metrics) as metric;
```