---
title: "Steampipe Table: aws_sfn_state_machine - Query AWS Step Functions State Machines using SQL"
description: "Allows users to query AWS Step Functions State Machines to retrieve detailed information about each state machine."
folder: "Step Functions"
---

# Table: aws_sfn_state_machine - Query AWS Step Functions State Machines using SQL

The AWS Step Functions State Machine is a component of AWS Step Functions service that represents the workflow of a distributed application. It defines a series of steps, their inputs, outputs, and how they interact. With this, you can coordinate microservices, automate processes, and build workflows to orchestrate AWS services and respond to changes.

## Table Usage Guide

The `aws_sfn_state_machine` table in Steampipe provides you with information about State Machines within AWS Step Functions. This table allows you, as a DevOps engineer, to query state machine-specific details, including ARN, name, type, status, creation date, and associated metadata. You can utilize this table to gather insights on state machines, such as their current status, type, role ARN, and more. The schema outlines the various attributes of the state machine for you, including the state machine ARN, creation date, definition, role ARN, and status.

## Examples

### Basic info
Explore which AWS Step Functions state machines are currently active, by identifying their name, status, and associated role. This can be useful in managing your AWS resources and tracking the status of your workflows.

```sql+postgres
select
  name,
  arn,
  status,
  type,
  role_arn
from
  aws_sfn_state_machine;
```

```sql+sqlite
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
Discover the active state machines in your AWS Step Functions to manage and maintain your workflow executions effectively. This can help you focus on workflows that are currently operational, ensuring resources are allocated efficiently.

```sql+postgres
select
  name,
  arn,
  status
from
  aws_sfn_state_machine
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  name,
  arn,
  status
from
  aws_sfn_state_machine
where
  status = 'ACTIVE';
```