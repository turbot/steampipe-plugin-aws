# Table: aws_cloudwatch_log_group

A log group is a group of log streams that share the same retention, monitoring, and access control settings.

**Note**: We recommend you specify a `name_prefix` to speed up queries, especially in regions that have a large number of them.

## Examples

### Basic info

```sql
select
  name,
  kms_key_id,
  metric_filter_count,
  retention_in_days
from
  aws_cloudwatch_log_group
where
  name_prefix = '/aws/lambda/';
```

### List log groups that are not encrypted

```sql
select
  name,
  kms_key_id,
  metric_filter_count,
  retention_in_days
from
  aws_cloudwatch_log_group
where
  kms_key_id is null
  and name_prefix = '/aws/lambda/';
```

### List log groups whose retention period is less than 7 days

```sql
select
  name,
  retention_in_days
from
  aws_cloudwatch_log_group
where
  retention_in_days < 7
  and name_prefix = '/aws/lambda/';
```

### Get metric filter info for each log group

```sql
select
  g.name as log_group_name,
  m.name as metric_filter_name,
  m.filter_pattern,
  m.metric_transformation_name,
  m.metric_transformation_value
from
  aws_cloudwatch_log_group g
  join aws_cloudwatch_log_metric_filter m on g.name = m.log_group_name
where
  g.name_prefix = '/aws/lambda/';
```
