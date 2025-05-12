---
title: "Steampipe Table: aws_wafregional_rule_group - Query AWS WAF Regional Rule Groups using SQL"
description: "Allows users to query AWS WAF Regional Rule Groups to gather information about each rule group's metadata, associated rules, and other relevant details."
folder: "Region"
---

# Table: aws_wafregional_rule_group - Query AWS WAF Regional Rule Groups using SQL

The AWS WAF Regional Rule Groups are a feature of the AWS WAF service that allows you to categorize and manage similar rules. These groups are used to consolidate rules and simplify the process of adding multiple rules to a web ACL. Rule groups help in enhancing security by enabling you to specify which AWS resources are in scope for a rule, thereby restricting access and reducing potential threats.

## Table Usage Guide

The `aws_wafregional_rule_group` table in Steampipe provides you with information about rule groups within AWS WAF Regional. This table allows you, as a DevOps engineer, to query rule group-specific details, including the rule group ARN, associated rules, and metadata. You can utilize this table to gather insights on rule groups, such as the activated rules in each group, the metric names associated with each rule, and more. The schema outlines the various attributes of the rule group for you, including the rule group ID, name, ARN, metric name, and associated tags.

## Examples

### Basic info
Explore the configuration of AWS WAF regional rule groups to understand the security measures in place across different regions. This can be useful for auditing security protocols and identifying potential areas for improvement.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in your AWS security setup where rule groups lack associated rules, allowing you to identify potential vulnerabilities and improve your overall security posture.

```sql+postgres
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

```sql+sqlite
select
  name,
  arn,
  rule_group_id,
  metric_name,
  activated_rules
from
  aws_wafregional_rule_group
where
  activated_rules is null or json_array_length(activated_rules) = 0;
```

### List details of rules associated with the rule group
Explore the specifics of rules linked to a particular rule group in AWS WAF Regional. This can help you understand the nature and function of each rule, aiding in security management and threat mitigation.

```sql+postgres
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

```sql+sqlite
select
  name as rule_group_name,
  rule_group_id,
  json_extract(a.value, '$.RuleId') as rule_id,
  json_extract(json_extract(a.value, '$.Action'), '$.Type') as rule_action_type,
  json_extract(a.value, '$.Type') as rule_type
from
  aws_wafregional_rule_group,
  json_each(activated_rules) as a;
```