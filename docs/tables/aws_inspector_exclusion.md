# Table: aws_inspector_exclusion

The AWS Inspector Exclusions are an output of assessment runs. Exclusions show which of your security checks can't be completed and how to resolve the issues.

## Examples

### Basic info

```sql
select
  arn,
  attributes,
  description,
  title,
  region
from
  aws_inspector_exclusion;
```

### List exclusions associated with an assessment run

```sql
select
  arn,
  attributes,
  description,
  title,
  region
from
  aws_inspector_exclusion
where
  assessment_run_arn = 'arn:aws:inspector:us-east-1:012345678912:target/0-ywdTAdRg/template/0-rY1J4B4f/run/0-LRRwpQFz';
```

### Get the attribute and scope details for each exclusion

```sql
select
  arn,
  jsonb_pretty(attributes) as attributes,
  jsonb_pretty(scopes) as scopes
from
  aws_inspector_exclusion;
```

### Count the number of exclusions whose type is 'Agent not found'

```sql
select
  arn,
  region,
  title,
  count(arn)
from
  aws_inspector_exclusion
group by
  arn,
  region,
  title
order by
  count desc;
```

### Get the exclusion details of each assessment template that have run at least once

```sql
select 
  e.arn, 
  e.title, 
  jsonb_pretty(e.attributes) as attributes, 
  e.recommendation 
from 
  test_aab.aws_inspector_exclusion e, 
  test_aab.aws_inspector_assessment_run r, 
  test_aab.aws_inspector_assessment_template t 
where 
  e.assessment_run_arn = r.arn 
and 
  r.assessment_template_arn = t.arn;
```