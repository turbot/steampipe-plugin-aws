---
title: "Table: aws_sfn_state_machine_execution_history - Query AWS Step Functions State Machine Execution History using SQL"
description: "Allows users to query AWS Step Functions State Machine Execution History to fetch information about the execution history of a state machine."
---

# Table: aws_sfn_state_machine_execution_history - Query AWS Step Functions State Machine Execution History using SQL

The `aws_sfn_state_machine_execution_history` table in Steampipe provides information about the execution history of a state machine within AWS Step Functions. This table allows DevOps engineers to query execution-specific details, including execution status, start and end dates, input and output data, and associated metadata. Users can utilize this table to gather insights on state machine executions, such as the status of executions, duration of executions, and verification of input and output data. The schema outlines the various attributes of the state machine execution history, including the execution ARN, state entered time, state exited time, and state name.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sfn_state_machine_execution_history` table, you can use the `.inspect aws_sfn_state_machine_execution_history` command in Steampipe.

**Key columns**:

- `execution_arn` - The Amazon Resource Name (ARN) that identifies the execution. This can be used to join with other tables that contain execution-specific information.
- `state_entered_time` - The date and time when the state was entered. This can be useful for tracking the duration of state executions.
- `state_exited_time` - The date and time when the state was exited. This can be useful for tracking the duration of state executions.

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
