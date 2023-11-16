---
title: "Table: aws_cloudformation_stack - Query AWS CloudFormation Stack using SQL"
description: "Allows users to query AWS CloudFormation Stack data, including stack name, status, creation time, and associated tags."
---

# Table: aws_cloudformation_stack - Query AWS CloudFormation Stack using SQL

The `aws_cloudformation_stack` table in Steampipe provides information about stacks within AWS CloudFormation. This table allows DevOps engineers to query stack-specific details, including stack name, status, creation time, and associated tags. Users can utilize this table to gather insights on stacks, such as stack status, stack resources, stack capabilities, and more. The schema outlines the various attributes of the CloudFormation stack, including stack ID, stack name, creation time, stack status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudformation_stack` table, you can use the `.inspect aws_cloudformation_stack` command in Steampipe.

**Key columns**:

- `stack_name`: The name associated with the stack. This can be used to join this table with other tables that contain stack information.
- `stack_id`: The unique identifier for the stack. This can be used to join this table with other tables that contain stack-specific details.
- `stack_status`: The current status of the stack. This can be useful in determining the health and state of the stack.

## Examples

### Find the status of each cloudformation stack

```sql
select
  name,
  id,
  status
from
  aws_cloudformation_stack;
```


### List of cloudformation stack where rollback is disabled

```sql
select
  name,
  disable_rollback
from
  aws_cloudformation_stack
where
  disable_rollback;
```


### List of stacks where termination protection is not enabled

```sql
select
  name,
  enable_termination_protection
from
  aws_cloudformation_stack
where
  not enable_termination_protection;
```


### Rollback configuration info for each cloudformation stack

```sql
select
  name,
  rollback_configuration ->> 'MonitoringTimeInMinutes' as monitoring_time_in_min,
  rollback_configuration ->> 'RollbackTriggers' as rollback_triggers
from
  aws_cloudformation_stack;
```


### Resource ARNs where notifications about stack actions will be sent

```sql
select
  name,
  jsonb_array_elements_text(notification_arns) as resource_arns
from
  aws_cloudformation_stack;
```