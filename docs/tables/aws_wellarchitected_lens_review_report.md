# Table: aws_wellarchitected_lens_review_report

AWS Well-Architected Lens review report.

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
  milestone_number,
  base64_string
from
  aws_wellarchitected_lens_review_report
where
  lens_alias = 'wellarchitected'
  and workload_id = '4fca39b680a31bb118be6bc0d177849d';
```

## Get workload details for the review report

```sql
select
  r.workload_name,
  r.workload_id,
  r.base64_string,
  w.environment,
  w.is_review_owner_update_acknowledged
from
  aws_wellarchitected_lens_review_report as r,
  aws_wellarchitected_workload as w
where
  r.lens_alias = 'wellarchitected'
  and r.workload_id = '4fca39b680a31bb118be6bc0d177849d'
  and r.workload_id = w.workload_id;
```