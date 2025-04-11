---
title: "Steampipe Table: aws_sfn_state_machine_execution_history - Query AWS Step Functions State Machine Execution History using SQL"
description: "Allows users to query AWS Step Functions State Machine Execution History to fetch information about the execution history of a state machine."
folder: "Step Functions"
---

# Table: aws_sfn_state_machine_execution_history - Query AWS Step Functions State Machine Execution History using SQL

The AWS Step Functions State Machine Execution History is a feature of AWS Step Functions that allows you to track the execution history of your state machines. It provides a detailed, near real-time, view of the step-by-step progress of each execution. This helps in monitoring workflows, diagnosing issues, and understanding the state transitions in your applications.

## Table Usage Guide

The `aws_sfn_state_machine_execution_history` table in Steampipe provides you with information about the execution history of a state machine within AWS Step Functions. This table allows you, as a DevOps engineer, to query execution-specific details, including execution status, start and end dates, input and output data, and associated metadata. You can utilize this table to gather insights on state machine executions, such as the status of executions, duration of executions, and verification of input and output data. The schema outlines the various attributes of the state machine execution history for you, including the execution ARN, state entered time, state exited time, and state name.

## Examples

### Basic info
Analyze the history of AWS Step Functions state machine executions to understand the sequence and timing of events. This can be useful for debugging purposes or gaining insights into workflow performance and efficiency.

```sql+postgres
select
  id,
  execution_arn,
  previous_event_id,
  timestamp,
  type
from
  aws_sfn_state_machine_execution_history;
```

```sql+sqlite
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
This query can be used to gain insights into the specific details of events that initiated the execution of a process in a state machine. It's particularly useful in monitoring and debugging scenarios, where understanding the input and role associated with the start of an execution can help identify potential issues or inefficiencies.

```sql+postgres
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

```sql+sqlite
select
  id,
  execution_arn,
  json_extract(execution_started_event_details, '$.Input') as event_input,
  json_extract(execution_started_event_details, '$.InputDetails') as event_input_details,
  json_extract(execution_started_event_details, '$.RoleArn') as event_role_arn
from
  aws_sfn_state_machine_execution_history
where
  type = 'ExecutionStarted';
```