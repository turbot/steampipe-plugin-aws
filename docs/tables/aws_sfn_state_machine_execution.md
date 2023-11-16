---
title: "Table: aws_sfn_state_machine_execution - Query AWS Step Functions State Machine Execution using SQL"
description: "Allows users to query AWS Step Functions State Machine Execution data, including execution status, start and end times, and associated state machine details."
---

# Table: aws_sfn_state_machine_execution - Query AWS Step Functions State Machine Execution using SQL

The `aws_sfn_state_machine_execution` table in Steampipe provides information about the execution of state machines within AWS Step Functions. This table allows DevOps engineers to query execution-specific details, including execution status, start and end times, and associated state machine details. Users can utilize this table to gather insights on state machine executions, such as execution duration, status, and associated state machine ARN. The schema outlines the various attributes of the state machine execution, including the state machine ARN, execution ARN, status, start time, end time, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sfn_state_machine_execution` table, you can use the `.inspect aws_sfn_state_machine_execution` command in Steampipe.

### Key columns:

- `execution_arn`: The ARN that identifies the execution. This can be used to join with other tables that provide more information on the specific execution.
- `state_machine_arn`: The ARN that identifies the state machine associated with the execution. This can be used to join with other tables that provide more information on the state machine.
- `status`: The current status of the execution. This is useful for filtering executions based on their status.

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
