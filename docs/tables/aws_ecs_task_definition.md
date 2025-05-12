---
title: "Steampipe Table: aws_ecs_task_definition - Query AWS ECS Task Definitions using SQL"
description: "Allows users to query AWS ECS Task Definitions to gain insights into the configuration of running tasks in an ECS service. The table provides details such as task definition ARN, family, network mode, revision, status, and more."
folder: "ECS"
---

# Table: aws_ecs_task_definition - Query AWS ECS Task Definitions using SQL

The AWS ECS Task Definition is a blueprint that describes how a Docker container should launch. It specifies the Docker image to use for the container, the required resources, and other configurations. Task Definitions are used in conjunction with the Amazon Elastic Container Service (ECS) to run containers reliably on AWS.

## Table Usage Guide

The `aws_ecs_task_definition` table in Steampipe provides you with information about the task definitions within AWS Elastic Container Service (ECS). This table allows you, as a DevOps engineer, to query task-specific details, including the task definition ARN, family, network mode, revision, and status. You can utilize this table to gather insights on task definitions, such as their configuration, associated IAM roles, container definitions, volumes, and more. The schema outlines the various attributes of the ECS task definition for you, including the task definition ARN, family, requires compatibility, and associated tags.

## Examples

### Basic info
Explore the configuration and status of task definitions in AWS ECS to understand their processing power and network configuration. This can be useful for optimizing resource allocation and network settings for better system performance.

```sql+postgres
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

```sql+sqlite
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
Explore the distribution of containers across various task definitions to better manage and optimize the use of resources in an AWS ECS environment.

```sql+postgres
select
  task_definition_arn,
  jsonb_array_length(container_definitions) as num_of_conatiners
from
  aws_ecs_task_definition;
```

```sql+sqlite
select
  task_definition_arn,
  json_array_length(container_definitions) as num_of_conatiners
from
  aws_ecs_task_definition;
```

### List containers with elevated privileges on the host container instance
Determine the areas in which containers are operating with elevated privileges within your host container instance. This is useful to identify potential security risks and ensure secure configuration of your container infrastructure.

```sql+postgres
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

```sql+sqlite
select
  task_definition_arn,
  json_extract(cd.value, '$.Privileged') as privileged,
  json_extract(cd.value, '$.Name') as container_name
from
  aws_ecs_task_definition,
  json_each(container_definitions) as cd
where
  json_extract(cd.value, '$.Privileged') = 'true';
```


### List task definitions with containers where logging is disabled
This query is useful in identifying all task definitions with containers where logging has been disabled in the AWS ECS system. This can aid in improving security and compliance by enabling you to quickly pinpoint areas where logging should be enabled for better tracking and auditing.

```sql+postgres
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

```sql+sqlite
select
  task_definition_arn,
  json_extract(cd.value, '$.Name') as container_name,
  json_extract(cd.value, '$.LogConfiguration') as log_configuration
from
  aws_ecs_task_definition,
  json_each(container_definitions) as cd
where
 json_extract(cd.value, '$.LogConfiguration') is null;
```