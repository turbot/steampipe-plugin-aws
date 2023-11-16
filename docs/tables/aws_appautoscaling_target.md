---
title: "Table: aws_appautoscaling_target - Query AWS Application Auto Scaling Targets using SQL"
description: "Allows users to query AWS Application Auto Scaling Targets. This table provides information about each target, including the service namespace, scalable dimension, resource ID, and the associated scaling policies."
---

# Table: aws_appautoscaling_target - Query AWS Application Auto Scaling Targets using SQL

The `aws_appautoscaling_target` table in Steampipe provides information about each target within AWS Application Auto Scaling. This table allows DevOps engineers to query target-specific details, including the service namespace, scalable dimension, resource ID, and the associated scaling policies. Users can utilize this table to gather insights on scaling targets, such as the min and max capacity, role ARN, and more. The schema outlines the various attributes of the scaling target, including the resource ID, scalable dimension, creation time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_appautoscaling_target` table, you can use the `.inspect aws_appautoscaling_target` command in Steampipe.

*Key columns:*

- `resource_id`: The unique identifier of the scalable target. This can be used to join this table with other tables to get more detailed information about the resource.
- `service_namespace`: The namespace of the AWS service. This is useful for filtering targets based on the service they belong to.
- `scalable_dimension`: The scalable dimension associated with the target. This can be used to join with other tables to get more detailed information about the scalable dimension.

## Examples

### Basic info

```sql
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

```sql
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
