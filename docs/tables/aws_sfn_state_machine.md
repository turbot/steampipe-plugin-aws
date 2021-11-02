# Table: aws_sfn_state_machine

AWS Step Functions makes it easy to coordinate the components of distributed applications as a series of steps in a visual workflow. You can quickly build and run state machines to execute the steps of your application in a reliable and scalable fashion.

## Examples

### Basic info

```sql
select
  name,
  arn,
  status,
  type,
  role_arn
from
  aws_sfn_state_machine;
```

### List active state machines

```sql
select
  name,
  arn,
  status
from
  aws_sfn_state_machine
where
  status = 'ACTIVE';
```