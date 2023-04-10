# Table: aws_wellarchitected_lens_review

Review the state of your applications and workloads against architectural best practices, identify opportunities for improvement, and track progress over time

## Examples

## Basic info

```sql
select
  lens_name,
  workload_id,
  lens_arn,
  lens_alias,
  lens_version,
  updated_at
from
  aws_wellarchitected_lens_review;
```

## List reviews for deprecated lenses

```sql
select
  lens_name,
  workload_id,
  lens_alias,
  lens_status
from
  aws_wellarchitected_lens_review
where
  lens_status = 'DEPRECATED';
```

## Get high risk issue counts for each review

```sql
select
  lens_name,
  workload_id,
  risk_counts -> 'HIGH' as high_risk_counts
from
  aws_wellarchitected_lens_review;
```

## Get workload details for each lwns review

```sql
select
  r.lens_name,
  r.workload_id,
  r.lens_status,
  r.lens_version,
  w.architectural_design,
  w.environment,
  w.review_restriction_date
from
  aws_wellarchitected_lens_review as r,
  aws_wellarchitected_workload as w
where
  r.workload_id = w.workload_id;
```