---
title: "Table: aws_inspector_assessment_template - Query AWS Inspector Assessment Templates using SQL"
description: "Allows users to query AWS Inspector Assessment Templates to gain insights into each template's configuration, including ARN, duration, rules package ARNs, and user attributes for findings."
---

# Table: aws_inspector_assessment_template - Query AWS Inspector Assessment Templates using SQL

The `aws_inspector_assessment_template` table in Steampipe provides information about assessment templates within AWS Inspector. This table allows DevOps engineers, security analysts, and other technical professionals to query template-specific details, including the ARN, duration, rules package ARNs, and user attributes for findings. Users can utilize this table to gather insights on assessment templates, such as identifying templates with specific rules, verifying template configurations, and more. The schema outlines the various attributes of the assessment template, including the template ARN, duration, rules package ARNs, user attributes for findings, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_inspector_assessment_template` table, you can use the `.inspect aws_inspector_assessment_template` command in Steampipe.

**Key columns**:

- `arn`: The ARN of the assessment template. This is the unique identifier of the template and can be used to join this table with other AWS Inspector tables.
- `duration_in_seconds`: The duration of the assessment run. This information can be useful in planning and scheduling assessment runs.
- `rules_package_arns`: The ARNs of the rules packages that are specified for the assessment template. These can be used to join this table with the `aws_inspector_rules_package` table for a comprehensive view of the rules applied in the assessment templates.

## Examples

### Basic info

```sql
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  region
from
  aws_inspector_assessment_template;
```


### List assessment templates that have no assigned finding attributes

```sql
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  user_attributes_for_findings = '[]';
```


### List assessment templates that have no assessment runs

```sql
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  assessment_run_count = 0;
```


### List assessment templates with run duration less than 1 hour

```sql
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  duration_in_seconds,
  region
from
  aws_inspector_assessment_template
where
  duration_in_seconds < 3600;
```
