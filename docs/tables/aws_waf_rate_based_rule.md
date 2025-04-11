---
title: "Steampipe Table: aws_waf_rate_based_rule - Query AWS WAF RateBasedRule using SQL"
description: "Allows users to query AWS WAF RateBasedRule to retrieve information about rate-based security rules that AWS Web Application Firewall (WAF) uses."
folder: "WAF"
---

# Table: aws_waf_rate_based_rule - Query AWS WAF RateBasedRule using SQL

The AWS WAF RateBasedRule is a feature within AWS Web Application Firewall (WAF) service that helps protect your web applications or APIs against common web exploits. This rule allows you to specify the maximum number of requests that a client can make in a five-minute period. If the number of requests exceeds the specified limit, AWS WAF blocks further requests from the client.

## Table Usage Guide

The `aws_waf_rate_based_rule` table in Steampipe provides you with information about the rate-based security rules that AWS Web Application Firewall (WAF) uses to identify potentially malicious requests and manage how they are handled. This table allows you, as a security administrator, to query rule-specific details, including the rule ARN, creation and modification dates, associated metrics, and associated predicates. You can utilize this table to gather insights on rate-based rules, such as the number of requests that arrive from a single IP address over a five-minute period, the rule action (BLOCK or COUNT), and more. The schema outlines the various attributes of the AWS WAF rate-based rule for you, including the rule ID, metric name, rate limit, and associated tags.

## Examples

### Basic info
This query allows you to examine the metrics associated with different rate-based rules in your AWS Web Application Firewall. It can be particularly useful for understanding how these rules are performing and identifying potential areas for improvement.

```sql+postgres
select
  name,
  rule_id,
  metric_name
from
  aws_waf_rate_based_rule;
```

```sql+sqlite
select
  name,
  rule_id,
  metric_name
from
  aws_waf_rate_based_rule;
```


### List rate-based rules that allow a request based on the negation of the settings in predicates
This query is used to identify rate-based rules in AWS Web Application Firewall that permit requests based on the reversal of certain settings. This can be useful in pinpointing potential security vulnerabilities where requests are being allowed contrary to the intended configuration.

```sql+postgres
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

```sql+sqlite
select
  name,
  rule_id,
  json_extract(p.value, '$.DataId') as data_id,
  json_extract(p.value, '$.Negated') as negated,
  json_extract(p.value, '$.Type') as type
from
  aws_waf_rate_based_rule,
  json_each(predicates) as p
where
  json_extract(p.value, '$.Negated') = 'True';
```