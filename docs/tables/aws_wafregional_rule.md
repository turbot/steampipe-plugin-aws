---
title: "Steampipe Table: aws_wafregional_rule - Query AWS WAF Regional Rules using SQL"
description: "Allows users to query AWS WAF Regional Rules for detailed information about each rule, including its ID, metric name, name, and the predicates associated with it."
folder: "Region"
---

# Table: aws_wafregional_rule - Query AWS WAF Regional Rules using SQL

The AWS WAF Regional Rule is a feature of AWS WAF, a web application firewall that helps protect your web applications from common web exploits. It allows you to create custom rules that block common attack patterns, such as SQL injection or cross-site scripting (XSS), and rules that are designed for your specific application. These rules can be used in AWS WAF to block or allow requests based on conditions that you specify.

## Table Usage Guide

The `aws_wafregional_rule` table in Steampipe provides you with information about AWS WAF Regional Rules. This table allows you as a DevOps engineer to query rule-specific details, including its ID, metric name, name, and the predicates associated with it. You can utilize this table to gather insights on rules, such as the types of patterns that AWS WAF searches for, whether AWS WAF is set to allow, block, or count web requests, and more. The schema outlines the various attributes of the AWS WAF Regional Rule for you, including the rule ARN, rule ID, metric name, and associated predicates.

## Examples

### Basic info
Determine the areas in which specific rules and associated metrics are applied within your AWS WAF regional setup. This can help you understand the reach and impact of your security configurations.

```sql+postgres
select
  name,
  rule_id,
  metric_name
from
  aws_wafregional_rule;
```

```sql+sqlite
select
  name,
  rule_id,
  metric_name
from
  aws_wafregional_rule;
```

### Get predicate details for each rule
Determine the specifics of each rule in your AWS WAF Regional setup, including whether conditions are negated and the type of data being evaluated. This allows for a comprehensive review of your security settings, helping identify potential weak points or areas for improvement.

```sql+postgres
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

```sql+sqlite
select
  name,
  rule_id,
  json_extract(p.value, '$.DataId') as data_id,
  json_extract(p.value, '$.Negated') as negated,
  json_extract(p.value, '$.Type') as type
from
  aws_wafregional_rule,
  json_each(predicates) as p;
```