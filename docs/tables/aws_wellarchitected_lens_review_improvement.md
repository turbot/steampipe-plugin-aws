# Table: aws_wellarchitected_lens_review_improvement

The improvement plan items for your workloads.

**Note:**
- `workload_id` and `lens_alias` are required in the query parameter to get the improvement plans of the lens review.
- For AWS official lenses, this is either the lens alias, such as serverless, or the lens ARN, such as arn:aws:wellarchitected:us-east-1::lens/serverless. Note that some operations (such as ExportLens and CreateLensShare) are not permitted on AWS official lenses.
- For custom lenses, this is the lens ARN, such as arn:aws:wellarchitected:us-west-2:123456789012:lens/0123456789abcdef01234567890abcdef.

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
  aws_wellarchitected_lens_review_improvement
where
  lens_alias = 'wellarchitected'
  and workload_id = '4fca39b680a31bb118be6bc0d177849d';
```

## List review improvements with risk high

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
  aws_wellarchitected_workload
where
  lens_alias = 'wellarchitected'
  and workload_id = '4fca39b680a31bb118be6bc0d177849d'
  and risk = 'HIGH';
```

## Get review improvement risk counts

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

## Get improvement plan details of the review improvements

```sql
select
  lens_alias,
  workload_id,
  p ->> 'ChoiceId' as choice_id,
  p ->> 'DisplayText' as display_text,
  p ->> 'ImprovementPlanUrl' as improvement_plan_url
from
  aws_wellarchitected_lens_review_improvement,
  jsonb_array_elements(improvement_plans) as p
where
  lens_alias = 'wellarchitected'
  and workload_id = '4fca39b680a31bb118be6bc0d177849d';
```