# Table: aws_ecs_task

A task is the instantiation of a task definition within a cluster. After you have created a task definition for your application within Amazon ECS, you can specify the number of tasks that will run on your cluster.

## Examples

### Basic info

```sql
select
  cluster_name,
  desired_status,
  launch_type,
  task_arn
from
  aws_ecs_task;
```

### List task attachment details

```sql
select
  cluster_name,
  task_arn,
  a ->> 'Id' as attachment_id,
  a ->> 'Status' as attachment_status,
  a ->> 'Type' as attachment_type,
  jsonb_pretty(a -> 'Details') as attachment_details
from
  aws_ecs_task,
  jsonb_array_elements(attachments) as a;
```

### List task protection details

```sql
select
  cluster_name,
  task_arn,
  protection ->> 'ProtectionEnabled' as protection_enabled,
  protection ->> 'ExpirationDate' as protection_expiration_date
from
  aws_ecs_task
```
