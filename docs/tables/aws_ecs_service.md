---
title: "Steampipe Table: aws_ecs_service - Query AWS Elastic Container Service using SQL"
description: "Allows users to query AWS Elastic Container Service (ECS) to retrieve information about the services within the ECS clusters."
folder: "ECS"
---

# Table: aws_ecs_service - Query AWS Elastic Container Service using SQL

The AWS Elastic Container Service (ECS) is a highly scalable, high-performance container orchestration service that supports Docker containers and allows you to easily run and scale containerized applications on AWS. ECS eliminates the need for you to install, operate, and scale your own cluster management infrastructure. With simple API calls, you can launch and stop Docker-enabled applications, query the complete state of your application, and access many familiar features like security groups, Elastic Load Balancing, EBS volumes, and IAM roles.

## Table Usage Guide

The `aws_ecs_service` table in Steampipe provides you with information about the services within the AWS Elastic Container Service (ECS) clusters. This table lets you, as a DevOps engineer, query service-specific details, including service status, task definitions, and associated metadata. You can utilize this table to gather insights on services, such as service health status, task definitions being used, and more. The schema outlines the various attributes of the ECS service for you, including the service ARN, cluster ARN, task definition, desired count, running count, and associated tags.

## Examples

### Basic info
Explore the status and details of various tasks within your AWS ECS service. This can help you understand the state of your tasks and identify any potential issues or anomalies.

```sql+postgres
select
  service_name,
  arn,
  cluster_arn,
  task_definition,
  status
from
  aws_ecs_service;
```

```sql+sqlite
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
Determine the areas in which your services are not utilizing the latest version of the AWS Fargate platform. This can be useful in identifying outdated services that may potentially benefit from an upgrade for enhanced performance and security.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are inactive within your AWS ECS services. This can be particularly useful when cleaning up or troubleshooting your environment.

```sql+postgres
select
  service_name,
  arn,
  status
from
  aws_ecs_service
where
  status = 'INACTIVE';
```

```sql+sqlite
select
  service_name,
  arn,
  status
from
  aws_ecs_service
where
  status = 'INACTIVE';
```