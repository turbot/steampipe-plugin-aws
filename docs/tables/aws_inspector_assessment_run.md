# Table: aws_inspector_assessment_run

The AWS Inspector Assessment Template resource specifies the Inspector assessment targets that will be evaluated by an assessment run and its related configurations. You use the template to start an assessment run, which is the monitoring and analysis process that results in a set of findings.

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
