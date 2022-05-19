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
  assessment_run_arn = 'arn:aws:inspector:us-east-1:533793682495:target/0-ywdTAdRg/template/0-rY1J4B4f/run/0-LRRwpQFz';
```

### Get the attributes for an exclusion

```sql
select
  arn,
  jsonb_pretty(attributes) as attributes
from
  aws_inspector_exclusion
where
  arn = 'arn:aws:inspector:us-east-1:533793682495:target/0-ywdTAdRg/template/0-rY1J4B4f/run/0-LRRwpQFz/exclusion/0-xNJPDc3o';
```

### Get the scopes for an exclusion

```sql
select
  arn,
  jsonb_pretty(attributes) as attributes
from
  aws_inspector_exclusion
where
  arn = 'arn:aws:inspector:us-east-1:533793682495:target/0-ywdTAdRg/template/0-rY1J4B4f/run/0-LRRwpQFz/exclusion/0-xNJPDc3o';
```