---
title: "Table: aws_cloudwatch_alarm - Query AWS CloudWatch Alarms using SQL"
description: "Allows users to query AWS CloudWatch Alarms, providing detailed information about each alarm, including its configuration, state, and associated actions."
---

# Table: aws_cloudwatch_alarm - Query AWS CloudWatch Alarms using SQL

The `aws_cloudwatch_alarm` table in Steampipe provides information about alarms within AWS CloudWatch. This table allows DevOps engineers to query alarm-specific details, including its current state, configuration, and actions associated with each alarm. Users can utilize this table to gather insights on alarms, such as alarms in a particular state, alarms associated with specific AWS resources, and understanding the actions that will be triggered when an alarm state changes. The schema outlines the various attributes of the CloudWatch alarm, including the alarm name, alarm description, metric name, comparison operator, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudwatch_alarm` table, you can use the `.inspect aws_cloudwatch_alarm` command in Steampipe.

### Key columns:

- `alarm_name`: The name of the alarm. This is the primary key column and can be used to join with other tables.
- `state_value`: The state value for the alarm. This can be useful for filtering alarms based on their current state.
- `actions_enabled`: Indicates whether actions should be executed during any changes to the alarm state. This can be useful for understanding the potential impact of an alarm state change.

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
  arn,
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
  arn,
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
