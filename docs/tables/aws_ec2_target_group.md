# Table: aws_ec2_target_group

Target group is used to route requests to one or more registered targets for a load balancer.

## Examples

### Basic target group info

```sql
select
  target_group_name,
  target_type,
  load_balancer_arns,
  vpc_id
from
  aws_ec2_target_group;
```


### Health check info of target groups

```sql
select
  health_check_enabled,
  protocol,
  matcher_http_code,
  healthy_threshold_count,
  unhealthy_threshold_count,
  health_check_enabled,
  health_check_interval_seconds,
  health_check_path,
  health_check_port,
  health_check_protocol,
  health_check_timeout_seconds
from
  aws_ec2_target_group;
```


### Registered target for each target group

```sql
select
  target_group_name,
  target_type,
  target -> 'Target' ->> 'AvailabilityZone' as availability_zone,
  target -> 'Target' ->> 'Id' as id,
  target -> 'Target' ->> 'Port' as port
from
  aws_ec2_target_group
  cross join jsonb_array_elements(target_health_descriptions) as target;
```


### Health status of registered targets

```sql
select
  target_group_name,
  target_type,
  target -> 'TargetHealth' ->> 'Description' as description,
  target -> 'TargetHealth' ->> 'Reason' reason,
  target -> 'TargetHealth' ->> 'State' as state
from
  aws_ec2_target_group
  cross join jsonb_array_elements(target_health_descriptions) as target;
```
