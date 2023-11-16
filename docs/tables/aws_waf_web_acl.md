---
title: "Table: aws_waf_web_acl - Query AWS WAF WebACLs using SQL"
description: "Allows users to query AWS WAF WebACLs to retrieve information about their configuration, rules, and associated metadata."
---

# Table: aws_waf_web_acl - Query AWS WAF WebACLs using SQL

The `aws_waf_web_acl` table in Steampipe provides information about Web Access Control Lists (WebACLs) within AWS WAF. This table allows security engineers to query WebACL-specific details, including associated rules, actions, and metadata. Users can utilize this table to gather insights on WebACLs, such as what rules are applied, what actions are taken when a rule is matched, and more. The schema outlines the various attributes of the WebACL, including the WebACL ARN, ID, default action, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_waf_web_acl` table, you can use the `.inspect aws_waf_web_acl` command in Steampipe.

### Key columns:

- `web_acl_id`: The identifier for the WebACL. This can be used to join this table with other tables that contain WebACL IDs.
- `name`: The name of the WebACL. This can be useful for filtering or ordering results based on the WebACL name.
- `default_action`: The action that AWS WAF takes when a web request does not match any of the rules in a WebACL. This can be used to identify the default behavior of the WebACL.

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
  aws_waf_web_acl;
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
  aws_waf_web_acl,
  jsonb_array_elements(rules) as r;
```

### Get web ACLs with no rule defined

```sql
select
  name,
  web_acl_id,
  arn,
  region,
  default_action,
  tags
from
  aws_waf_web_acl
where
  rules is null;
```

### Get web ACLs with default action as allow

```sql
select
  name,
  web_acl_id,
  arn,
  region,
  default_action
from
  aws_waf_web_acl
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
  aws_waf_web_acl
where
  logging_configuration is null;
```
