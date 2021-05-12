# Table: aws_inspector_assessment_template

The AWS Inspector Assessment Template resource specifies the Inspector assessment targets that will be evaluated by an assessment run and its related configurations.

## Examples

### Basic info

```sql
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  region
from
  aws_inspector_assessment_template;
```


### List assessment templates that have no assigned finding attributes

```sql
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  user_attributes_for_findings = '[]';
```


### List assessment templates that have no assessment runs

```sql
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  user_attributes_for_findings,
  region
from
  aws_inspector_assessment_template
where
  assessment_run_count = 0;
```


### List assessment templates with run duration less than 1 hour

```sql
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  duration_in_seconds,
  region
from
  aws_inspector_assessment_template
where
  duration_in_seconds < 3600;
```
