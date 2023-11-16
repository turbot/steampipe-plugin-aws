---
title: "Table: aws_eventbridge_bus - Query AWS EventBridge Buses using SQL"
description: "Allows users to query AWS EventBridge Buses for detailed information about each bus, including its name, ARN, policy, and more."
---

# Table: aws_eventbridge_bus - Query AWS EventBridge Buses using SQL

The `aws_eventbridge_bus` table in Steampipe provides information about buses within AWS EventBridge. This table allows DevOps engineers to query bus-specific details, including the bus name, ARN, policy, and associated metadata. Users can utilize this table to gather insights on buses, such as their policies, the events they can handle, and more. The schema outlines the various attributes of the EventBridge bus, including the name, ARN, policy, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_eventbridge_bus` table, you can use the `.inspect aws_eventbridge_bus` command in Steampipe.

**Key columns**:

- `name`: The name of the EventBridge bus. This can be used to join this table with others that contain EventBridge bus names.
- `arn`: The Amazon Resource Name (ARN) of the EventBridge bus. This is a unique identifier for the bus and can be used to join this table with others that contain EventBridge bus ARNs.
- `policy`: The policy of the EventBridge bus. This provides information about what actions can be performed on the bus and can be used to join this table with others that contain EventBridge bus policies.

## Examples

### Basic info

```sql
select
  name,
  arn,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_eventbridge_bus;
```
