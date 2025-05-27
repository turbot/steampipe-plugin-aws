---
title: "Steampipe Table: aws_wellarchitected_lens_review_report - Query AWS Well-Architected Tool Lens Review Report using SQL"
description: "Allows users to query Lens Review Reports in the AWS Well-Architected Tool."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_lens_review_report - Query AWS Well-Architected Tool Lens Review Report using SQL

The AWS Well-Architected Tool Lens Review Report is a feature of the AWS Well-Architected Tool. It helps you review your workloads against AWS architectural best practices, and provides guidance on improving your cloud architectures. This tool conducts a comprehensive review of your applications, identifying any high-risk issues and providing strategies to mitigate them.

## Table Usage Guide

The `aws_wellarchitected_lens_review_report` table in Steampipe provides you with information about Lens Review Reports within the AWS Well-Architected Tool. This table allows you, as a DevOps engineer, architect, or system administrator, to query details about the lens review reports, including the lens alias, lens name, lens version, and associated metadata. You can utilize this table to gather insights on lens review reports, such as identifying the lens version, verifying the lens status, understanding the lens notes, and more. The schema outlines the various attributes of the Lens Review Report for you, including the lens version, lens status, lens notes, and associated tags.

**Important Notes**
- `workload_id`, `lens_alias` and `milstone_number` are optional key column quals for the query parameter to get the improvement plans of the lens review.
- For AWS official lenses, this is either the lens alias, such as serverless, or the lens ARN, such as arn:aws:wellarchitected:us-east-1::lens/serverless. Note that some operations (such as ExportLens and CreateLensShare) are not permitted on AWS official lenses.
- For custom lenses, this is the lens ARN, such as arn:aws:wellarchitected:us-west-2:123456789012:lens/0123456789abcdef01234567890abcdef.
- The `base64_string` column value can be used to get the PDF format of the review report.
The tool(https://base64.guru/converter/decode/pdf) can be used to decode the `base64_string` value to a PDF format.

## Examples

### Basic info
Explore the milestones and workloads associated with different lenses in AWS Well-Architected tool to understand the progress and status of your architectural reviews. This can assist in identifying areas for improvement and planning future actions.

```sql+postgres
select
  lens_alias,
  lens_arn,
  workload_id,
  milestone_number,
  base64_string
from
  aws_wellarchitected_lens_review_report;
```

```sql+sqlite
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
Determine the areas in which workload details are needed for a review report. This is useful for understanding the environment and whether the review owner has acknowledged updates, which can help in managing workloads effectively.

```sql+postgres
select
  w.workload_name,
  r.workload_id,
  r.base64_string,
  w.environment,
  w.is_review_owner_update_acknowledged
from
  aws_wellarchitected_lens_review_report as r,
  aws_wellarchitected_workload as w
where
  r.workload_id = w.workload_id;
```

```sql+sqlite
select
  w.workload_name,
  r.workload_id,
  r.base64_string,
  w.environment,
  w.is_review_owner_update_acknowledged
from
  aws_wellarchitected_lens_review_report as r
join
  aws_wellarchitected_workload as w
on
  r.workload_id = w.workload_id;
```

### Get the review report of custom lenses
Explore the review reports of custom lenses in AWS Well-Architected Tool, focusing on lenses that are not officially provided by AWS. This can help in assessing the performance of custom lenses and identifying potential areas for improvement.

```sql+postgres
select
  r.lens_alias,
  r.lens_arn,
  r.base64_string,
  l.lens_type
from
  aws_wellarchitected_lens_review_report as r,
  aws_wellarchitected_lens as l
where
  l.lens_type <> 'aws_OFFICIAL';
```

```sql+sqlite
select
  r.lens_alias,
  r.lens_arn,
  r.base64_string,
  l.lens_type
from
  aws_wellarchitected_lens_review_report as r,
  aws_wellarchitected_lens as l
where
  l.lens_type <> 'aws_OFFICIAL';
```