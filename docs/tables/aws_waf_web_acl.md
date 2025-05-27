---
title: "Steampipe Table: aws_waf_web_acl - Query AWS WAF WebACLs using SQL"
description: "Allows users to query AWS WAF WebACLs to retrieve information about their configuration, rules, and associated metadata."
folder: "WAF"
---

# Table: aws_waf_web_acl - Query AWS WAF WebACLs using SQL

The AWS WAF WebACL is a resource in AWS WAF service that provides control over how AWS WAF handles a request for a web application. It contains a set of rules that dictate which traffic to allow, block, or count. These rules can be based on IP addresses, HTTP headers, HTTP body, or URI strings, providing a flexible and powerful security layer for your web applications.

## Table Usage Guide

The `aws_waf_web_acl` table in Steampipe provides you with information about Web Access Control Lists (WebACLs) within AWS WAF. This table allows you, as a security engineer, to query WebACL-specific details, including associated rules, actions, and metadata. You can utilize this table to gather insights on WebACLs, such as what rules are applied, what actions are taken when a rule is matched, and more. The schema outlines the various attributes of the WebACL for you, including the WebACL ARN, ID, default action, and associated tags.

## Examples

### Basic info
Explore the settings of your AWS WAF web access control list (ACL) to understand its default actions and associated regions. This is useful for assessing the security configuration of your AWS resources and identifying potential areas for improvement.

```sql+postgres
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

```sql+sqlite
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
Determine the specifics of each web access control list (ACL) rule, including its type, any excluded rules, and its action type. This can help in understanding the security configuration and identifying any potential vulnerabilities or areas for improvement.

```sql+postgres
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

```sql+sqlite
select
  name,
  web_acl_id,
  json_extract(r.value, '$.RuleId') as rule_id,
  json_extract(r.value, '$.Type') as rule_type,
  json_extract(r.value, '$.ExcludedRules') as excluded_rules,
  json_extract(r.value, '$.OverrideAction') as override_action,
  json_extract(json_extract(r.value, '$.Action'), '$.Type') as action_type
from
  aws_waf_web_acl,
  json_each(rules) as r;
```

### Get web ACLs with no rule defined
Identify instances where web access control lists (ACLs) have no defined rules. This is beneficial in pinpointing potential security gaps in your AWS WAF configuration.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which web access control lists (ACLs) are set to allow by default. This is useful to identify potential security vulnerabilities where unrestricted access is granted.

```sql+postgres
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

```sql+sqlite
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
This example helps identify web access control lists (ACLs) within your AWS infrastructure that have logging disabled. This can be useful in enhancing security measures by ensuring all web ACLs have logging enabled for better tracking and auditing.

```sql+postgres
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

```sql+sqlite
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