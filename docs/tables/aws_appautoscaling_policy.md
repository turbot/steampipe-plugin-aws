# Table: aws_appautoscaling_policy

Scaling policies defined in Application Auto Scaling

## Examples

### Basic info

```sql
select
  service_namespace,
  scalable_dimension,
  policy_type,
  resource_id,
  creation_time
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs';
```


### List policies for ECS services with policy type Step scaling

```sql
select
  resource_id,
  policy_type
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs'
  and policy_type = 'StepScaling';
```
