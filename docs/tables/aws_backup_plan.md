# Table: aws_backup_plan

AWS Backup plan is a policy expression that defines when and how you want to back up your AWS resources, such as Amazon DynamoDB tables or Amazon Elastic File System (Amazon EFS) file systems.
You can assign resources to backup plans, and AWS Backup automatically backs up and retains backups for those resources according to the backup plan. You can create multiple backup plans if you have workloads with different backup requirements.

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

### List plans older than 90 days

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

### List plans order by creation date

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

### List plans which are deleted in last 7 days

```sql
select
  name,
  arn,
  creation_date,
  deletion_date
from
  aws_backup_plan
where
  deletion_date > current_date - 7
order by
  deletion_date;
```
