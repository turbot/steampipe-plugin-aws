---
title: "Table: aws_ecs_service - Query AWS Elastic Container Service using SQL"
description: "Allows users to query AWS Elastic Container Service (ECS) to retrieve information about the services within the ECS clusters."
---

# Table: aws_ecs_service - Query AWS Elastic Container Service using SQL

The `aws_ecs_service` table in Steampipe provides information about the services within the AWS Elastic Container Service (ECS) clusters. This table allows DevOps engineers to query service-specific details, including service status, task definitions, and associated metadata. Users can utilize this table to gather insights on services, such as service health status, task definitions being used, and more. The schema outlines the various attributes of the ECS service, including the service ARN, cluster ARN, task definition, desired count, running count, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecs_service` table, you can use the `.inspect aws_ecs_service` command in Steampipe.

**Key columns**:

- `service_arn`: The Amazon Resource Name (ARN) that identifies the service. This column can be used to join with other tables that contain service ARN information.
- `cluster_arn`: The ARN of the cluster where the service is running. This can be used to join with the `aws_ecs_cluster` table.
- `task_definition`: The task definition that the service uses. This can be used to join with the `aws_ecs_task_definition` table to get detailed information about the task definitions being used by the service.

## Examples

### Basic info

```sql
select
  service_name,
  arn,
  cluster_arn,
  task_definition,
  status
from
  aws_ecs_service;
```

### List services not using the latest version of AWS Fargate platform

```sql
select
  service_name,
  arn,
  launch_type,
  platform_version
from
  aws_ecs_service
where
  launch_type = 'FARGATE'
  and platform_version is not null;
```

### List inactive services

```sql
select
  service_name,
  arn,
  status
from
  aws_ecs_service
where
  status = 'INACTIVE';
```
