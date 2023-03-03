# Table: aws_emr_block_public_access_configuration

The `aws_emr_block_public_access_configuration` table provides information on the Amazon EMR block public access configuration for your AWS account in the current region.

## Examples

### Basic info

```sql
select
  created_by_arn,
  block_public_security_group_rules,
  creation_date,
  classification,
  configurations,
  permitted_public_security_group_rule_ranges
from
  aws_emr_block_public_access_configuration
order by
  created_by_arn,
  creation_date;
```

### List block public access configurations that block public security group rules

```sql
select
  created_by_arn,
  creation_date,
  configurations
from
  aws_emr_block_public_access_configuration
where
  block_public_security_group_rules;
```

### List permitted public security group rule maximum and minimum port ranges

```sql
select
  created_by_arn,
  creation_date,
  rules ->> 'MaxRange' as max_range,
  rules ->> 'MinRange' as min_range
from
  aws_emr_block_public_access_configuration
  cross join jsonb_array_elements(permitted_public_security_group_rule_ranges) as rules;
```

### List block public access configurations created in last 90 days

```sql
select
  created_by_arn,
  creation_date,
  configurations
from
  aws_emr_block_public_access_configuration
where
  date_part('day', now() - creation_date) < 90;
```
