# Table: aws_waf_rule

Amazon WAF Rule defines how to inspect web requests and what to do when a web request matches the inspection criteria.

## Example

### Basic info

```sql
select
  name,
  rule_id,
  metric_name
from
  aws_waf_rule;
```

### Get predicate details for each rule

```sql
select
  name,
  rule_id,
  p ->> 'DataId' as data_id,
  p ->> 'Negated' as negated,
  p ->> 'Type' as type
from
  aws_waf_rule,
  jsonb_array_elements(predicates) as p;
```
