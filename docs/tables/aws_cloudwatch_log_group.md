# Table: aws_cloudwatch_log_group

A log group is a group of log streams that share the same retention, monitoring, and access control settings.

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
