# Table: aws_ec2_load_balancer_listener

A listener is a process that checks for connection requests. It is configured with a protocol and a port for front-end (client to load balancer) connections, and a protocol and a port for back-end (load balancer to back-end instance) connections.

Note: This table lists the listeners for application load balancers and network load balancers only.

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
