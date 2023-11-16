---
title: "Table: aws_eventbridge_rule - Query AWS EventBridge Rule using SQL"
description: "Allows users to query AWS EventBridge Rule to access information regarding the EventBridge rules defined within an AWS account."
---

# Table: aws_eventbridge_rule - Query AWS EventBridge Rule using SQL

The `aws_eventbridge_rule` table in Steampipe provides information about EventBridge rules within AWS EventBridge. This table allows DevOps engineers to query rule-specific details, including the rule name, ARN, state, description, schedule expression, and associated metadata. Users can utilize this table to gather insights on rules, such as the rules associated with a specific event bus, the state of the rules (whether they are enabled or disabled), and more. The schema outlines the various attributes of the EventBridge rule, including the rule ARN, event bus name, description, state, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_eventbridge_rule` table, you can use the `.inspect aws_eventbridge_rule` command in Steampipe.

### Key columns:

- `name`: The name of the rule. This can be used to join with other tables that need to reference the rule.
- `arn`: The Amazon Resource Name (ARN) of the rule. This can be used to join with other tables that use ARNs for identification.
- `event_bus_name`: The name of the event bus associated with the rule. This can be used to join with tables that contain information about event buses.

## Examples

### Basic info

```sql
select
  name,
  arn,
  state,
  created_by,
  event_bus_name
from
  aws_eventbridge_rule;
```


### List disabled rules

```sql
select
  name,
  arn,
  state,
  created_by
from
  aws_eventbridge_rule
where
  state != 'ENABLED';
```


### Get the target information for each rule

```sql
select
  name,
  cd ->> 'Id' as target_id,
  cd ->> 'Arn' as target_arn,
  cd ->> 'RoleArn' as role_arn
from
  aws_eventbridge_rule,
  jsonb_array_elements(targets) as cd;
```
