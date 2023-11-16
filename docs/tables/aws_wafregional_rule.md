---
title: "Table: aws_wafregional_rule - Query AWS WAF Regional Rules using SQL"
description: "Allows users to query AWS WAF Regional Rules for detailed information about each rule, including its ID, metric name, name, and the predicates associated with it."
---

# Table: aws_wafregional_rule - Query AWS WAF Regional Rules using SQL

The `aws_wafregional_rule` table in Steampipe provides information about AWS WAF Regional Rules. This table allows DevOps engineers to query rule-specific details, including its ID, metric name, name, and the predicates associated with it. Users can utilize this table to gather insights on rules, such as the types of patterns that AWS WAF searches for, whether AWS WAF is set to allow, block, or count web requests, and more. The schema outlines the various attributes of the AWS WAF Regional Rule, including the rule ARN, rule ID, metric name, and associated predicates.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wafregional_rule` table, you can use the `.inspect aws_wafregional_rule` command in Steampipe.

**Key columns**:

- `rule_id`: The identifier for the rule. This can be used to join this table with other tables that contain rule-specific information.
- `name`: The name of the rule. This can be useful for joining with tables that reference rules by name.
- `metric_name`: The name or description for the Amazon CloudWatch metric for this rule. This can be used to join with CloudWatch metric tables for additional insights.

## Examples

### Basic info

```sql
select
  name,
  rule_id,
  metric_name
from
  aws_wafregional_rule;
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
  aws_wafregional_rule,
  jsonb_array_elements(predicates) as p;
```