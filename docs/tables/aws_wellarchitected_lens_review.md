---
title: "Table: aws_wellarchitected_lens_review - Query AWS Well-Architected Tool Lens Reviews using SQL"
description: "Allows users to query AWS Well-Architected Tool Lens Reviews to obtain detailed information about each review, including its associated workload, lens, and milestone information."
---

# Table: aws_wellarchitected_lens_review - Query AWS Well-Architected Tool Lens Reviews using SQL

The `aws_wellarchitected_lens_review` table in Steampipe provides information about lens reviews within the AWS Well-Architected Tool. This table allows DevOps engineers to query review-specific details, including the associated workload, lens, and milestone information. Users can utilize this table to gather insights on lens reviews, such as the risk level of each review, the number of high-risk issues, and the number of improvement plans. The schema outlines the various attributes of the lens review, including the workload ID, lens alias, milestone number, and other associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_lens_review` table, you can use the `.inspect aws_wellarchitected_lens_review` command in Steampipe.

### Key columns:

- `workload_id`: This is the identifier of the workload. It can be used to join this table with other tables that contain workload-related information.
- `lens_alias`: This is the alias of the lens. It can be used to join this table with other tables that contain lens-related information.
- `milestone_number`: This is the number of the milestone. It can be used to join this table with other tables that contain milestone-related information.

## Examples

### Basic info

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

### List reviews of deprecated lenses

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

### Get high-risk issue counts for each review

```sql
select
  lens_name,
  workload_id,
  risk_counts -> 'HIGH' as high_risk_counts
from
  aws_wellarchitected_lens_review;
```

### Get workload details of each lens review

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

### Get the pillar review summary of lens reviews

```sql
select
  lens_name,
  lens_arn,
  s ->> 'Notes' as pillar_review_summary_note,
  s ->> 'PillarId' as pillar_id,
  s ->> 'PillarName' as pillar_name,
  s ->> 'RiskCounts' as RiskCounts
from
  aws_wellarchitected_lens_review,
  jsonb_array_elements(pillar_review_summaries) as s;
```

### Get risk count details of the lens review

```sql
select
  lens_name,
  lens_arn,
  jsonb_pretty(risk_counts) as risk_counts
from
  aws_wellarchitected_lens_review;
```