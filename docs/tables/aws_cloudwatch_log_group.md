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

### List of cloudwatch log groups whose retention period is less than 7 days

```sql
select
  name,
  retention_in_days
from
  aws_cloudwatch_log_group
where
  retention_in_days < 7;
```

### Metric filters info attached to cloudwatch log groups

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

### List the cloudwatch log groups' data protection policy documents

```sql
select
  name,
  jsonb_pretty(data_protection_policy_document) as policy_document
from
  aws_cloudwatch_log_group;
```
