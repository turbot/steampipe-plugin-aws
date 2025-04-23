---
title: "Steampipe Table: aws_ec2_load_balancer_listener_rule - Query AWS EC2 Load Balancer Listener Rules using SQL"
description: "Allows users to query AWS EC2 Load Balancer Listener Rules, providing detailed information on rule actions, conditions, and priority."
folder: "ELB"
---

# Table: aws_ec2_load_balancer_listener_rule - Query AWS EC2 Load Balancer Listener Rules using SQL

AWS Elastic Load Balancing automatically distributes incoming application traffic across multiple targets, such as EC2 instances. The listener rules determine how traffic is routed based on conditions specified by the user. The `aws_ec2_load_balancer_listener_rule` table in Steampipe allows you to query information about these listener rules, including their actions, conditions, priorities, and more.

## Table Usage Guide

The `aws_ec2_load_balancer_listener_rule` table enables cloud administrators and DevOps engineers to gather detailed insights into their load balancer listener rules. You can query various aspects of the rules, such as their actions, conditions, priorities, and associated listeners. This table is particularly useful for monitoring traffic routing, ensuring compliance with traffic rules, and managing load balancer configurations.

**Important Notes**
- You **_must_** specify `arn` or `listener_arn` in a `where` clause in order to use this table.

## Examples

### Basic info
Retrieve basic information about your AWS EC2 Load Balancer Listener Rules, including their ARN, priority, and associated listener.

```sql+postgres
select
  arn,
  listener_arn,
  priority,
  is_default,
  region
from
  aws_ec2_load_balancer_listener_rule
where
  arn = 'arn:aws:elasticloadbalancing:us-east-1:123456789098:listener-rule/app/test53333/f7cc8cdc44ff910b/c9418b57592205f0/a8fe6d8842838dfa';
```

```sql+sqlite
select
  arn,
  listener_arn,
  priority,
  is_default,
  region
from
  aws_ec2_load_balancer_listener_rule
where
  arn = 'arn:aws:elasticloadbalancing:us-east-1:123456789098:listener-rule/app/test53333/f7cc8cdc44ff910b/c9418b57592205f0/a8fe6d8842838dfa';
```

### List rules for a specific listener
Fetch all the rules associated with a specific listener by providing the listener ARN.

```sql+postgres
select
  arn,
  priority,
  is_default,
  actions,
  conditions
from
  aws_ec2_load_balancer_listener_rule
where
  listener_arn = 'arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-load-balancer/50dc6c495c0c9188/70d7923f8398b272';
```

```sql+sqlite
select
  arn,
  priority,
  is_default,
  actions,
  conditions
from
  aws_ec2_load_balancer_listener_rule
where
  listener_arn = 'arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-load-balancer/50dc6c495c0c9188/70d7923f8398b272';
```

### Get rule action details
Retrieve rules action details.

```sql+postgres
select
  arn,
  a ->> 'Type' as action_type,
  a ->> 'Order' as action_order,
  a ->> 'TargetGroupArn' as target_group_arn,
  a -> 'RedirectConfig' as redirect_config,
  a -> 'ForwardConfig' as forward_config,
  a -> 'FixedResponseConfig' as fixed_response_config,
  a -> 'AuthenticateOidcConfig' as authenticate_oidc_config,
  a -> 'AuthenticateCognitoConfig' as authenticate_cognito_config
from
  aws_ec2_load_balancer_listener_rule,
  jsonb_array_elements(actions) as a
where
  listener_arn = 'arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-load-balancer/50dc6c495c0c9188/70d7923f8398b272';
```

```sql+sqlite
select
  arn,
  json_extract(a.value, '$.Type') as action_type,
  json_extract(a.value, '$.Order') as action_order,
  json_extract(a.value, '$.TargetGroupArn') as target_group_arn,
  json_extract(a.value, '$.RedirectConfig') as redirect_config,
  json_extract(a.value, '$.ForwardConfig') as forward_config,
  json_extract(a.value, '$.FixedResponseConfig') as fixed_response_config,
  json_extract(a.value, '$.AuthenticateOidcConfig') as authenticate_oidc_config,
  json_extract(a.value, '$.AuthenticateCognitoConfig') as authenticate_cognito_config
from
  aws_ec2_load_balancer_listener_rule,
  json_each(actions) as a
where
  listener_arn = 'arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-load-balancer/50dc6c495c0c9188/70d7923f8398b272';
```

### List default listener rules
Identify the listener rules that are set as the default rule for a listener.

```sql+postgres
select
  arn,
  listener_arn,
  priority
from
  aws_ec2_load_balancer_listener_rule
where
  listener_arn = 'arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-load-balancer/50dc6c495c0c9188/70d7923f8398b272'
  and is_default = true;
```

```sql+sqlite
select
  arn,
  listener_arn,
  priority
from
  aws_ec2_load_balancer_listener_rule
where
  listener_arn = 'arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-load-balancer/50dc6c495c0c9188/70d7923f8398b272'
  and is_default = 1;
```

### Get all rules for listeners
Retrieve detailed information about the rules associated with AWS EC2 load balancer listeners.

```sql+postgres
select
  r.arn,
  r.listener_arn,
  l.load_balancer_arn,
  l.protocol as listener_protocol,
  l.ssl_policy,
  r.priority,
  r.is_default,
  r.actions,
  r.conditions
from
  aws_ec2_load_balancer_listener_rule as r
  join aws_ec2_load_balancer_listener as l on r.listener_arn = l.arn;
```

```sql+sqlite
select
  r.arn,
  r.listener_arn,
  l.load_balancer_arn,
  l.protocol as listener_protocol,
  l.ssl_policy,
  r.priority,
  r.is_default,
  r.actions,
  r.conditions
from
  aws_ec2_load_balancer_listener_rule as r
  join aws_ec2_load_balancer_listener as l on r.listener_arn = l.arn;
```

### Get listener rules for application load balancers
Retrieve detailed information about the rules associated with AWS EC2 Application load balancer listeners.

```sql+postgres
select
  r.arn,
  r.listener_arn,
  l.load_balancer_arn,
  l.protocol as listener_protocol,
  l.ssl_policy,
  a.canonical_hosted_zone_id,
  a.dns_name,
  a.ip_address_type,
  r.priority,
  r.is_default,
  r.actions,
  r.conditions
from
  aws_ec2_load_balancer_listener_rule as r
  join aws_ec2_load_balancer_listener as l on r.listener_arn = l.arn
  join aws_ec2_application_load_balancer as a on l.load_balancer_arn = a.arn;
```

```sql+sqlite
 select
  r.arn,
  r.listener_arn,
  l.load_balancer_arn,
  l.protocol as listener_protocol,
  l.ssl_policy,
  a.canonical_hosted_zone_id,
  a.dns_name,
  a.ip_address_type,
  r.priority,
  r.is_default,
  r.actions,
  r.conditions
from
  aws_ec2_load_balancer_listener_rule as r
  join aws_ec2_load_balancer_listener as l on r.listener_arn = l.arn
  join aws_ec2_application_load_balancer as a on l.load_balancer_arn = a.arn;
```