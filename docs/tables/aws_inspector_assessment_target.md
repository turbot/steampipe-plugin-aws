---
title: "Table: aws_inspector_assessment_target - Query AWS Inspector Assessment Targets using SQL"
description: "Allows users to query AWS Inspector Assessment Targets. The `aws_inspector_assessment_target` table in Steampipe provides information about assessment targets within AWS Inspector. This table allows DevOps engineers to query target-specific details, including ARN, name, and associated resource group ARN. Users can utilize this table to gather insights on assessment targets, such as their creation time, last updated time, and more. The schema outlines the various attributes of the assessment target, including the target ARN, creation date, and associated tags."
---

# Table: aws_inspector_assessment_target - Query AWS Inspector Assessment Targets using SQL

The `aws_inspector_assessment_target` table in Steampipe provides information about assessment targets within AWS Inspector. This table allows DevOps engineers to query target-specific details, including ARN, name, and associated resource group ARN. Users can utilize this table to gather insights on assessment targets, such as their creation time, last updated time, and more. The schema outlines the various attributes of the assessment target, including the target ARN, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the aws_inspector_assessment_target table, you can use the `.inspect aws_inspector_assessment_target` command in Steampipe.

**Key columns**:

- `arn`: The ARN of the assessment target. It can be used to join this table with other AWS tables.
- `name`: The name of the assessment target. It provides a human-readable identifier for the target.
- `resource_group_arn`: The ARN of the resource group associated with the assessment target. This can be used to join with the `aws_inspector_resource_group` table for more detailed information.

## Examples

### Basic info

```sql
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target;
```


### List assessment targets created within the last 7 days

```sql
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target
where
  created_at > (current_date - interval '7' day);
```


### List assessment targets that were updated after creation

```sql
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target
where
  created_at != updated_at;
```
