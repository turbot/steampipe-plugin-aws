---
title: "Table: aws_waf_rule - Query AWS WAF Rule using SQL"
description: "Allows users to query AWS Web Application Firewall (WAF) rules."
---

# Table: aws_waf_rule - Query AWS WAF Rule using SQL

The `aws_waf_rule` table in Steampipe provides information about AWS WAF rules. These rules are used to block common web-based attacks. The table allows security professionals and developers to query rule-specific details, including the rule action (block, allow, or count), the predicates that make up the rule, and associated metadata. Users can utilize this table to gather insights on rules, such as rules that are currently in effect, the conditions under which a rule is triggered, and more. The schema outlines the various attributes of the WAF rule, including the rule ID, type, metric name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_waf_rule` table, you can use the `.inspect aws_waf_rule` command in Steampipe.

**Key columns**:

- `rule_id`: The identifier for the rule. This can be used to join this table with other tables that reference WAF rules.
- `name`: The name of the rule. This can be useful for human-readable queries and joins with other tables that reference WAF rules by name.
- `type`: The type of the rule (REGULAR, RATE_BASED, or GROUP). This can be useful for filtering or categorizing rules.

## Examples

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