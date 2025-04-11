---
title: "Steampipe Table: aws_sfn_state_machine_execution - Query AWS Step Functions State Machine Execution using SQL"
description: "Allows users to query AWS Step Functions State Machine Execution data, including execution status, start and end times, and associated state machine details."
folder: "Step Functions"
---

# Table: aws_sfn_state_machine_execution - Query AWS Step Functions State Machine Execution using SQL

The AWS Step Functions State Machine Execution is a feature of AWS Step Functions that allows you to coordinate multiple AWS services into serverless workflows so you can build and update apps quickly. Using Step Functions, you can design and run workflows that stitch together services, such as AWS Lambda, AWS Fargate, and Amazon SageMaker, into feature-rich applications. Workflows are made up of a series of steps, with the output of one step acting as input into the next.

## Table Usage Guide

The `aws_sfn_state_machine_execution` table in Steampipe provides you with information about the execution of state machines within AWS Step Functions. This table allows you, as a DevOps engineer, to query execution-specific details, including execution status, start and end times, and associated state machine details. You can utilize this table to gather insights on state machine executions, such as execution duration, status, and associated state machine ARN. The schema outlines for you the various attributes of the state machine execution, including the state machine ARN, execution ARN, status, start time, end time, and more.

## Examples

### Basic info
Gain insights into the status of your AWS Step Functions state machine executions to understand their current operational state and monitor their progress. This can help in identifying any potential issues or bottlenecks in your workflows.

```sql+postgres
select
  name,
  execution_arn,
  status,
  state_machine_arn
from
  aws_sfn_state_machine_execution;
```

```sql+sqlite
select
  name,
  execution_arn,
  status,
  state_machine_arn
from
  aws_sfn_state_machine_execution;
```

### List failed executions
Determine the areas in which AWS Step Functions executions have failed, allowing you to identify problematic workflows and troubleshoot accordingly.

```sql+postgres
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

```sql+sqlite
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