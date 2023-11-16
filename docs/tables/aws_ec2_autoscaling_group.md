---
title: "Table: aws_ec2_autoscaling_group - Query AWS EC2 Auto Scaling Groups using SQL"
description: "Allows users to query AWS EC2 Auto Scaling Groups and access detailed information about each group's configuration, instances, policies, and more."
---

# Table: aws_ec2_autoscaling_group - Query AWS EC2 Auto Scaling Groups using SQL

The `aws_ec2_autoscaling_group` table in Steampipe provides information about Auto Scaling Groups within AWS EC2. This table allows DevOps engineers to query group-specific details, including configuration, associated instances, scaling policies, and associated metadata. Users can utilize this table to gather insights on groups, such as their desired, minimum and maximum sizes, default cooldown periods, load balancer names, and more. The schema outlines the various attributes of the Auto Scaling Group, including the ARN, creation date, health check type and grace period, launch configuration name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_autoscaling_group` table, you can use the `.inspect aws_ec2_autoscaling_group` command in Steampipe.

**Key columns**:

- `auto_scaling_group_arn`: The Amazon Resource Name (ARN) of the Auto Scaling group. This can be used to join with other tables that store ARN information.
- `auto_scaling_group_name`: The name of the Auto Scaling group. This is a unique identifier and can be used to join with other tables that store Auto Scaling group names.
- `vpc_zone_identifier`: The identifiers of the subnets to which the Auto Scaling group is assigned. This can be used to join with other tables that store subnet identifiers.

## Examples

### Basic info

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
  aws_ec2_autoscaling_group,
  jsonb_array_elements(instances) as ins_detail;
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
