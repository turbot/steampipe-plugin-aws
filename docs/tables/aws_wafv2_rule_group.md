---
title: "Table: aws_wafv2_rule_group - Query AWS WAFv2 Rule Groups using SQL"
description: "Allows users to query AWS WAFv2 Rule Groups and gather information such as the group's ARN, capacity, description, rules, visibility configuration, and more."
---

# Table: aws_wafv2_rule_group - Query AWS WAFv2 Rule Groups using SQL

The `aws_wafv2_rule_group` table in Steampipe allows users to query details related to rule groups in AWS WAFv2 (Web Application Firewall version 2). This table can be used to gather information such as the ARN, capacity, description, rules, visibility configuration, and more about each rule group. DevOps engineers and security professionals can utilize this table to analyze and manage rule groups, monitor rule group capacity, review rule configurations, and ensure visibility settings are correctly configured. The schema outlines the various attributes of the rule group, including the ARN, capacity, description, rules, visibility configuration, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wafv2_rule_group` table, you can use the `.inspect aws_wafv2_rule_group` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the rule group. This can be used to join this table with other tables that also have AWS ARN information.
- `capacity`: The capacity of the rule group. This is important to monitor and manage the capacity of your rule groups.
- `rules`: The rules defined in the rule group. This can be useful to analyze and review the configuration of each rule within the group.

## Examples

### Basic info

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  capacity,
  rules,
  region
from
  aws_wafv2_rule_group;
```

### List global (CloudFront) rule groups

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  capacity,
  region
from
  aws_wafv2_rule_group
where
  scope = 'CLOUDFRONT';
```

### List rule groups with fewer than 5 web ACL capacity units (WCUs)

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  capacity,
  region
from
  aws_wafv2_rule_group
where
  capacity < 5;
```
