---
title: "Steampipe Table: aws_waf_rule - Query AWS WAF Rule using SQL"
description: "Allows users to query AWS Web Application Firewall (WAF) rules."
folder: "WAF"
---

# Table: aws_waf_rule - Query AWS WAF Rule using SQL

The AWS WAF Rule is a component of AWS Web Application Firewall (WAF) service. It allows you to protect your web applications from common web exploits that could affect application availability, compromise security, or consume excessive resources. AWS WAF gives you control over which traffic to allow or block to your web applications by defining customizable web security rules.

## Table Usage Guide

The `aws_waf_rule` table in Steampipe provides you with information about AWS WAF rules. These rules are used to block common web-based attacks. This table allows you, as a security professional or developer, to query rule-specific details, including the rule action (block, allow, or count), the predicates that make up the rule, and associated metadata. You can utilize this table to gather insights on rules, such as rules that are currently in effect, the conditions under which a rule is triggered, and more. The schema outlines the various attributes of the WAF rule for you, including the rule ID, type, metric name, and associated tags.

## Examples

### Basic info
This query allows you to analyze the rules associated with your AWS Web Application Firewall (WAF). It helps in understanding the effectiveness of your security measures by identifying the specific rules and metrics applied.

```sql+postgres
select
  name,
  rule_id,
  metric_name
from
  aws_waf_rule;
```

```sql+sqlite
select
  name,
  rule_id,
  metric_name
from
  aws_waf_rule;
```

### Get predicate details for each rule
Explore the specifics of each rule within AWS WAF, including whether the rule is negated and its type. This information can be useful for assessing the configuration and effectiveness of your web application firewall rules.

```sql+postgres
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

```sql+sqlite
select
  name,
  rule_id,
  json_extract(p.value, '$.DataId') as data_id,
  json_extract(p.value, '$.Negated') as negated,
  json_extract(p.value, '$.Type') as type
from
  aws_waf_rule,
  json_each(predicates) as p;
```