# Table: aws_networkfirewall_rule_group

A rule group to inspect and control network traffic. You define stateless rule groups to inspect individual packets and you define stateful rule groups to inspect packets in the context of their traffic flow.

## Examples

### Basic info

```sql
select
  rule_group_name,
  rule_group_status,
  type,
  jsonb_pretty(rules_source) as rules_source
from
  aws_networkfirewall_rule_group;
```

### List rule groups with no associations

```sql
select
  rule_group_name,
  rule_group_status
from
  aws_networkfirewall_rule_group
where
  number_of_associations = 0;
```

### Get rules for stateful rule groups

```sql
select
  rule_group_name,
  rule_group_status,
  jsonb_pretty(rules_source -> 'StatefulRules') as stateful_rules,
  jsonb_pretty(rule_variables) as rule_variables,
  stateful_rule_options
from
  aws_networkfirewall_rule_group
where
  type = 'STATEFUL';
```

### Get rules and custom actions for stateless rule groups

```sql
select
  rule_group_name,
  rule_group_status,
  jsonb_pretty(rules_source -> 'StatelessRulesAndCustomActions' -> 'StatelessRules') as stateless_rules,
  jsonb_pretty(rules_source -> 'StatelessRulesAndCustomActions' -> 'CustomActions') as custom_actions
from
  aws_networkfirewall_rule_group
where
  type = 'STATELESS';
```

### List rule groups with no rules

```sql
select
  rule_group_name,
  rule_group_status,
  number_of_associations
from
  aws_networkfirewall_rule_group
where
  type = 'STATELESS' and jsonb_array_length(rules_source -> 'StatelessRulesAndCustomActions' -> 'StatelessRules') = 0
  or type = 'STATEFUL' and jsonb_array_length(rules_source -> 'StatefulRules') = 0;
```
