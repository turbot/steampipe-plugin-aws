---
title: "Steampipe Table: aws_ec2_load_balancer_listener - Query AWS EC2 Load Balancer Listeners using SQL"
description: "Allows users to query AWS EC2 Load Balancer Listener data, which provides information about listeners for an Application Load Balancer or Network Load Balancer."
folder: "ELB"
---

# Table: aws_ec2_load_balancer_listener - Query AWS EC2 Load Balancer Listeners using SQL

An AWS EC2 Load Balancer Listener is a component of the AWS Elastic Load Balancing service that checks for connection requests. It is configured with a protocol and port for the front-end (client to load balancer) connections, and a protocol and port for the back-end (load balancer to back-end instance) connections. Listeners are crucial in routing requests from clients to the registered instances based on the configured routing policies.

## Table Usage Guide

The `aws_ec2_load_balancer_listener` table in Steampipe provides you with information about listeners for an Application Load Balancer or Network Load Balancer in Amazon Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query listener-specific details, including protocol, port, SSL policy, and associated actions. You can utilize this table to gather insights on listeners, such as their current state, default actions, and certificates. The schema outlines the various attributes of the Load Balancer Listener for you, including the listener ARN, load balancer ARN, default actions, and associated tags.

## Examples

### Load balancer listener basic info
Determine the areas in which your AWS EC2 load balancer listeners operate, by examining crucial details such as port and protocol. This information can be beneficial in optimizing network traffic management and troubleshooting connectivity issues.

```sql+postgres
select
  title,
  arn,
  port,
  protocol
from
  aws_ec2_load_balancer_listener;
```

```sql+sqlite
select
  title,
  arn,
  port,
  protocol
from
  aws_ec2_load_balancer_listener;
```


### Action configuration details of each load balancer
Explore the configuration details of each load balancer's actions to understand how they are set up for authentication, fixed responses, and target group stickiness. This can be useful in assessing the security and efficiency of your load balancing setup.

```sql+postgres
select
  title,
  arn,
  action ->> 'AuthenticateCognitoConfig' as authenticate_cognito_config,
  action ->> 'AuthenticateOidcConfig' as authenticate_Oidc_config,
  action ->> 'FixedResponseConfig' as fixed_response_config,
  action -> 'ForwardConfig' -> 'TargetGroupStickinessConfig' ->> 'DurationSeconds' as duration_seconds,
  action -> 'ForwardConfig' -> 'TargetGroupStickinessConfig' ->> 'Enabled' as target_group_stickiness_config_enabled
from
  aws_ec2_load_balancer_listener
  cross join jsonb_array_elements(default_actions) as action;
```

```sql+sqlite
select
  title,
  arn,
  json_extract(action.value, '$.AuthenticateCognitoConfig') as authenticate_cognito_config,
  json_extract(action.value, '$.AuthenticateOidcConfig') as authenticate_Oidc_config,
  json_extract(action.value, '$.FixedResponseConfig') as fixed_response_config,
  json_extract(action.value, '$.ForwardConfig.TargetGroupStickinessConfig.DurationSeconds') as duration_seconds,
  json_extract(action.value, '$.ForwardConfig.TargetGroupStickinessConfig.Enabled') as target_group_stickiness_config_enabled
from
  aws_ec2_load_balancer_listener,
  json_each(default_actions) as action;
```

### List of load balancer listeners which listen to HTTP protocol
Discover the segments that are using the HTTP protocol for load balancing. This is useful for identifying potential security risks, as HTTP traffic is unencrypted and can be intercepted.

```sql+postgres
select
  title,
  arn,
  port,
  protocol
from
  aws_ec2_load_balancer_listener
where
  protocol = 'HTTP';
```

```sql+sqlite
select
  title,
  arn,
  port,
  protocol
from
  aws_ec2_load_balancer_listener
where
  protocol = 'HTTP';
```