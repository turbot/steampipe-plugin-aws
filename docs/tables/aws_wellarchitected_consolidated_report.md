---
title: "Steampipe Table: aws_wellarchitected_consolidated_report - Query AWS Well-Architected Tool Consolidated Reports using SQL"
description: "Allows users to query consolidated reports from the AWS Well-Architected Tool, providing a comprehensive view of a workload's alignment with AWS architectural best practices."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_consolidated_report - Query AWS Well-Architected Tool Consolidated Reports using SQL

The AWS Well-Architected Tool is a service that helps you review the state of your workloads and compares them to the latest AWS architectural best practices. The service provides a consolidated report that includes the workload's risks and improvement plan. This tool is designed to provide high-level guidance and best practices for architects and developers, helping ensure the efficiency, cost-effectiveness, and reliability of your applications.

## Table Usage Guide

The `aws_wellarchitected_consolidated_report` table in Steampipe provides you with information about consolidated reports within the AWS Well-Architected Tool. This table allows you, as a DevOps engineer, architect, or other technical professional, to query report-specific details, including findings, risks, and improvement plans. You can utilize this table to gather insights on workloads, such as high-risk issues, improvement status, and architectural alignment with AWS best practices. The schema outlines the various attributes of the consolidated report for you, including the workload ID, risk counts, lens name, and associated metadata.

**Important Notes**
- The column `base64_string` value is a Base64-encoded string representation of a lens review report. This data can be used to create a PDF file.
- The tool https://base64.guru/converter/decode/pdf can be used for converting the Base64-encoded string to a PDF format.

## Examples

### Basic info
Explore the key details of your AWS workloads, including the number of applied lenses and the type of metrics used. This insight can help in understanding the overall configuration and recent updates to your workloads, thus aiding in efficient management and optimization.

```sql+postgres
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

```sql+sqlite
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
Explore the details of each workload in your consolidated reports, including the applied lenses count, environment, improvement status, and review restriction date. This can help you understand the current state of your workloads and identify areas for improvement.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in your AWS workloads where high-risk issues are prevalent. This query helps you understand where potential vulnerabilities exist, allowing you to prioritize and address these risks effectively.

```sql+postgres
select
  workload_name,
  workload_id,
  risk_counts -> 'HIGH' as high_risk_counts
from
  aws_wellarchitected_consolidated_report;
```

```sql+sqlite
select
  workload_name,
  workload_id,
  json_extract(risk_counts, '$.HIGH') as high_risk_counts
from
  aws_wellarchitected_consolidated_report;
```

### Get lens details for each consolidated report
Determine the areas in which each lens contributes to a consolidated report within the AWS well-architected framework. This allows for a comprehensive analysis of workload risks and potential improvements.

```sql+postgres
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

```sql+sqlite
select
  workload_name,
  workload_id,
  json_extract(l.value, '$.LensArn') as lens_arn,
  json_extract(l.value, '$.Pillars') as pillars,
  json_extract(l.value, '$.RiskCounts') as risk_counts
from
  aws_wellarchitected_consolidated_report,
  json_each(aws_wellarchitected_consolidated_report.lenses) as l;
```