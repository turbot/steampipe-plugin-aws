# Table: aws_backup_plan

In AWS Backup, a backup plan is a policy expression that defines when and how you want to back up your AWS resources, such as Amazon DynamoDB tables or Amazon Elastic File System (Amazon EFS) file systems. You can assign resources to backup plans, and AWS Backup automatically backs up and retains backups for those resources according to the backup plan. You can create multiple backup plans if you have workloads with different backup requirements.

## Examples

### Basic Info

```sql
select
  name,
  backup_plan_id,
  arn,
  creation_date,
  last_execution_date
from
  aws_backup_plan;
```


### List of backup_plans older than 90 days

```sql
select
  name,
  backup_plan_id,
  arn,
  creation_date,
  last_execution_date
from
  aws_backup_plan
where
  creation_date <= (current_date - interval '90' day)
  order by
  creation_date;
```

### List of backup_plans order by creation date

```sql
select
  name,
  backup_plan,
  creation_date
from
  aws_backup_plan
  order by
  creation_date;
```

