---
title: "Steampipe Table: aws_wellarchitected_lens_review_improvement - Query AWS Well-Architected Framework Lens Review using SQL"
description: "Allows users to query Lens Review Improvements in the AWS Well-Architected Framework."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_lens_review_improvement - Query AWS Well-Architected Framework Lens Review using SQL

The AWS Well-Architected Framework Lens Review is a tool that helps you analyze and improve your workloads by comparing them against AWS best practices. It uses a set of architectural best practices across five pillars: operational excellence, security, reliability, performance efficiency, and cost optimization. This review provides you with a consistent approach to evaluate architectures and implement designs that can scale over time.

## Table Usage Guide

The `aws_wellarchitected_lens_review_improvement` table in Steampipe provides you with information about Lens Review Improvements within the AWS Well-Architected Framework. This table allows you, as a DevOps engineer, architect, or developer, to query improvement-specific details, including the improvement status, improvement summary, and associated metadata. You can utilize this table to gather insights on improvements, such as those associated with a specific lens review, their status, and a brief summary of the improvement. The schema outlines the various attributes of the Lens Review Improvement for you, including the workload ID, lens alias, improvement ID, and improvement status.

**Important Notes**
-  `workload_id` and `lens_alias` are optional query parameters for filtering out the review improvements with given workload id or lens alias.
- For AWS official lenses, this is either the lens alias, such as serverless, or the lens ARN, such as arn:aws:wellarchitected:us-east-1::lens/serverless. Note that some operations (such as ExportLens and CreateLensShare) are not permitted on AWS official lenses.
- For custom lenses, this is the lens ARN, such as arn:aws:wellarchitected:us-west-2:123456789012:lens/0123456789abcdef01234567890abcdef.

## Examples

## Basic info

```sql+postgres
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

```sql+sqlite
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

```sql+postgres
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

```sql+sqlite
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

```sql+postgres
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

```sql+sqlite
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

```sql+postgres
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

```sql+sqlite
select
  lens_alias,
  workload_id,
  json_extract(p.value, '$.ChoiceId') as choice_id,
  json_extract(p.value, '$.DisplayText') as display_text,
  json_extract(p.value, '$.ImprovementPlanUrl') as improvement_plan_url
from
  aws_wellarchitected_lens_review_improvement,
  json_each(improvement_plans) as p;
```