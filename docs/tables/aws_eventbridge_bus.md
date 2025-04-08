---
title: "Steampipe Table: aws_eventbridge_bus - Query AWS EventBridge Buses using SQL"
description: "Allows users to query AWS EventBridge Buses for detailed information about each bus, including its name, ARN, policy, and more."
folder: "EventBridge"
---

# Table: aws_eventbridge_bus - Query AWS EventBridge Buses using SQL

The AWS EventBridge Bus is a component of AWS EventBridge service, which enables the routing of events between AWS services and applications. It acts as a central hub where these events are filtered and routed to specific targets, based on rules you define. The EventBridge Bus helps in building loosely coupled, distributed applications with nearly real-time event-driven architectures.

## Table Usage Guide

The `aws_eventbridge_bus` table in Steampipe provides you with information about buses within AWS EventBridge. This table allows you, as a DevOps engineer, to query bus-specific details, including the bus name, ARN, policy, and associated metadata. You can utilize this table to gather insights on buses, such as their policies, the events they can handle, and more. The schema outlines the various attributes of the EventBridge bus for you, including the name, ARN, policy, and associated tags.

## Examples

### Basic info
Gain insights into the configurations of your AWS EventBridge Bus to ensure they align with your security and operational requirements. This can help you manage and monitor your AWS resources effectively.

```sql+postgres
select
  name,
  arn,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_eventbridge_bus;
```

```sql+sqlite
select
  name,
  arn,
  policy,
  policy_std
from
  aws_eventbridge_bus;
```