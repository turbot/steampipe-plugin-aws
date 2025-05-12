---
title: "Steampipe Table: aws_securityhub_action_target - Query AWS Security Hub Action Targets using SQL"
description: "Allows users to query AWS Security Hub Action Targets, providing detailed information about each action target within AWS Security Hub, including its ARN, name, and description."
folder: "Security Hub"
---

# Table: aws_securityhub_action_target - Query AWS Security Hub Action Targets using SQL

AWS Security Hub Action Targets are specific response actions that can be taken in response to findings. These actions can be custom actions, which you define for your own needs, or AWS managed actions, which are predefined by AWS. They provide a systematic way to initiate a response to specific types of findings.

## Table Usage Guide

The `aws_securityhub_action_target` table in Steampipe provides you with information about Action Targets within AWS Security Hub. This table allows you, as a DevOps engineer, to query Action Target-specific details, including its ARN, name, and description. You can utilize this table to gather insights on Action Targets, such as understanding the purpose of each Action Target, verifying their names and descriptions, and more. The schema outlines the various attributes of the Action Target for you, including the action target ARN, name, and description.

## Examples

### Basic info
Determine the areas in which specific action targets in your AWS Security Hub are located. This can help you manage and prioritize security tasks based on their geographical locations.

```sql+postgres
select
  name,
  arn,
  region
from
  aws_securityhub_action_target;
```

```sql+sqlite
select
  name,
  arn,
  region
from
  aws_securityhub_action_target;
```

### Get details of a specific action target
This example helps you identify specific security actions within your AWS Security Hub, particularly useful when you need to understand the details of a certain action for security auditing or compliance purposes. It's a handy tool for pinpointing actions in a specific region.

```sql+postgres
select
  name,
  arn,
  region
from
  aws_securityhub_action_target
where
  arn = 'arn:aws:securityhub:ap-south-1:*****:action/custom/test';
```

```sql+sqlite
select
  name,
  arn,
  region
from
  aws_securityhub_action_target
where
  arn = 'arn:aws:securityhub:ap-south-1:*****:action/custom/test';
```