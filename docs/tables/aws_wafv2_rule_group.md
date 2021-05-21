# Table: aws_wafv2_rule_group

An AWS WAFv2 rule group is a collection of rules for inspecting and controlling web requests.

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

### List global (CloudFront) rule groups

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

### List rule groups with fewer than 5 web ACL capacity units (WCUs)

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
