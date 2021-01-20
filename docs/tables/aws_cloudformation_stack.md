# Table: aws_cloudformation_stack

A stack is a collection of AWS resources that you can manage as a single unit.

## Examples

### Find the status of each cloudformation stack

```sql
select
  name,
  id,
  status
from
  aws_cloudformation_stack;
```


### List of cloudformation stack where rollback is disabled

```sql
select
  name,
  disable_rollback
from
  aws_cloudformation_stack
where
  disable_rollback;
```


### List of stacks where termination protection is not enabled

```sql
select
  name,
  enable_termination_protection
from
  aws_cloudformation_stack
where
  not enable_termination_protection;
```


### Rollback configuration info for each cloudformation stack

```sql
select
  name,
  rollback_configuration ->> 'MonitoringTimeInMinutes' as monitoring_time_in_min,
  rollback_configuration ->> 'RollbackTriggers' as rollback_triggers
from
  aws_cloudformation_stack;
```


### Resource ARNs where notifications about stack actions will be sent

```sql
select
  name,
  jsonb_array_elements_text(notification_arns) as resource_arns
from
  aws_cloudformation_stack;
```