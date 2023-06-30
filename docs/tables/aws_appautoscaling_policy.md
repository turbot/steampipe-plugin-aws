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

### List policies for ECS services created in the last 30 days

```sql
select
  resource_id,
  policy_type
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs'
  and creation_time > now() - interval '30 days';
```

### Get the CloudWatch alarms associated with the Auto Scaling policy

```sql
select
  resource_id,
  policy_type,
  jsonb_array_elements(alarms) -> 'AlarmName' as alarm_name
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs';
```

### Get the configuration for Step scaling type policies

```sql
select
  resource_id,
  policy_type,
  step_scaling_policy_configuration
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs'
  and policy_type = 'StepScaling';
```
