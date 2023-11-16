---
title: "Table: aws_cloudwatch_log_group - Query AWS CloudWatch Log Groups using SQL"
description: "Allows users to query AWS CloudWatch Log Groups and retrieve their attributes such as ARN, creation time, stored bytes, metric filter count, and more."
---

# Table: aws_cloudwatch_log_group - Query AWS CloudWatch Log Groups using SQL

The `aws_cloudwatch_log_group` table in Steampipe provides information about Log Groups within AWS CloudWatch. This table allows DevOps engineers to query Log Group-specific details, including the ARN, creation time, stored bytes, metric filter count, retention period, and associated tags. Users can utilize this table to gather insights on Log Groups, such as their size, age, and associated metrics. The schema outlines the various attributes of the Log Group, including the ARN, creation time, stored bytes, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudwatch_log_group` table, you can use the `.inspect aws_cloudwatch_log_group` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Number (ARN) of the Log Group. This can be used to join this table with other tables that reference AWS resources by their ARN.
- `name`: The name of the Log Group. This can be used to join this table with other tables that reference AWS resources by their name.
- `creation_time`: The creation time of the Log Group. This can be useful for joining with other tables that track resource creation times, to analyze resource age or lifecycle.

## Examples

### List all the log groups that are not encrypted

```sql
select
  name,
  kms_key_id,
  metric_filter_count,
  retention_in_days
from
  aws_cloudwatch_log_group
where
  kms_key_id is null;
```

### List of log groups whose retention period is less than 7 days

```sql
select
  name,
  retention_in_days
from
  aws_cloudwatch_log_group
where
  retention_in_days < 7;
```

### Metric filters info attached log groups

```sql
select
  groups.name as log_group_name,
  metric.name as metric_filter_name,
  metric.filter_pattern,
  metric.metric_transformation_name,
  metric.metric_transformation_value
from
  aws_cloudwatch_log_group groups
  join aws_cloudwatch_log_metric_filter metric on groups.name = metric.log_group_name;
```

### List data protection audit policies and their destinations for each log group

```sql
select
  i as data_identifier,
  s -> 'Operation' -> 'Audit' -> 'FindingsDestination' -> 'S3' -> 'Bucket' as  destination_bucket,
  s -> 'Operation' -> 'Audit' -> 'FindingsDestination' -> 'CloudWatchLogs' -> 'LogGroup'as destination_log_group,
  s -> 'Operation' -> 'Audit' -> 'FindingsDestination' -> 'Firehose' -> 'DeliveryStream'as destination_delivery_stream
from
  aws_cloudwatch_log_group,
  jsonb_array_elements(data_protection_policy -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'DataIdentifier') as i
where
  s ->> 'Sid' = 'audit-policy'
  and name = 'log-group-name'
```

### List log groups with no data protection policy

```sql
select
  arn,
  name,
  creation_time
from
  aws_cloudwatch_log_group
where
  data_protection_policy is null
```
