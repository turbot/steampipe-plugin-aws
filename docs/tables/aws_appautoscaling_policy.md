---
title: "Table: aws_appautoscaling_policy - Query AWS Application Auto Scaling Policies using SQL"
description: "Allows users to query AWS Application Auto Scaling Policies to obtain information about their configuration, attached resources, and other metadata."
---

# Table: aws_appautoscaling_policy - Query AWS Application Auto Scaling Policies using SQL

The `aws_appautoscaling_policy` table in Steampipe provides information about Application Auto Scaling policies in AWS. This table allows DevOps engineers, system administrators, and other technical professionals to query policy-specific details, including the scaling target, scaling dimensions, and associated metadata. Users can utilize this table to gather insights on policies, such as policy configurations, attached resources, scaling activities, and more. The schema outlines the various attributes of the Application Auto Scaling policy, including the policy ARN, policy type, creation time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_appautoscaling_policy` table, you can use the `.inspect aws_appautoscaling_policy` command in Steampipe.

### Key columns:

- `policy_name`: The name of the scaling policy. This is a key identifier and can be used to join this table with other tables that contain policy name information.

- `resource_id`: The identifier of the resource associated with the scaling policy. This can be used to join this table with other tables that contain resource ID information.

- `policy_arn`: The Amazon Resource Name (ARN) of the scaling policy. This can be used to join this table with other tables that contain policy ARN information.

## Examples

### Basic info

```sql
select
  service_namespace,
  scalable_dimension,
  policy_type,
  resource_id,
  creation_time
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs';
```

### List policies for ECS services with policy type Step scaling

```sql
select
  resource_id,
  policy_type
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs'
  and policy_type = 'StepScaling';
```

### List policies for ECS services created in the last 30 days

```sql
select
  resource_id,
  policy_type
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs'
  and creation_time > now() - interval '30 days';
```

### Get the CloudWatch alarms associated with the Auto Scaling policy

```sql
select
  resource_id,
  policy_type,
  jsonb_array_elements(alarms) -> 'AlarmName' as alarm_name
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs';
```

### Get the configuration for Step scaling type policies

```sql
select
  resource_id,
  policy_type,
  step_scaling_policy_configuration
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs'
  and policy_type = 'StepScaling';
```
