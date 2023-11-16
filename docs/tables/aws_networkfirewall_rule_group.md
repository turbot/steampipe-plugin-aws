---
title: "Table: aws_networkfirewall_rule_group - Query AWS Network Firewall Rule Group using SQL"
description: "Allows users to query AWS Network Firewall Rule Group details, including rule group ARN, capacity, rule group name, and associated tags."
---

# Table: aws_networkfirewall_rule_group - Query AWS Network Firewall Rule Group using SQL

The `aws_networkfirewall_rule_group` table in Steampipe provides information about rule groups within AWS Network Firewall. This table allows DevOps engineers to query rule group-specific details, including the rule group ARN, capacity, rule group name, and associated tags. Users can utilize this table to gather insights on rule groups, such as the rule group's capacity, the rule group's name, and more. The schema outlines the various attributes of the Network Firewall rule group, including the rule group ARN, capacity, rule group name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_networkfirewall_rule_group` table, you can use the `.inspect aws_networkfirewall_rule_group` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Name (ARN) of the rule group. This can be used to join this table with other tables.
- `capacity`: The capacity setting of the rule group. This can be useful for understanding the scale of the rule group.
- `rule_group_name`: The name of the rule group. This can be useful for identifying specific rule groups in queries.

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
