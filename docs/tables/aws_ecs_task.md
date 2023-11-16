---
title: "Table: aws_ecs_task - Query AWS ECS Tasks using SQL"
description: "Allows users to query AWS ECS Tasks to obtain detailed information about each task, including its status, task definition, cluster, and other related metadata."
---

# Table: aws_ecs_task - Query AWS ECS Tasks using SQL

The `aws_ecs_task` table in Steampipe provides information about tasks within Amazon Elastic Container Service (ECS). This table allows DevOps engineers to query task-specific details, including the current task status, task definition, associated cluster, and other metadata. Users can utilize this table to gather insights on tasks, such as tasks that are running, stopped, or pending, tasks associated with specific clusters, and more. The schema outlines the various attributes of the ECS task, including the task ARN, last status, task definition ARN, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecs_task` table, you can use the `.inspect aws_ecs_task` command in Steampipe.

### Key columns:

- `task_arn`: The Amazon Resource Name (ARN) that identifies the task. This is a key column as it uniquely identifies each task and can be used to join this table with other tables that contain task-specific information.
- `cluster_arn`: The ARN of the cluster that hosts the task. This column is important for identifying which tasks belong to which clusters and can be used to join this table with the `aws_ecs_cluster` table.
- `task_definition_arn`: The ARN of the task definition for the task. This column is useful for identifying the task definition associated with each task and can be used to join this table with the `aws_ecs_task_definition` table.

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
  aws_ecs_task;
```
