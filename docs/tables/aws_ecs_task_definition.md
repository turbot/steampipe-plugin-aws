---
title: "Table: aws_ecs_task_definition - Query AWS ECS Task Definitions using SQL"
description: "Allows users to query AWS ECS Task Definitions to gain insights into the configuration of running tasks in an ECS service. The table provides details such as task definition ARN, family, network mode, revision, status, and more."
---

# Table: aws_ecs_task_definition - Query AWS ECS Task Definitions using SQL

The `aws_ecs_task_definition` table in Steampipe provides information about the task definitions within AWS Elastic Container Service (ECS). This table allows DevOps engineers to query task-specific details, including the task definition ARN, family, network mode, revision, and status. Users can utilize this table to gather insights on task definitions, such as their configuration, associated IAM roles, container definitions, volumes, and more. The schema outlines the various attributes of the ECS task definition, including the task definition ARN, family, requires compatibility, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecs_task_definition` table, you can use the `.inspect aws_ecs_task_definition` command in Steampipe.

### Key columns:

- `task_definition_arn`: The Amazon Resource Name (ARN) that identifies the task definition. This can be used to join with other tables that reference task definitions.
- `family`: The family of the task definition. This allows for easy grouping and querying of related task definitions.
- `network_mode`: The Docker networking mode to use for the containers in the task. This can provide insight into the networking configuration of the task.

## Examples

### Basic info

```sql
select
  task_definition_arn,
  cpu,
  network_mode,
  title,
  status,
  tags
from
  aws_ecs_task_definition;
```


### Count the number of containers attached to each task definitions

```sql
select
  task_definition_arn,
  jsonb_array_length(container_definitions) as num_of_conatiners
from
  aws_ecs_task_definition;
```


### List containers with elevated privileges on the host container instance

```sql
select
  task_definition_arn,
  cd ->> 'Privileged' as privileged,
  cd ->> 'Name' as container_name
from
  aws_ecs_task_definition,
  jsonb_array_elements(container_definitions) as cd
where
  cd ->> 'Privileged' = 'true';
```


### List task definitions with containers where logging is disabled

```sql
select
  task_definition_arn,
  cd ->> 'Name' as container_name,
  cd ->> 'LogConfiguration' as log_configuration
from
  aws_ecs_task_definition,
  jsonb_array_elements(container_definitions) as cd
where
 cd ->> 'LogConfiguration' is null;
```
