# Table: aws_waf_web_acl

A web access control list (web ACL) gives you fine-grained control over all of the HTTP(S) web requests that your protected resource responds to. You can protect Amazon CloudFront, Amazon API Gateway, Application Load Balancer, and AWS AppSync resources.

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
