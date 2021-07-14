# Table: aws_ecs_service

An Amazon ECS service allows you to run and maintain a specified number of instances of a task definition simultaneously in an Amazon ECS cluster. If any of your tasks should fail or stop for any reason, the Amazon ECS service scheduler launches another instance of your task definition to replace it in order to maintain the desired number of tasks in the service.

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
  status <> 'ACTIVE';
```
