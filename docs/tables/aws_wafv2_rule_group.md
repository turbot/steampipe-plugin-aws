# Table: aws_wafv2_rule_group

An AWS WAFv2 Rule Group a collection of rules for inspecting and controlling web requests.

## Examples

### Basic info

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  capacity,
  rules,
  region
from
  aws_wafv2_rule_group;
```

### List global (CloudFront) Rule Groups

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  capacity,
  region
from
  aws_wafv2_rule_group
where
  scope = 'CLOUDFRONT';
```

### List Rule Groups having having capacity less than 5

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  capacity,
  region
from
  aws_wafv2_rule_group
where
  capacity < 5;
```
