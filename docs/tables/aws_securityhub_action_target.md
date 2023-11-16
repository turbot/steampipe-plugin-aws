---
title: "Table: aws_securityhub_action_target - Query AWS Security Hub Action Targets using SQL"
description: "Allows users to query AWS Security Hub Action Targets, providing detailed information about each action target within AWS Security Hub, including its ARN, name, and description."
---

# Table: aws_securityhub_action_target - Query AWS Security Hub Action Targets using SQL

The `aws_securityhub_action_target` table in Steampipe provides information about Action Targets within AWS Security Hub. This table allows DevOps engineers to query Action Target-specific details, including its ARN, name, and description. Users can utilize this table to gather insights on Action Targets, such as understanding the purpose of each Action Target, verifying their names and descriptions, and more. The schema outlines the various attributes of the Action Target, including the action target ARN, name, and description.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_securityhub_action_target` table, you can use the `.inspect aws_securityhub_action_target` command in Steampipe.

### Key columns:

- `arn`: The ARN (Amazon Resource Name) of the action target. This is a unique identifier that can be used to join this table with other tables.
- `name`: The name of the action target. This can be used to filter or sort action targets by name, making it easier to find specific targets.
- `description`: The description of the action target. This provides context about the purpose of each action target, which can be useful when auditing or reviewing security configurations.

## Examples

### Basic info

```sql
select
  name,
  arn,
  region
from
  aws_securityhub_action_target;
```

### Get details of a specific action target

```sql
select
  name,
  arn,
  region
from
  aws_securityhub_action_target
where
  arn = 'arn:aws:securityhub:ap-south-1:*****:action/custom/test';
```
