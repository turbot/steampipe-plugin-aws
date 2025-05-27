---
title: "Steampipe Table: aws_cloudwatch_log_metric_filter - Query AWS CloudWatch log metric filters using SQL"
description: "Allows users to query AWS CloudWatch log metric filters to obtain detailed information about each filter, including its name, creation date, associated log group, filter pattern, metric transformations and more."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_metric_filter - Query AWS CloudWatch log metric filters using SQL

The AWS CloudWatch Log Metric Filter is a feature within AWS CloudWatch that enables you to extract information from the logs and create custom metrics. These custom metrics can be used for detailed monitoring and alarming based on patterns that might appear in your logs. This is a powerful tool for identifying trends, troubleshooting issues, and setting up real-time monitoring across your AWS resources.

## Table Usage Guide

The `aws_cloudwatch_log_metric_filter` table in Steampipe provides you with information about log metric filters within AWS CloudWatch. This table allows you, as a DevOps engineer, to query filter-specific details, including the associated log group, filter pattern, and metric transformations. You can utilize this table to gather insights on filters, such as filter patterns used, metrics generated from log data, and more. The schema outlines for you the various attributes of the log metric filter, including the filter name, creation date, filter pattern, and associated log group.

## Examples

### Basic AWS cloudwatch log metric info
Explore the essential characteristics and setup of your AWS CloudWatch log metrics. This query can help you assess the overall configuration and performance metrics of your logs, providing valuable insights for monitoring and optimizing your AWS environment.

```sql+postgres
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

```sql+sqlite
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
Identify instances where specific metric filters are configured to send error logs to Cloudwatch log groups. This allows for effective error tracking and proactive issue resolution in cloud environments.

```sql+postgres
select
  name,
  log_group_name,
  filter_pattern
from
  aws_cloudwatch_log_metric_filter
where
  filter_pattern ilike '%error%';
```

```sql+sqlite
select
  name,
  log_group_name,
  filter_pattern
from
  aws_cloudwatch_log_metric_filter
where
  filter_pattern like '%error%';
```


### Number of metric filters attached to each cloudwatch log group
Determine the areas in which Cloudwatch log groups have multiple metric filters attached. This can help in managing and optimizing your AWS Cloudwatch setup by understanding the distribution of metric filters across different log groups.

```sql+postgres
select
  log_group_name,
  count(name) as metric_filter_count
from
  aws_cloudwatch_log_metric_filter
group by
  log_group_name;
```

```sql+sqlite
select
  log_group_name,
  count(name) as metric_filter_count
from
  aws_cloudwatch_log_metric_filter
group by
  log_group_name;
```