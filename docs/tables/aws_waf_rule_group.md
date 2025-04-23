---
title: "Steampipe Table: aws_waf_rule_group - Query AWS WAF Rule Groups using SQL"
description: "Allows users to query AWS WAF Rule Groups to provide information about Web Application Firewall (WAF) rule groups within AWS WAF. This table enables security and DevOps engineers to query rule group-specific details, including rules, actions, and associated metadata."
folder: "WAF"
---

# Table: aws_waf_rule_group - Query AWS WAF Rule Groups using SQL

The AWS WAF Rule Group is a component of the AWS Web Application Firewall (WAF) service that allows you to bundle rules that identify common patterns of malicious web requests. These rule groups can then be associated with resources to protect them from these identified threats. This facilitates the management and organization of security rules, improving the overall protection of your web applications.

## Table Usage Guide

The `aws_waf_rule_group` table in Steampipe provides you with information about Web Application Firewall (WAF) rule groups within AWS WAF. This table allows you, as a security or DevOps engineer, to query rule group-specific details, including the rule group ID, name, metric name, and associated rules. You can utilize this table to gather insights on rule groups, such as the types of rules within a rule group, the actions for each rule, and more. The schema outlines the various attributes of the WAF rule group for you, including the rule group ID, name, metric name, and associated rules.

## Examples

### Basic info
Analyze the settings to understand the activated rules of your AWS WAF rule groups. This can help you gain insights into the security measures in place and identify any potential areas for improvement.

```sql+postgres
select
  name,
  arn,
  rule_group_id,
  metric_name,
  activated_rules
from
  aws_waf_rule_group;
```

```sql+sqlite
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
Discover the segments that have rule groups with no associated rules in AWS WAF. This is useful in identifying potential security gaps, as these rule groups are not actively filtering web traffic.

```sql+postgres
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

```sql+sqlite
select
  name,
  arn,
  rule_group_id,
  metric_name,
  activated_rules
from
  aws_waf_rule_group
where
  activated_rules is null or json_array_length(activated_rules) = 0;
```

### List details of rules associated with the rule group
Identify the specific rules linked to a particular rule group in AWS WAF. This can be useful in understanding the security actions and types associated with each rule within the group.

```sql+postgres
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

```sql+sqlite
select
  name as rule_group_name,
  rule_group_id,
  json_extract(a.value, '$.RuleId') as rule_id,
  json_extract(json_extract(a.value, '$.Action'), '$.Type') as rule_action_type,
  json_extract(a.value, '$.Type') as rule_type
from
  aws_waf_rule_group,
  json_each(activated_rules) as a;
```