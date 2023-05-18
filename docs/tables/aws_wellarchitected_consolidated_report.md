# Table: aws_wellarchitected_consolidated_report

The AWS Well-Architected Consolidated Report is an enhanced version of the Well-Architected Review report. It consolidates the results of multiple workload reviews into a single comprehensive report. This allows you to gain insights across multiple workloads, identify common patterns or issues, and track improvement progress over time.

**Note:** The column `base64_string` value is a Base64-encoded string representation of a lens review report. This data can be used to create a PDF file.
The tool https://base64.guru/converter/decode/pdf can be used for converting the Base64-encoded string to a PDF format.

## Examples

### Basic info

```sql
select
  workload_name,
  workload_arn,
  workload_id,
  lenses_applied_count,
  metric_type,
  updated_at
from
  aws_wellarchitected_consolidated_report;
```

### Get workload details for each consolidated report

```sql
select
  r.workload_name,
  r.workload_arn,
  r.workload_id,
  r.lenses_applied_count,
  w.environment as workload_environment,
  w.improvement_status as workload_improvement_status,
  w.review_restriction_date as workload_review_restriction_date
from
  aws_wellarchitected_consolidated_report as r,
  aws_wellarchitected_workload as w
where
  w.workload_id = r.workload_id;
```

### Get high-risk issue counts for each consolidated report

```sql
select
  workload_name,
  workload_id,
  risk_counts -> 'HIGH' as high_risk_counts
from
  aws_wellarchitected_consolidated_report;
```

### Get lens details for each consolidated report

```sql
select
  workload_name,
  workload_id,
  l ->> 'LensArn' as lens_arn,
  l -> 'Pillars' as pillars,
  l -> 'RiskCounts' as risk_counts
from
  aws_wellarchitected_consolidated_report,
  jsonb_array_elements(lenses) as l;
```