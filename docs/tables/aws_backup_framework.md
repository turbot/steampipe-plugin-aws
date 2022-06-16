# Table: aws_backup_framework

AWS backup framework is a collection of controls that you can use to evaluate your backup practices.
The AWS backup framework will then evaluate whether your backup practices comply with your policies and highlights which
resources are not yet in compliance.

## Examples

### Basic Info

```sql
select
  framework_name,
  arn,
  framework_description,
  deployment_status,
  creation_time,
  number_of_controls,
from
  aws_backup_framework;
```

### List backup frameworks older than 90 days

```sql
select
  framework_name,
  arn,
  framework_description,
  deployment_status,
  creation_time,
  number_of_controls,
from
  aws_backup_plan
where
  creation_date <= (current_date - interval '90' day)
order by
  creation_date;
```
