# Table: aws_autoscaling_group

An Auto Scaling group contains a collection of Amazon EC2 instances that are treated as a logical grouping for the purposes of automatic scaling and management.

## Examples

### Basic auto scaling group info

```sql
select
  name,
  load_balancer_names,
  availability_zones,
  service_linked_role_arn,
  default_cooldown,
  max_size,
  min_size,
  new_instances_protected_from_scale_in
from
  aws_ec2_autoscaling_group;
```


### Autoscaling groups with availability zone count less than 2

```sql
select
  name,
  jsonb_array_length(availability_zones) as az_count
from
  aws_ec2_autoscaling_group
where
  jsonb_array_length(availability_zones) < 2;
```


### Instances' information attached to the autoscaling group

```sql
select
  name as autoscaling_group_name,
  ins_detail ->> 'InstanceId' as instance_id,
  ins_detail ->> 'InstanceType' as instance_type,
  ins_detail ->> 'AvailabilityZone' as az,
  ins_detail ->> 'HealthStatus' as health_status,
  ins_detail ->> 'LaunchConfigurationName' as launch_configuration_name,
  ins_detail -> 'LaunchTemplate' ->> 'LaunchTemplateName' as launch_template_name,
  ins_detail -> 'LaunchTemplate' ->> 'Version' as launch_template_version,
  ins_detail ->> 'ProtectedFromScaleIn' as protected_from_scale_in
from
  aws_ec2_autoscaling_group
  cross join jsonb_array_elements(instances) as ins_detail;
```


### Auto scaling group health check info

```sql
select
  name,
  health_check_type,
  health_check_grace_period
from
  aws_ec2_autoscaling_group;
```