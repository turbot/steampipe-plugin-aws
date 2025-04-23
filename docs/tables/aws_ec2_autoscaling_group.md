---
title: "Steampipe Table: aws_ec2_autoscaling_group - Query AWS EC2 Auto Scaling Groups using SQL"
description: "Allows users to query AWS EC2 Auto Scaling Groups and access detailed information about each group's configuration, instances, policies, and more."
folder: "Auto Scaling"
---

# Table: aws_ec2_autoscaling_group - Query AWS EC2 Auto Scaling Groups using SQL

The AWS EC2 Auto Scaling Groups service allows you to ensure that you have the correct number of Amazon EC2 instances available to handle the load for your applications. Auto Scaling Groups contain a collection of EC2 instances that share similar characteristics and are treated as a logical grouping for the purposes of instance scaling and management. This service automatically increases or decreases the number of instances depending on the demand, ensuring optimal performance and cost management.

## Table Usage Guide

The `aws_ec2_autoscaling_group` table in Steampipe provides you with information about Auto Scaling Groups within AWS EC2. This table allows you, as a DevOps engineer, to query group-specific details, including configuration, associated instances, scaling policies, and associated metadata. You can utilize this table to gather insights on groups, such as their desired, minimum and maximum sizes, default cooldown periods, load balancer names, and more. The schema outlines for you the various attributes of the Auto Scaling Group, including the ARN, creation date, health check type and grace period, launch configuration name, and associated tags.

## Examples

### Basic info
Explore the configuration of your AWS EC2 autoscaling group to understand its operational parameters, such as the default cooldown period and size limitations. This can help you optimize resource allocation and improve cost efficiency in your cloud environment.

```sql+postgres
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

```sql+sqlite
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
Identify autoscaling groups that may not be optimally configured for high availability due to having less than two availability zones. This can be useful to improve fault tolerance and ensure uninterrupted service.

```sql+postgres
select
  name,
  jsonb_array_length(availability_zones) as az_count
from
  aws_ec2_autoscaling_group
where
  jsonb_array_length(availability_zones) < 2;
```

```sql+sqlite
select
  name,
  json_array_length(availability_zones) as az_count
from
  aws_ec2_autoscaling_group
where
  json_array_length(availability_zones) < 2;
```


### Instances' information attached to the autoscaling group
Explore the health and configuration status of instances within an autoscaling group. This is useful to monitor and manage the scalability and availability of your AWS EC2 resources.

```sql+postgres
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

```sql+sqlite
select
  name as autoscaling_group_name,
  json_extract(ins_detail, '$.InstanceId') as instance_id,
  json_extract(ins_detail, '$.InstanceType') as instance_type,
  json_extract(ins_detail, '$.AvailabilityZone') as az,
  json_extract(ins_detail, '$.HealthStatus') as health_status,
  json_extract(ins_detail, '$.LaunchConfigurationName') as launch_configuration_name,
  json_extract(ins_detail, '$.LaunchTemplate.LaunchTemplateName') as launch_template_name,
  json_extract(ins_detail, '$.LaunchTemplate.Version') as launch_template_version,
  json_extract(ins_detail, '$.ProtectedFromScaleIn') as protected_from_scale_in
from
  aws_ec2_autoscaling_group,
  json_each(instances) as ins_detail;
```

### Auto scaling group health check info
Explore the health check settings of your auto scaling groups to understand their operational readiness and grace periods. This can help you assess the resilience of your system and plan for contingencies.

```sql+postgres
select
  name,
  health_check_type,
  health_check_grace_period
from
  aws_ec2_autoscaling_group;
```

```sql+sqlite
select
  name,
  health_check_type,
  health_check_grace_period
from
  aws_ec2_autoscaling_group;
```
