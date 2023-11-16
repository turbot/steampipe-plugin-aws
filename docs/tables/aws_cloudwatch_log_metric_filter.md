---
title: "Table: aws_cloudwatch_log_metric_filter - Query AWS CloudWatch log metric filters using SQL"
description: "Allows users to query AWS CloudWatch log metric filters to obtain detailed information about each filter, including its name, creation date, associated log group, filter pattern, metric transformations and more."
---

# Table: aws_cloudwatch_log_metric_filter - Query AWS CloudWatch log metric filters using SQL

The `aws_cloudwatch_log_metric_filter` table in Steampipe provides information about log metric filters within AWS CloudWatch. This table allows DevOps engineers to query filter-specific details, including the associated log group, filter pattern, and metric transformations. Users can utilize this table to gather insights on filters, such as filter patterns used, metrics generated from log data, and more. The schema outlines the various attributes of the log metric filter, including the filter name, creation date, filter pattern, and associated log group.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudwatch_log_metric_filter` table, you can use the `.inspect aws_cloudwatch_log_metric_filter` command in Steampipe.

### Key columns:

- `name`: The name of the metric filter. This can be used to join with other tables that need filter-specific information.
- `log_group_name`: The name of the log group to which the metric filter is associated. This can be used to join with other tables that need log group-specific information.
- `metric_transformation`: The metric transformation operations associated with the filter. This can be used to join with other tables that need information about the metrics generated from log data.

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