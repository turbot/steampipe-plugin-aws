---
title: "Steampipe Table: aws_networkfirewall_rule_group - Query AWS Network Firewall Rule Group using SQL"
description: "Allows users to query AWS Network Firewall Rule Group details, including rule group ARN, capacity, rule group name, and associated tags."
folder: "Network Firewall"
---

# Table: aws_networkfirewall_rule_group - Query AWS Network Firewall Rule Group using SQL

The AWS Network Firewall Rule Group is a component of AWS Network Firewall, a managed service that makes it easy to deploy essential network protections for all of your Amazon Virtual Private Clouds (VPCs). The rule group acts as a container for the stateless and stateful rule sets that make up the firewall policy for a network resource. It enables you to mix and match sets of rules to meet the specific security requirements of each individual resource.

## Table Usage Guide

The `aws_networkfirewall_rule_group` table in Steampipe provides you with information about rule groups within AWS Network Firewall. This table allows you, as a DevOps engineer, to query rule group-specific details, including the rule group ARN, capacity, rule group name, and associated tags. You can utilize this table to gather insights on rule groups, such as the rule group's capacity, the rule group's name, and more. The schema outlines for you the various attributes of the Network Firewall rule group, including the rule group ARN, capacity, rule group name, and associated tags.

## Examples

### Basic info
Explore the status and type of rule groups in your AWS Network Firewall to understand their configurations and ensure your network security measures are functioning as expected. This can be particularly useful in identifying areas of vulnerability or inefficiency within your firewall setup.

```sql+postgres
select
  rule_group_name,
  rule_group_status,
  type,
  jsonb_pretty(rules_source) as rules_source
from
  aws_networkfirewall_rule_group;
```

```sql+sqlite
select
  rule_group_name,
  rule_group_status,
  type,
  rules_source
from
  aws_networkfirewall_rule_group;
```

### List rule groups with no associations
Determine the areas in which rule groups within the AWS Network Firewall service are not associated with any entities. This can be useful in identifying unused rule groups that may be unnecessarily incurring costs or cluttering the system.

```sql+postgres
select
  rule_group_name,
  rule_group_status
from
  aws_networkfirewall_rule_group
where
  number_of_associations = 0;
```

```sql+sqlite
select
  rule_group_name,
  rule_group_status
from
  aws_networkfirewall_rule_group
where
  number_of_associations = 0;
```

### Get rules for stateful rule groups
This query is used to explore the rules for stateful rule groups in AWS Network Firewall. It's a useful tool for security administrators who want to analyze the status and options of these groups, providing insights into their configuration and potential vulnerabilities.

```sql+postgres
select
  rule_group_name,
  rule_group_status,
  jsonb_pretty(rules_source -> 'StatefulRules') as stateful_rules,
  jsonb_pretty(rule_variables) as rule_variables,
  stateful_rule_options
from
  aws_networkfirewall_rule_group
where
  type = 'STATEFUL';
```

```sql+sqlite
select
  rule_group_name,
  rule_group_status,
  json_extract(rules_source, '$.StatefulRules') as stateful_rules,
  rule_variables,
  stateful_rule_options
from
  aws_networkfirewall_rule_group
where
  type = 'STATEFUL';
```

### Get rules and custom actions for stateless rule groups
Determine the areas in which rules and custom actions apply for stateless rule groups. This information can be useful for understanding the configuration and status of your network firewall.

```sql+postgres
select
  rule_group_name,
  rule_group_status,
  jsonb_pretty(rules_source -> 'StatelessRulesAndCustomActions' -> 'StatelessRules') as stateless_rules,
  jsonb_pretty(rules_source -> 'StatelessRulesAndCustomActions' -> 'CustomActions') as custom_actions
from
  aws_networkfirewall_rule_group
where
  type = 'STATELESS';
```

```sql+sqlite
select
  rule_group_name,
  rule_group_status,
  json_extract(rules_source, '$.StatelessRulesAndCustomActions.StatelessRules') as stateless_rules,
  json_extract(rules_source, '$.StatelessRulesAndCustomActions.CustomActions') as custom_actions
from
  aws_networkfirewall_rule_group
where
  type = 'STATELESS';
```

### List rule groups with no rules
Determine the areas in your network firewall where rule groups are defined but contain no rules. This can help you identify potential vulnerabilities or inefficiencies in your network security setup.

```sql+postgres
select
  rule_group_name,
  rule_group_status,
  number_of_associations
from
  aws_networkfirewall_rule_group
where
  type = 'STATELESS' and jsonb_array_length(rules_source -> 'StatelessRulesAndCustomActions' -> 'StatelessRules') = 0
  or type = 'STATEFUL' and jsonb_array_length(rules_source -> 'StatefulRules') = 0;
```

```sql+sqlite
select
  rule_group_name,
  rule_group_status,
  number_of_associations
from
  aws_networkfirewall_rule_group
where
  (type = 'STATELESS' and json_array_length(json_extract(rules_source, '$.StatelessRulesAndCustomActions.StatelessRules')) = 0)
  or (type = 'STATEFUL' and json_array_length(json_extract(rules_source, '$.StatefulRules')) = 0);
```