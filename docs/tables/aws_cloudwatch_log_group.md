---
title: "Steampipe Table: aws_cloudwatch_log_group - Query AWS CloudWatch Log Groups using SQL"
description: "Allows users to query AWS CloudWatch Log Groups and retrieve their attributes such as ARN, creation time, stored bytes, metric filter count, and more."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_group - Query AWS CloudWatch Log Groups using SQL

The AWS CloudWatch Log Group is a resource that encapsulates your AWS CloudWatch Logs. These log groups are used to monitor, store, and access your log events. It allows you to specify a retention period to automatically expire old log events, thus aiding in managing your log data efficiently.

## Table Usage Guide

The `aws_cloudwatch_log_group` table in Steampipe provides you with information about Log Groups within AWS CloudWatch. This table allows you, as a DevOps engineer, to query Log Group-specific details, including the ARN, creation time, stored bytes, metric filter count, retention period, and associated tags. You can utilize this table to gather insights on Log Groups, such as their size, age, and associated metrics. The schema outlines the various attributes of the Log Group for you, including the ARN, creation time, stored bytes, and associated tags.

## Examples

### List all the log groups that are not encrypted
Identify instances where log groups in AWS CloudWatch are not encrypted. This is beneficial in assessing security measures and ensuring encryption is applied where necessary for data protection.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in your AWS Cloudwatch where log groups are set to retain data for less than a week. This query is useful for identifying potential data loss risks due to short retention periods.

```sql+postgres
select
  name,
  retention_in_days
from
  aws_cloudwatch_log_group
where
  retention_in_days < 7;
```

```sql+sqlite
select
  name,
  retention_in_days
from
  aws_cloudwatch_log_group
where
  retention_in_days < 7;
```

### Metric filters info attached log groups
Uncover the details of how your AWS CloudWatch log groups relate to metric filters, providing a comprehensive view of your logging and monitoring setup. This can be helpful in auditing your CloudWatch configurations, ensuring that important log data is being correctly processed and monitored.

```sql+postgres
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

```sql+sqlite
select
  groups.name as log_group_name,
  metric.name as metric_filter_name,
  metric.filter_pattern,
  metric.metric_transformation_name,
  metric.metric_transformation_value
from
  aws_cloudwatch_log_group as groups
  join aws_cloudwatch_log_metric_filter as metric on groups.name = metric.log_group_name;
```

### List data protection audit policies and their destinations for each log group
Explore the configuration of your data protection audit policies to understand how and where your log data is being sent. This can be useful for ensuring that your logs are being directed to the correct destinations, making it easier to manage and monitor your data.

```sql+postgres
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
  and name = 'log-group-name';
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List log groups with no data protection policy
Determine the areas in which data protection policies are not applied to AWS Cloudwatch log groups. This can be useful for identifying potential security vulnerabilities and ensuring all log data is adequately protected.

```sql+postgres
select
  arn,
  name,
  creation_time
from
  aws_cloudwatch_log_group
where
  data_protection_policy is null;
```

```sql+sqlite
select
  arn,
  name,
  creation_time
from
  aws_cloudwatch_log_group
where
  data_protection_policy is null;
```