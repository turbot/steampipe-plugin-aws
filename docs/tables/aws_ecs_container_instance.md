---
title: "Table: aws_ecs_container_instance - Query AWS ECS Container Instance using SQL"
description: "Allows users to query AWS ECS Container Instance to retrieve data about the Amazon Elastic Container Service (ECS) container instances. This includes information about the container instance ARN, status, running tasks count, pending tasks count, agent connected status, and more."
---

# Table: aws_ecs_container_instance - Query AWS ECS Container Instance using SQL

The `aws_ecs_container_instance` table in Steampipe provides information about the Amazon Elastic Container Service (ECS) container instances. This table allows DevOps engineers to query container-specific details, including the container instance ARN, status, running tasks count, pending tasks count, agent connected status, and more. Users can utilize this table to gather insights on container instances, such as the number of running or pending tasks, the status of the agent connection, and more. The schema outlines the various attributes of the ECS container instance, including the instance ARN, instance type, launch type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecs_container_instance` table, you can use the `.inspect aws_ecs_container_instance` command in Steampipe.

### Key columns:

- `container_instance_arn`: The Amazon Resource Name (ARN) that identifies the container instance. This is the primary key of the table and can be used to join with other tables.
- `cluster_arn`: The ARN of the cluster to which the container instance belongs. This can be used to join with the `aws_ecs_cluster` table.
- `status`: The status of the container instance. This can be used to filter container instances based on their status.

## Examples

### Basic info

```sql
select
  arn,
  ec2_instance_id,
  status,
  status_reason,
  running_tasks_count,
  pending_tasks_count
from
  aws_ecs_container_instance;
```


### List container instances that have failed registration

```sql
select
  arn,
  status,
  status_reason
from
  aws_ecs_container_instance
where
  status = 'REGISTRATION_FAILED';
```


### Get details of resources attached to each container instance

```sql
select
  arn,
  attachment ->> 'id' as attachment_id,
  attachment ->> 'status' as attachment_status,
  attachment ->> 'type' as attachment_type
from
  aws_ecs_container_instance,
  jsonb_array_elements(attachments) as attachment;
```


### List container instances with using a given AMI

```sql
select
  arn,
  setting ->> 'Name' as name,
  setting ->> 'Value' as value
from
  aws_ecs_container_instance,
  jsonb_array_elements(attributes) as setting
where
  setting ->> 'Name' = 'ecs.ami-id' and
  setting ->> 'Value' = 'ami-0babb0c4a4e5769b8';
```
