---
title: "Steampipe Table: aws_ecs_task - Query AWS ECS Tasks using SQL"
description: "Allows users to query AWS ECS Tasks to obtain detailed information about each task, including its status, task definition, cluster, and other related metadata."
folder: "ECS"
---

# Table: aws_ecs_task - Query AWS ECS Tasks using SQL

AWS Elastic Container Service (ECS) Tasks are a running instance of an Amazon ECS task. They are a scalable unit of computing that contain everything needed to run an application on Amazon ECS. ECS tasks can be used to run applications on a managed cluster of Amazon EC2 instances.

## Table Usage Guide

The `aws_ecs_task` table in Steampipe provides you with information about tasks within Amazon Elastic Container Service (ECS). This table enables you, as a DevOps engineer, to query task-specific details, including the current task status, task definition, associated cluster, and other metadata. You can utilize this table to gather insights on tasks, such as tasks that are running, stopped, or pending, tasks associated with specific clusters, and more. The schema outlines the various attributes of the ECS task for you, including the task ARN, last status, task definition ARN, and associated tags.

## Examples

### Basic info
Determine the status and launch type of tasks within your AWS Elastic Container Service (ECS) to manage and optimize your ECS resources effectively. This can help in maintaining the desired state of your tasks and ensuring they are running as expected.

```sql+postgres
select
  cluster_name,
  desired_status,
  launch_type,
  task_arn
from
  aws_ecs_task;
```

```sql+sqlite
select
  cluster_name,
  desired_status,
  launch_type,
  task_arn
from
  aws_ecs_task;
```

### List task attachment details
This query is useful for gaining insights into the status and types of attachments associated with specific tasks within a cluster. This can help in managing and troubleshooting tasks effectively in a real-world scenario.

```sql+postgres
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

```sql+sqlite
select
  cluster_name,
  task_arn,
  json_extract(a.value, '$.Id') as attachment_id,
  json_extract(a.value, '$.Status') as attachment_status,
  json_extract(a.value, '$.Type') as attachment_type,
  json(a.value, '$.Details') as attachment_details
from
  aws_ecs_task,
  json_each(attachments) as a;
```

### List task protection details
Explore the protection status and expiry dates of tasks within your AWS ECS clusters. This can help ensure all tasks are adequately protected and any expiring protections are promptly renewed.

```sql+postgres
select
  cluster_name,
  task_arn,
  protection ->> 'ProtectionEnabled' as protection_enabled,
  protection ->> 'ExpirationDate' as protection_expiration_date
from
  aws_ecs_task;
```

```sql+sqlite
select
  cluster_name,
  task_arn,
  json_extract(protection, '$.ProtectionEnabled') as protection_enabled,
  json_extract(protection, '$.ExpirationDate') as protection_expiration_date
from
  aws_ecs_task;
```