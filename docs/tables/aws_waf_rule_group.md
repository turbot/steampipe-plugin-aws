---
title: "Table: aws_waf_rule_group - Query AWS WAF Rule Groups using SQL"
description: "Allows users to query AWS WAF Rule Groups to provide information about Web Application Firewall (WAF) rule groups within AWS WAF. This table enables security and DevOps engineers to query rule group-specific details, including rules, actions, and associated metadata."
---

# Table: aws_waf_rule_group - Query AWS WAF Rule Groups using SQL

The `aws_waf_rule_group` table in Steampipe provides information about Web Application Firewall (WAF) rule groups within AWS WAF. This table allows security and DevOps engineers to query rule group-specific details, including the rule group ID, name, metric name, and associated rules. Users can utilize this table to gather insights on rule groups, such as the types of rules within a rule group, the actions for each rule, and more. The schema outlines the various attributes of the WAF rule group, including the rule group ID, name, metric name, and associated rules.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_waf_rule_group` table, you can use the `.inspect aws_waf_rule_group` command in Steampipe.

**Key columns**:

- `name`: The name of the rule group. This column is useful for joining with other tables that reference rule groups by name.
- `rule_group_id`: The identifier of the rule group. This column can be used to join with other tables that reference rule groups by their ID.
- `metric_name`: The name of the metric for the rule group. This column can be used to join with other tables that reference rule groups by their metric name.

## Examples

### Basic info

```sql
select
  name,
  arn,
  rule_group_id,
  metric_name,
  activated_rules
from
  aws_waf_rule_group;
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
  aws_waf_rule_group
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
  aws_waf_rule_group,
  jsonb_array_elements(activated_rules) as a;
```
