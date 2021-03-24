# Table: aws_inspector_assessment_template

The AWS Inspector Assessment Template resource creates an Amazon Inspector assessment template, which specifies the Inspector assessment targets that will be evaluated by an assessment run and its related configurations.

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


### List Assessment Templates that have no user attributes for findings tags

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


### List Assessment Templates that have zero assessment runs

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


### List Assessment Templates with assessment run duration less than 1 hour

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


### List Assessment Templates that are created within last 7 days

```sql
select
  name,
  arn,
  assessment_run_count,
  created_at,
  assessment_target_arn,
  region
from
  aws_inspector_assessment_template
where
  created_at > (current_date - interval '7' day);
```