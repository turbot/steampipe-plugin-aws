---
title: "Table: aws_inspector_assessment_run - Query AWS Inspector Assessment Runs using SQL"
description: "Allows users to query AWS Inspector Assessment Runs to get detailed information about each assessment run, including its state, duration, findings, and more."
---

# Table: aws_inspector_assessment_run - Query AWS Inspector Assessment Runs using SQL

The `aws_inspector_assessment_run` table in Steampipe provides information about assessment runs within AWS Inspector. This table allows DevOps engineers to query run-specific details, including its state, duration, findings, and associated metadata. Users can utilize this table to gather insights on runs, such as the number of findings, the state of the run, and the time it took for the run to complete. The schema outlines the various attributes of the assessment run, including the run ARN, creation date, state, duration, findings, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_inspector_assessment_run` table, you can use the `.inspect aws_inspector_assessment_run` command in Steampipe.

**Key columns**:

- `arn`: The ARN of the assessment run. This is the unique identifier of the run and can be used to join this table with other tables that also contain AWS Inspector Assessment Run ARNs.
- `state`: The state of the assessment run. This can be used to filter runs based on their state (e.g., COMPLETED, STARTED, STOPPED, etc.).
- `duration_in_seconds`: The duration of the assessment run. This can be used to analyze the time taken by different runs and optimize accordingly.

## Examples

### Basic info

```sql
select
  name,
  arn,
  assessment_template_arn,
  created_at,
  state,
  region
from
  aws_inspector_assessment_run;
```

### List finding counts by severity

```sql
select
  name,
  finding_counts ->> 'High' as high,
  finding_counts ->> 'Low' as low,
  finding_counts ->> 'Medium' as medium,
  finding_counts ->> 'Informational' as informational,
  state
from
  aws_inspector_assessment_run;
```

### List assessment runs for each assessment template

```sql
select
  t.name as assessment_template_name,
  r.name as assessment_run_name,
  r.created_at as assessment_run_created_at,
  r.state,
  r.region
from
  aws_inspector_assessment_run as r,
  aws_inspector_assessment_template as t
where
  r.assessment_template_arn = t.arn;
```

### List assessment runs which are not completed

```sql
select
  name,
  arn,
  assessment_template_arn,
  created_at,
  state,
  region
from
  aws_inspector_assessment_run
where
  state <> 'COMPLETED';
```

### List state changes for each assessment run

```sql
select
  name,
  arn,
  state,
  jsonb_pretty(state_changes) as state_changes
from
  aws_inspector_assessment_run;
```

### List assessment runs in the last 7 days

```sql
select
  name,
  arn,
  assessment_template_arn,
  created_at,
  state,
  region
from
  aws_inspector_assessment_run
where
  created_at >= (now() - interval '7' day);
```
