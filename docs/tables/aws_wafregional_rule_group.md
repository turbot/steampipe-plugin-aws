# Table: aws_waf_rule_group

An AWS WAF Regional rule group is a collection of rules for inspecting and controlling web requests.

## Examples

### Basic info

```sql
select
  name,
  arn,
  rule_group_id,
  metric_name,
  activated_rules,
  region
from
  aws_wafregional_rule_group;
```

### List rule groups with no associated rules

```sql
select
  name,
  arn,
  rule_group_id,
  metric_name,
  activated_rules
from
  aws_wafregional_rule_group
where
  activated_rules is null or jsonb_array_length(activated_rules) = 0;
```

### List details of rules associated with the rule group

```sql
select
  name as rule_group_name,
  rule_group_id,
  a ->> 'RuleId' as rule_id,
  a -> 'Action' ->> 'Type' as rule_action_type,
  a ->> 'Type' as rule_type
from
  aws_wafregional_rule_group,
  jsonb_array_elements(activated_rules) as a;
```
