# Table: aws_backup_selection

AWS Backup selections manage selection conditions for AWS Backup plan resources. A backup plan is a policy expression that defines when and how you want to back up your AWS resources, such as Amazon DynamoDB tables or Amazon Elastic File System (Amazon EFS) file systems.

## Examples

### Basic Info

```sql
select
  selection_name,
  backup_plan_id,
  iam_role_arn,
  region,
  account_id
from
  aws_backup_selection;
```

### List EBS volumes that are in a backup plan

```sql
with filtered_data as (
  select
    backup_plan_id,
    jsonb_agg(r) as assigned_resource
  from
    aws_backup_selection,
    jsonb_array_elements(resources) as r
  group by backup_plan_id
)
select
  v.volume_id,
  v.region,
  v.account_id
from
  aws_ebs_volume as v
  join filtered_data t on t.assigned_resource ?| array[v.arn];
```
