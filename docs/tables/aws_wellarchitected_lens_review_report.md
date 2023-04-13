# Table: aws_wellarchitected_lens_review_report

AWS Well-Architected Lens review report.

**Note:**
- `workload_id`, `lens_alias` and `milstone_number` are optional key coulmn quals for the query parameter to get the improvement plans of the lens review.
- For AWS official lenses, this is either the lens alias, such as serverless, or the lens ARN, such as arn:aws:wellarchitected:us-east-1::lens/serverless. Note that some operations (such as ExportLens and CreateLensShare) are not permitted on AWS official lenses.
- For custom lenses, this is the lens ARN, such as arn:aws:wellarchitected:us-west-2:123456789012:lens/0123456789abcdef01234567890abcdef.

`base64_string` value can be used to get the PDF format of the review report.
The tool(https://base64.guru/converter/decode/pdf) can be used to decode the `base64_string` value to a PDF format.
## Examples

### Basic info

```sql
select
  lens_alias,
  lens_arn,
  workload_id,
  milestone_number,
  base64_string
from
  aws_wellarchitected_lens_review_report;
```

### Get workload details for the review report

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
  and r.workload_id = w.workload_id;
```

### Get custom lenses review reports

```sql
select
  r.lens_alias,
  r.lens_arn,
  r.base64_string,
  l.lens_type
from
  aws_wellarchitected_lens_review_report as r,
  aws_wellarchitected_lens as l
where
  l.lens_type <> 'AWS_OFFICIAL';
```