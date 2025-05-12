---
title: "Steampipe Table: aws_wellarchitected_lens_review - Query AWS Well-Architected Tool Lens Reviews using SQL"
description: "Allows users to query AWS Well-Architected Tool Lens Reviews to obtain detailed information about each review, including its associated workload, lens, and milestone information."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_lens_review - Query AWS Well-Architected Tool Lens Reviews using SQL

The AWS Well-Architected Tool Lens Review is a feature of the AWS Well-Architected Tool. It allows you to review your workloads against the best practices defined in AWS Well-Architected Framework and improve your cloud architectures. It uses SQL queries to provide insights into the performance, cost efficiency, operational excellence, reliability, and security of your workloads.

## Table Usage Guide

The `aws_wellarchitected_lens_review` table in Steampipe provides you with information about lens reviews within the AWS Well-Architected Tool. This table allows you, as a DevOps engineer, to query review-specific details, including the associated workload, lens, and milestone information. You can utilize this table to gather insights on lens reviews, such as the risk level of each review, the number of high-risk issues, and the number of improvement plans. The schema outlines the various attributes of the lens review for you, including the workload ID, lens alias, milestone number, and other associated metadata.

## Examples

### Basic info
Explore the details of your AWS Well-Architected Lens reviews to gain insights into the different lenses you have applied and their respective updates. This can help you manage and optimize your workloads effectively.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that include reviews of outdated lenses in the AWS Well-Architected tool, which can be useful to identify areas for potential updates or replacements.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which high-risk issues are prevalent within each review, in order to prioritize risk mitigation efforts.

```sql+postgres
select
  lens_name,
  workload_id,
  risk_counts -> 'HIGH' as high_risk_counts
from
  aws_wellarchitected_lens_review;
```

```sql+sqlite
select
  lens_name,
  workload_id,
  json_extract(risk_counts, '$.HIGH') as high_risk_counts
from
  aws_wellarchitected_lens_review;
```

### Get workload details of each lens review
Explore the status and version of each lens review in your AWS environment, along with its associated workload details. This helps in assessing the architectural design and environment of the workload, and in identifying any review restrictions.

```sql+postgres
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

```sql+sqlite
select
  r.lens_name,
  r.workload_id,
  r.lens_status,
  r.lens_version,
  w.architectural_design,
  w.environment,
  w.review_restriction_date
from
  aws_wellarchitected_lens_review as r
join
  aws_wellarchitected_workload as w
on
  r.workload_id = w.workload_id;
```

### Get the pillar review summary of lens reviews
Explore the summary of lens reviews in the AWS Well-Architected Tool to gain insights into the assessment of architectural decisions. This query is useful in identifying areas of risk and improvement within your AWS environment.

```sql+postgres
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

```sql+sqlite
select
  lens_name,
  lens_arn,
  json_extract(s.value, '$.Notes') as pillar_review_summary_note,
  json_extract(s.value, '$.PillarId') as pillar_id,
  json_extract(s.value, '$.PillarName') as pillar_name,
  json_extract(s.value, '$.RiskCounts') as RiskCounts
from
  aws_wellarchitected_lens_review,
  json_each(pillar_review_summaries) as s;
```

### Get risk count details of the lens review
Discover the segments that have potential risks within your AWS Well-Architected Lens Review. This is useful for identifying areas that need improvement to ensure your architecture is well-optimized and secure.

```sql+postgres
select
  lens_name,
  lens_arn,
  jsonb_pretty(risk_counts) as risk_counts
from
  aws_wellarchitected_lens_review;
```

```sql+sqlite
select
  lens_name,
  lens_arn,
  json_pretty(risk_counts) as risk_counts
from
  aws_wellarchitected_lens_review;
```