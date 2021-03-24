# Table: aws_inspector_assessment_target

The AWS Inspector Assessment Target resource specify the Amazon EC2 instances that will be analyzed during an assessment run.

## Examples

### Basic info

```sql
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target;
```


### List assessment targets that are created within last 7 days

```sql
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target
where
  created_at > (current_date - interval '7' day);
```


### List assessment targets that are updated after creation

```sql
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target
where
  created_at != updated_at;
```