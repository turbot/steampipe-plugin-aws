---
title: "Steampipe Table: aws_wafregional_web_acl - Query AWS WAF Regional WebACL using SQL"
description: "Allows users to query AWS WAF Regional WebACL to get information about AWS WAF Regional WebACLs."
folder: "Region"
---

# Table: aws_wafregional_web_acl - Query AWS WAF Regional WebACL using SQL

The AWS WAF Regional WebACL is a resource within the AWS WAF service that allows you to protect your AWS resources against common web exploits. It provides control over which traffic to allow or block to your web applications by defining customizable web security rules. You can use AWS WAF to create custom rules that block common attack patterns, such as SQL injection or cross-site scripting.

## Table Usage Guide

The `aws_wafregional_web_acl` table in Steampipe provides you with information about Web Access Control Lists (WebACLs) in AWS WAF Regional. This table allows you, as a security professional, to query WebACL-specific details, including associated rules, default actions, metric names, and associated metadata. You can utilize this table to gather insights on WebACLs, such as their associated rules, default actions, and more. The schema outlines for you the various attributes of the WebACL, including the WebACL ID, ARN, name, metric name, default action, and associated tags.

## Examples

### Basic info
Explore the default actions and regional distribution of your AWS WAF web access control lists (ACLs) to gain insights into their configurations and associated tags. This is useful in identifying potential security gaps and ensuring that your ACLs are optimally configured for your specific use case.

```sql+postgres
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

```sql+sqlite
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
Explore specific rules applied to each web application firewall (WAF) to understand its function, including any rules that have been excluded and any actions taken when a rule is triggered. This can help in assessing the security configuration of your WAF.

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
  aws_wafregional_web_acl,
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
  aws_wafregional_web_acl,
  json_each(rules) as r;
```

### Get web ACLs with no rules defined
Identify instances where web access control lists (ACLs) are defined without any rules in the AWS WAF Regional service. This can help in pinpointing potential security vulnerabilities where traffic is not being properly filtered.

```sql+postgres
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

```sql+sqlite
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
Explore which web access control lists (ACLs) in your AWS WAF regional setup have been configured to allow all traffic by default. This can help you identify potential security vulnerabilities where access is not sufficiently restricted.

```sql+postgres
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

```sql+sqlite
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
Explore which web access control lists (ACLs) in your AWS infrastructure have logging disabled. This is useful for identifying potential security blind spots where unauthorized access could occur undetected.

```sql+postgres
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

```sql+sqlite
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