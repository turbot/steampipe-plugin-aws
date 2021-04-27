# Table: aws_appautoscaling_target

Application Auto Scaling allows you to automatically scale your scalable resources according to conditions that you define.

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


### Check whether autoscaling is enabled of both read-write capacity

```sql
select
  resource_id,
  scalable_dimension
from
  aws_appautoscaling_target
where
  service_namespace = 'dynamodb'
  and scalable_dimension = 'dynamodb:table:ReadCapacityUnits'
  or 'scalable_dimension' = 'dynamodb:table:WriteCapacityUnits';
```
