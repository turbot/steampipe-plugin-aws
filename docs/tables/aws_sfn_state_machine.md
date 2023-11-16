---
title: "Table: aws_sfn_state_machine - Query AWS Step Functions State Machines using SQL"
description: "Allows users to query AWS Step Functions State Machines to retrieve detailed information about each state machine."
---

# Table: aws_sfn_state_machine - Query AWS Step Functions State Machines using SQL

The `aws_sfn_state_machine` table in Steampipe provides information about State Machines within AWS Step Functions. This table allows DevOps engineers to query state machine-specific details, including ARN, name, type, status, creation date, and associated metadata. Users can utilize this table to gather insights on state machines, such as their current status, type, role ARN, and more. The schema outlines the various attributes of the state machine, including the state machine ARN, creation date, definition, role ARN, and status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sfn_state_machine` table, you can use the `.inspect aws_sfn_state_machine` command in Steampipe.

**Key columns**:

- **name**: The name of the state machine. This is a unique identifier and can be used to join this table with others that contain state machine information.
- **arn**: The Amazon Resource Number (ARN) of the state machine. This is a globally unique identifier that can be used to join this table with any other AWS resource table.
- **status**: The current status of the state machine. This column can be useful for filtering or sorting state machines based on their status.

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