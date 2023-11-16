---
title: "Table: aws_wellarchitected_lens_review_improvement - Query AWS Well-Architected Framework Lens Review using SQL"
description: "Allows users to query Lens Review Improvements in the AWS Well-Architected Framework."
---

# Table: aws_wellarchitected_lens_review_improvement - Query AWS Well-Architected Framework Lens Review using SQL

The `aws_wellarchitected_lens_review_improvement` table in Steampipe provides information about Lens Review Improvements within the AWS Well-Architected Framework. This table allows DevOps engineers, architects, and developers to query improvement-specific details, including the improvement status, improvement summary, and associated metadata. Users can utilize this table to gather insights on improvements, such as those associated with a specific lens review, their status, and a brief summary of the improvement. The schema outlines the various attributes of the Lens Review Improvement, including the workload ID, lens alias, improvement ID, and improvement status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_lens_review_improvement` table, you can use the `.inspect aws_wellarchitected_lens_review_improvement` command in Steampipe.

### Key columns:

- `workload_id`: This is the unique identifier of the workload. It is important as it can be used to join this table with other tables related to workloads.
- `lens_alias`: This is the alias of the lens in the workload. It is useful as it can be used to join this table with other tables related to lenses.
- `improvement_id`: This is the unique identifier of the improvement. It is important as it can be used to join this table with other tables related to improvements.

## Examples

## Basic info

```sql
select
  lens_alias,
  lens_arn,
  workload_id,
  improvement_plan_url,
  pillar_id,
  question_id,
  question_title
from
  aws_wellarchitected_lens_review_improvement;
```

## List review improvements with risk high for a workload

```sql
select
  lens_alias,
  lens_arn,
  workload_id,
  improvement_plan_url,
  question_id,
  question_title,
  risk
from
  aws_wellarchitected_lens_review_improvement
where
  workload_id = '4fca39b680a31bb118be6bc0d177849d'
  and risk = 'HIGH';
```

## Get review improvement risk counts for a particular workload and lens

```sql
select
  lens_arn,
  workload_id,
  risk,
  count(risk)
from
  aws_wellarchitected_lens_review_improvement
where
  lens_alias = 'wellarchitected'
  and workload_id = '4fca39b680a31bb118be6bc0d177849d'
group by
  risk,
  lens_arn,
  workload_id;
```

## Get improvement plan details of the review improvements for each workload

```sql
select
  lens_alias,
  workload_id,
  p ->> 'ChoiceId' as choice_id,
  p ->> 'DisplayText' as display_text,
  p ->> 'ImprovementPlanUrl' as improvement_plan_url
from
  aws_wellarchitected_lens_review_improvement,
  jsonb_array_elements(improvement_plans) as p;
```