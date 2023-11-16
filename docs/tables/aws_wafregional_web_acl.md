---
title: "Table: aws_wafregional_web_acl - Query AWS WAF Regional WebACL using SQL"
description: "Allows users to query AWS WAF Regional WebACL to get information about AWS WAF Regional WebACLs."
---

# Table: aws_wafregional_web_acl - Query AWS WAF Regional WebACL using SQL

The `aws_wafregional_web_acl` table in Steampipe provides information about Web Access Control Lists (WebACLs) in AWS WAF Regional. This table allows security professionals to query WebACL-specific details, including associated rules, default actions, metric names, and associated metadata. Users can utilize this table to gather insights on WebACLs, such as their associated rules, default actions, and more. The schema outlines the various attributes of the WebACL, including the WebACL ID, ARN, name, metric name, default action, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wafregional_web_acl` table, you can use the `.inspect aws_wafregional_web_acl` command in Steampipe.

**Key columns**:

- `web_acl_id`: The identifier of the WebACL. This can be used to join this table with other tables to get more information about the WebACL.
- `name`: The name of the WebACL. This can be used to filter the results to get information about a specific WebACL.
- `default_action`: The default action for the WebACL. This can be used to understand the default behavior when a request does not match any rules in the WebACL.

## Examples

### Basic info

```sql
select
  name,
  web_acl_id,
  arn,
  region,
  default_action,
  tags
from
  aws_wafregional_web_acl;
```

### Get rule details for each web ACL

```sql
select
  name,
  web_acl_id,
  r ->> 'RuleId' as rule_id,
  r ->> 'Type' as rule_type,
  r ->> 'ExcludedRules' as excluded_rules,
  r ->> 'OverrideAction' as override_action,
  r -> 'Action' ->> 'Type' as action_type
from
  aws_wafregional_web_acl,
  jsonb_array_elements(rules) as r;
```

### Get web ACLs with no rules defined

```sql
select
  name,
  web_acl_id,
  arn,
  region,
  default_action,
  tags
from
  aws_wafregional_web_acl
where
  rules is null;
```

### Get web ACLs with default action as 'ALLOW'

```sql
select
  name,
  web_acl_id,
  arn,
  region,
  default_action
from
  aws_wafregional_web_acl
where
  default_action = 'ALLOW';
```

### List web ACLs with logging disabled

```sql
select
  name,
  web_acl_id,
  arn,
  region
from
  aws_wafregional_web_acl
where
  logging_configuration is null;
```
