# Table: aws_cloudformation_stack_resource

A stack is a collection of AWS resources that you can manage as a single unit.

## Examples

### Basic info

```sql
select
  stack_name,
  stack_id,
  logical_resource_id,
  resource_type,
  resource_status
from
  aws_cloudformation_stack_resource;
```

### List of cloudformation stack resource where rollback is disabled

```sql
select
  s.name,
  s.disable_rollback,
  r.logical_resource_id,
  r.resource_status
from
  aws_cloudformation_stack_resource as r,
  aws_cloudformation_stack as s
where
  r.stack_id = s.id and s.disable_rollback;
```

### List of stack resources where termination protection is not enabled

```sql
select
  s.name,
  s.enable_termination_protection,
  s.disable_rollback,
  r.logical_resource_id,
  r.resource_status
from
  aws_cloudformation_stack_resource as r,
  aws_cloudformation_stack as s
where
  r.stack_id = s.id and not enable_termination_protection;
```

### List VPC resource type resources of the stacks

```sql
select
  stack_name,
  stack_id,
  logical_resource_id,
  resource_status,
  resource_type
from
  aws_cloudformation_stack_resource
where
  resource_type = 'AWS::EC2::VPC';
```

### List resources that are not deleted

```sql
select
  stack_name,
  logical_resource_id,
  resource_status,
  resource_type
from
  aws_cloudformation_stack_resource
where
  resource_status != 'DELETE_COMPLETE';
```
