# Table: aws_waf_rate_based_rule

AWS WAF RateBasedRule counts the number of requests that arrive from a specified IP address every five minutes.

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


### List ratebasedrules whose predicates are based on negation

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
