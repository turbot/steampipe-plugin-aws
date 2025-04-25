---
title: "Steampipe Table: aws_appautoscaling_target - Query AWS Application Auto Scaling Targets using SQL"
description: "Allows users to query AWS Application Auto Scaling Targets. This table provides information about each target, including the service namespace, scalable dimension, resource ID, and the associated scaling policies."
folder: "Application Auto Scaling"
---

# Table: aws_appautoscaling_target - Query AWS Application Auto Scaling Targets using SQL

The AWS Application Auto Scaling Targets are used to manage scalable targets within AWS services. These targets can be any resource that can scale in or out, such as an Amazon ECS service, an Amazon EC2 Spot Fleet request, or an Amazon RDS read replica. The Application Auto Scaling service automatically adjusts the resource's capacity to maintain steady, predictable performance at the lowest possible cost.

## Table Usage Guide

The `aws_appautoscaling_target` table in Steampipe provides you with information about each target within AWS Application Auto Scaling. This table allows you, as a DevOps engineer, to query target-specific details, including the service namespace, scalable dimension, resource ID, and the associated scaling policies. You can utilize this table to gather insights on scaling targets, such as the min and max capacity, role ARN, and more. The schema outlines the various attributes of the scaling target for you, including the resource ID, scalable dimension, creation time, and associated tags.

**Important Notes**
- You **_must_** specify `service_namespace` in a `where` clause in order to use this table.
- For supported values of the service namespace, please refer [Service Namespace](https://docs.aws.amazon.com/autoscaling/application/APIReference/API_ScalableTarget.html#autoscaling-Type-ScalableTarget-ServiceNamespace).
- This table supports optional quals. Queries with optional quals are optimised to use additional filtering provided by the AWS API function. Optional quals are supported for the following columns:
  - `resource_id`
  - `scalable_dimension`

## Examples

### Basic info
Explore the creation timeline of resources within the AWS DynamoDB service, which can help in understanding their scalability dimensions and facilitate efficient resource management.

```sql+postgres
select
  service_namespace,
  scalable_dimension,
  resource_id,
  creation_time
from
  aws_appautoscaling_target
where
  service_namespace = 'dynamodb';
```

```sql+sqlite
select
  service_namespace,
  scalable_dimension,
  resource_id,
  creation_time
from
  aws_appautoscaling_target
where
  service_namespace = 'dynamodb';
```


### List targets for DynamoDB tables with read or write auto scaling enabled
Determine the areas in which auto-scaling is enabled for read or write operations in DynamoDB tables. This is useful in managing resources efficiently and optimizing cost by ensuring that scaling only occurs when necessary.

```sql+postgres
select
  resource_id,
  scalable_dimension
from
  aws_appautoscaling_target
where
  service_namespace = 'dynamodb'
  and scalable_dimension = 'dynamodb:table:ReadCapacityUnits'
  or scalable_dimension = 'dynamodb:table:WriteCapacityUnits';
```

```sql+sqlite
select
  resource_id,
  scalable_dimension
from
  aws_appautoscaling_target
where
  service_namespace = 'dynamodb'
  and (scalable_dimension = 'dynamodb:table:ReadCapacityUnits'
  or scalable_dimension = 'dynamodb:table:WriteCapacityUnits');
```
