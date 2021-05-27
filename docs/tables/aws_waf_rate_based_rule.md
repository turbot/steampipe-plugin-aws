# Table: aws_waf_rate_based_rule

AWS WAF rate-based rules count the number of requests that arrive from a specified IP address every five minutes.

## Examples

### Basic info

```sql
select
  name,
  rule_id,
  metric_name
from
  aws_waf_rate_based_rule;
```


### List rate-based rules that allow a request based on the negation of the settings in predicates

```sql
select
  name,
  rule_id,
  p ->> 'DataId' as data_id,
  p ->> 'Negated' as negated,
  p ->> 'Type' as type
from
  aws_waf_rate_based_rule,
  jsonb_array_elements(predicates) as p
where
  p ->> 'Negated' = 'True';
```
