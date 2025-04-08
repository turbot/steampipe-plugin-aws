---
title: "Steampipe Table: aws_shield_protection_group - Query AWS Shield Advanced Protection Groups using SQL"
description: "Allows users to query AWS Shield Advanced Protection Groups and retrieve detailed information about each Protection Group."
folder: "Shield"
---

# Table: aws_shield_protection_group - Query AWS Shield Advanced Protection Groups using SQL

AWS Shield Advanced Protection Groups are logical collections of your Shield Advanced protected resources. AWS Shield Advanced protection groups give you a self-service way to customize the scope of detection and mitigation by treating multiple protected resources as a single unit. Protection groups can, for example, help reduce false positives in situations such as blue/green swap, where resources alternate between being near zero load and fully loaded.

## Table Usage Guide

The `aws_shield_protection` table in Steampipe allows you to query AWS Shield Advanced Protection Groups and retrieve information like the resources included in the group or the aggregation method used for the group. For more information about the individual columns and their values, please refer to the [official AWS documentation](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_DescribeProtectionGroup.html#API_DescribeProtectionGroup_ResponseSyntax).

## Examples

### Basic info

```sql+postgres
select
  protection_group_id,
  aggregation,
  pattern,
  resource_type
from
  aws_shield_protection_group;
```

```sql+sqlite
select
  protection_group_id,
  aggregation,
  pattern,
  resource_type
from
  aws_shield_protection_group;
```

### List all members of protection groups with the pattern `ARBITRARY`

```sql+postgres
select
  protection_group_id,
  member
from
  aws_shield_protection_group,
  jsonb_array_elements_text(members) as member
where
  pattern = 'ARBITRARY'
order by
  protection_group_id;
```

```sql+sqlite
select
  protection_group_id,
  member
from
  aws_shield_protection_group,
  json_each(members) as member
where
  pattern = 'ARBITRARY'
order by
  protection_group_id;
```
