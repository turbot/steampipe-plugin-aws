---
title: "Steampipe Table: aws_wafv2_rule_group - Query AWS WAFv2 Rule Groups using SQL"
description: "Allows users to query AWS WAFv2 Rule Groups and gather information such as the group's ARN, capacity, description, rules, visibility configuration, and more."
folder: "WAFv2"
---

# Table: aws_wafv2_rule_group - Query AWS WAFv2 Rule Groups using SQL

The AWS WAFv2 Rule Groups is a feature of AWS Web Application Firewall (WAF) that allows you to encapsulate a set of rules that you can reuse across multiple web ACLs. It helps you manage similar set of rules across your AWS resources without having to recreate them individually. This aids in maintaining consistency in your security posture and simplifies the management of your WAF configurations.

## Table Usage Guide

The `aws_wafv2_rule_group` table in Steampipe allows you to query details related to rule groups in AWS WAFv2 (Web Application Firewall version 2). You can use this table to gather information such as the ARN, capacity, description, rules, visibility configuration, and more about each rule group. As a DevOps engineer or security professional, you can utilize this table to analyze and manage rule groups, monitor rule group capacity, review rule configurations, and ensure visibility settings are correctly configured. The schema outlines the various attributes of the rule group for you, including the ARN, capacity, description, rules, visibility configuration, and associated tags.

## Examples

### Basic info
This query allows you to gain insights into various rule groups within your AWS Web Application Firewall (WAF) version 2. It's useful for understanding the configuration and capacity of each rule group, as well as their associated regions, helping to manage and optimize your security settings.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which global rule groups are applied within CloudFront. This is useful for assessing security measures and identifying potential vulnerabilities across your cloud-based content delivery network.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which rule groups are operating with less than optimal web access control (WAC) capacity. This is useful for identifying potential vulnerabilities or inefficiencies in your system's security settings.

```sql+postgres
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

```sql+sqlite
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