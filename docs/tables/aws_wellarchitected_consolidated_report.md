---
title: "Table: aws_wellarchitected_consolidated_report - Query AWS Well-Architected Tool Consolidated Reports using SQL"
description: "Allows users to query consolidated reports from the AWS Well-Architected Tool, providing a comprehensive view of a workload's alignment with AWS architectural best practices."
---

# Table: aws_wellarchitected_consolidated_report - Query AWS Well-Architected Tool Consolidated Reports using SQL

The `aws_wellarchitected_consolidated_report` table in Steampipe provides information about consolidated reports within the AWS Well-Architected Tool. This table allows DevOps engineers, architects, and other technical professionals to query report-specific details, including findings, risks, and improvement plans. Users can utilize this table to gather insights on workloads, such as high-risk issues, improvement status, and architectural alignment with AWS best practices. The schema outlines the various attributes of the consolidated report, including the workload ID, risk counts, lens name, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_consolidated_report` table, you can use the `.inspect aws_wellarchitected_consolidated_report` command in Steampipe.

### Key columns:

- `workload_id`: This is the unique identifier for the workload. It can be used to join this table with other workload-related tables.
- `lens_name`: This column represents the name of the lens used in the report. It can be useful when joining with other lens-specific tables or filtering reports by lens type.
- `risk_counts`: This column provides a count of risks identified in the report, categorized by level. It is critical for risk analysis and mitigation planning.

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