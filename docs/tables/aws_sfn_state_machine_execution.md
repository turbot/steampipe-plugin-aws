# Table: aws_sfn_state_machine_execution

A state machine execution occurs when an AWS Step Functions state machine runs and performs its tasks. Each Step Functions state machine can have multiple simultaneous executions.

## Examples

### Basic info

```sql
select
  name,
  execution_arn,
  status,
  state_machine_arn
from
  aws_sfn_state_machine_execution;
```

### List failed executions

```sql
select
  name,
  execution_arn,
  status,
  state_machine_arn
from
  aws_sfn_state_machine_execution
where
  status = 'FAILED';
```
