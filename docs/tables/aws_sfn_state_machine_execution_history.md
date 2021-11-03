# Table: aws_sfn_state_machine_execution_history

A state machine execution occurs when an AWS Step Functions state machine runs and performs its tasks. Each Step Functions state machine can have multiple simultaneous executions. Each execution event is stored in history. AWS Step Functions has a hard quota of 25,000 entries in the execution history.

## Examples

### Basic info

```sql
select
  id,
  execution_arn,
  previous_event_id,
  timestamp,
  type
from
  aws_sfn_state_machine_execution_history;
```

### List execution started event details

```sql
select
  id,
  execution_arn,
  execution_started_event_details -> 'Input' as event_input,
  execution_started_event_details -> 'InputDetails' as event_input_details,
  execution_started_event_details ->> 'RoleArn' as event_role_arn
from
  aws_sfn_state_machine_execution_history
where
  type = 'ExecutionStarted';
```
