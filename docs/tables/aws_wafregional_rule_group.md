---
title: "Table: aws_wafregional_rule_group - Query AWS WAF Regional Rule Groups using SQL"
description: "Allows users to query AWS WAF Regional Rule Groups to gather information about each rule group's metadata, associated rules, and other relevant details."
---

# Table: aws_wafregional_rule_group - Query AWS WAF Regional Rule Groups using SQL

The `aws_wafregional_rule_group` table in Steampipe provides information about rule groups within AWS WAF Regional. This table allows DevOps engineers to query rule group-specific details, including the rule group ARN, associated rules, and metadata. Users can utilize this table to gather insights on rule groups, such as the activated rules in each group, the metric names associated with each rule, and more. The schema outlines the various attributes of the rule group, including the rule group ID, name, ARN, metric name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wafregional_rule_group` table, you can use the `.inspect aws_wafregional_rule_group` command in Steampipe.

**Key columns**:

- `rule_group_id`: The identifier for the rule group. This can be used to join this table with other tables that contain rule group ID information.
- `name`: The name of the rule group. This is useful for querying specific rule groups by their name.
- `arn`: The Amazon Resource Name (ARN) of the rule group. This can be used to join this table with any other table that contains ARN information for AWS resources.

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
