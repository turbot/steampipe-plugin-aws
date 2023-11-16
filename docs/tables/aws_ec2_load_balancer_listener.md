---
title: "Table: aws_ec2_load_balancer_listener - Query AWS EC2 Load Balancer Listeners using SQL"
description: "Allows users to query AWS EC2 Load Balancer Listener data, which provides information about listeners for an Application Load Balancer or Network Load Balancer."
---

# Table: aws_ec2_load_balancer_listener - Query AWS EC2 Load Balancer Listeners using SQL

The `aws_ec2_load_balancer_listener` table in Steampipe provides information about listeners for an Application Load Balancer or Network Load Balancer in Amazon Elastic Compute Cloud (EC2). This table allows DevOps engineers to query listener-specific details, including protocol, port, SSL policy, and associated actions. Users can utilize this table to gather insights on listeners, such as their current state, default actions, and certificates. The schema outlines the various attributes of the Load Balancer Listener, including the listener ARN, load balancer ARN, default actions, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_load_balancer_listener` table, you can use the `.inspect aws_ec2_load_balancer_listener` command in Steampipe.

### Key columns:

- `listener_arn`: The Amazon Resource Name (ARN) of the listener. This can be used to join this table with other tables.
- `load_balancer_arn`: The ARN of the load balancer. This can be used to join this table with the `aws_ec2_load_balancer` table.
- `port`: The port on which the load balancer is listening. This column is useful in understanding the traffic flow of the load balancer.

## Examples

### Load balancer listener basic info

```sql
select
  title,
  arn,
  port,
  protocol
from
  aws_ec2_load_balancer_listener;
```


### Action configuration details of each load balancer

```sql
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


### List of load balancer listeners which listen to HTTP protocol

```sql
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
