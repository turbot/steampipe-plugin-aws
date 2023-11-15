# Table: aws_directory_service_log_subscription

This forward real-time Directory Service domain controller security logs to the specified Amazon CloudWatch log group in your AWS account.

## Examples

### Basic info

```sql
select
  log_group_name,
  partition,
  subscription_created_date_time,
  directory_id,
  title
from
  aws_directory_service_log_subscription;
```

### Get details of the directory associated to the log subscription

```sql
select
  s.log_group_name,
  d.name as directory_name,
  d.arn as directory_arn,
  d.directory_id,
  d.type as directory_type
from
  aws_directory_service_log_subscription as s
  left join aws_directory_service_directory as d on s.directory_id = d.directory_id;
```
