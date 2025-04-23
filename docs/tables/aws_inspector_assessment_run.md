---
title: "Steampipe Table: aws_inspector_assessment_run - Query AWS Inspector Assessment Runs using SQL"
description: "Allows users to query AWS Inspector Assessment Runs to get detailed information about each assessment run, including its state, duration, findings, and more."
folder: "Inspector"
---

# Table: aws_inspector_assessment_run - Query AWS Inspector Assessment Runs using SQL

The AWS Inspector Assessment Run is a feature of AWS Inspector that allows you to evaluate the behavior of the applications you have in AWS against the defined set of AWS security best practices. It provides detailed findings about security vulnerabilities and deviations from best practices, with a detailed list of steps for remediation. This helps to improve the security and compliance of applications deployed on AWS.

## Table Usage Guide

The `aws_inspector_assessment_run` table in Steampipe provides you with information about assessment runs within AWS Inspector. This table allows you, as a DevOps engineer, to query run-specific details, including its state, duration, findings, and associated metadata. You can utilize this table to gather insights on runs, such as the number of findings, the state of the run, and the time it took for the run to complete. The schema outlines the various attributes of the assessment run for you, including the run ARN, creation date, state, duration, findings, and associated tags.

## Examples

### Basic info
Determine the areas in which AWS Inspector assessment runs are active and when they were created to better manage and monitor your AWS resources.

```sql+postgres
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

```sql+sqlite
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
This query is used to uncover the details of security assessment findings, categorized by their severity levels. It helps to prioritize necessary actions, by highlighting areas with high severity issues.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(finding_counts, '$.High') as high,
  json_extract(finding_counts, '$.Low') as low,
  json_extract(finding_counts, '$.Medium') as medium,
  json_extract(finding_counts, '$.Informational') as informational,
  state
from
  aws_inspector_assessment_run;
```

### List assessment runs for each assessment template
Identify instances where each assessment run corresponds to a specific assessment template. This can be useful for tracking the progress and status of different assessments, and for understanding the distribution of assessments across different regions.

```sql+postgres
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

```sql+sqlite
select
  t.name as assessment_template_name,
  r.name as assessment_run_name,
  r.created_at as assessment_run_created_at,
  r.state,
  r.region
from
  aws_inspector_assessment_run as r
join
  aws_inspector_assessment_template as t
on
  r.assessment_template_arn = t.arn;
```

### List assessment runs which are not completed
Identify instances where AWS Inspector assessment runs are still in progress. This can help in tracking the progress of security assessments and identifying any potential delays or issues.

```sql+postgres
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

```sql+sqlite
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
  state != 'COMPLETED';
```

### List state changes for each assessment run
Analyze the transitions of each assessment run to understand its progress and status changes over time. This can help in tracking the evolution and completion status of various assessments.

```sql+postgres
select
  name,
  arn,
  state,
  jsonb_pretty(state_changes) as state_changes
from
  aws_inspector_assessment_run;
```

```sql+sqlite
select
  name,
  arn,
  state,
  state_changes
from
  aws_inspector_assessment_run;
```

### List assessment runs in the last 7 days
Gain insights into recent security assessment runs within the past week. This is useful for understanding the current state and region of recent assessments, helping to maintain and improve security standards across your AWS resources.

```sql+postgres
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

```sql+sqlite
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
  created_at >= datetime('now', '-7 day');
```