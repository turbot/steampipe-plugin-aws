---
title: "Table: aws_waf_rate_based_rule - Query AWS WAF RateBasedRule using SQL"
description: "Allows users to query AWS WAF RateBasedRule to retrieve information about rate-based security rules that AWS Web Application Firewall (WAF) uses."
---

# Table: aws_waf_rate_based_rule - Query AWS WAF RateBasedRule using SQL

The `aws_waf_rate_based_rule` table in Steampipe provides information about the rate-based security rules that AWS Web Application Firewall (WAF) uses to identify potentially malicious requests and manage how it handles them. This table allows security administrators to query rule-specific details, including the rule ARN, creation and modification dates, associated metrics, and associated predicates. Users can utilize this table to gather insights on rate-based rules, such as the number of requests that arrive from a single IP address over a five-minute period, the rule action (BLOCK or COUNT), and more. The schema outlines the various attributes of the AWS WAF rate-based rule, including the rule ID, metric name, rate limit, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_waf_rate_based_rule` table, you can use the `.inspect aws_waf_rate_based_rule` command in Steampipe.

**Key columns**:

- `rule_id`: The identifier for the rate-based rule. This can be used to join this table with other tables that contain information about AWS WAF rules.
- `name`: The name of the rate-based rule. This provides a human-readable identifier for the rule and can be used to join this table with other tables that contain information about AWS WAF rules.
- `metric_name`: The name or description for the Amazon CloudWatch metric for this rate-based rule. This can be used to join this table with CloudWatch metrics tables for further analysis and monitoring.

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
