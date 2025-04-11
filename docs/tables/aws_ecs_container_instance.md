---
title: "Steampipe Table: aws_ecs_container_instance - Query AWS ECS Container Instance using SQL"
description: "Allows users to query AWS ECS Container Instance to retrieve data about the Amazon Elastic Container Service (ECS) container instances. This includes information about the container instance ARN, status, running tasks count, pending tasks count, agent connected status, and more."
folder: "ECS"
---

# Table: aws_ecs_container_instance - Query AWS ECS Container Instance using SQL

The AWS ECS Container Instance is a resource within the Amazon Elastic Container Service (ECS). It refers to a single EC2 instance that is part of an ECS cluster, which runs containerized applications. It provides the necessary infrastructure to manage, schedule, and run Docker containers on a cluster.

## Table Usage Guide

The `aws_ecs_container_instance` table in Steampipe provides you with information about the Amazon Elastic Container Service (ECS) container instances. This table allows you, as a DevOps engineer, to query container-specific details, including the container instance ARN, status, running tasks count, pending tasks count, agent connected status, and more. You can utilize this table to gather insights on container instances, such as the number of running or pending tasks, the status of the agent connection, and more. The schema outlines the various attributes of the ECS container instance for you, including the instance ARN, instance type, launch type, and associated tags.

## Examples

### Basic info
Determine the areas in which your AWS Elastic Container Service (ECS) instances are running and their status. This can help you identify instances where there are pending tasks, providing insights for potential performance improvement.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which container instances have failed to register within the AWS ECS service. This is useful in diagnosing and resolving issues that could potentially disrupt your application's performance or availability.

```sql+postgres
select
  arn,
  status,
  status_reason
from
  aws_ecs_container_instance
where
  status = 'REGISTRATION_FAILED';
```

```sql+sqlite
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
Explore which resources are linked to each container instance in AWS ECS to better manage and track resources. This can help in identifying any potential issues or inefficiencies in resource allocation.

```sql+postgres
select
  arn,
  attachment ->> 'id' as attachment_id,
  attachment ->> 'status' as attachment_status,
  attachment ->> 'type' as attachment_type
from
  aws_ecs_container_instance,
  jsonb_array_elements(attachments) as attachment;
```

```sql+sqlite
select
  arn,
  json_extract(attachment.value, '$.id') as attachment_id,
  json_extract(attachment.value, '$.status') as attachment_status,
  json_extract(attachment.value, '$.type') as attachment_type
from
  aws_ecs_container_instance,
  json_each(attachments) as attachment;
```


### List container instances with using a given AMI
Determine the areas in which specific Amazon Machine Images (AMIs) are being used within your container instances. This is particularly useful for identifying potential security risks or for troubleshooting purposes.

```sql+postgres
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

```sql+sqlite
select
  arn,
  json_extract(setting.value, '$.Name') as name,
  json_extract(setting.value, '$.Value') as value
from
  aws_ecs_container_instance,
  json_each(attributes) as setting
where
  json_extract(setting, '$.Name') = 'ecs.ami-id' and
  json_extract(setting, '$.Value') = 'ami-0babb0c4a4e5769b8';
```