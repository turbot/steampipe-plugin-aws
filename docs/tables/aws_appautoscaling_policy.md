---
title: "Steampipe Table: aws_appautoscaling_policy - Query AWS Application Auto Scaling Policies using SQL"
description: "Allows users to query AWS Application Auto Scaling Policies to obtain information about their configuration, attached resources, and other metadata."
folder: "Application Auto Scaling"
---

# Table: aws_appautoscaling_policy - Query AWS Application Auto Scaling Policies using SQL

The AWS Application Auto Scaling Policies allow you to manage the scaling of your applications in response to their demand patterns. They enable automatic adjustments to the scalable target capacity as needed to maintain optimal resource utilization. This ensures that your applications always have the right resources at the right time, improving their performance and reducing costs.

## Table Usage Guide

The `aws_appautoscaling_policy` table in Steampipe provides you with information about Application Auto Scaling policies in AWS. This table allows you, as a DevOps engineer, system administrator, or other technical professional, to query policy-specific details, including the scaling target, scaling dimensions, and associated metadata. You can utilize this table to gather insights on policies, such as policy configurations, attached resources, scaling activities, and more. The schema outlines the various attributes of the Application Auto Scaling policy, including the policy ARN, policy type, creation time, and associated tags for you.

**Important Notes**
- You **_must_** specify `service_namespace` in a `where` clause in order to use this table.
- For supported values of the service namespace, please refer [Service Namespace](https://docs.aws.amazon.com/autoscaling/application/APIReference/API_ScalingPolicy.html#autoscaling-Type-ScalingPolicy-ServiceNamespace).
- This table supports optional quals. Queries with optional quals are optimised to use additional filtering provided by the AWS API function. Optional quals are supported for the following columns:
  - `policy_name`
  - `resource_id`

## Examples

### Basic info
Analyze the settings to understand the policies associated with your AWS ECS service. This can help in managing the scaling behavior of the resources in the ECS service, identifying the dimensions that are scalable, and the time of policy creation.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which step scaling policies are applied for ECS services. This can help in managing and optimizing resource allocation for your applications.

```sql+postgres
select
  resource_id,
  policy_type
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs'
  and policy_type = 'StepScaling';
```

```sql+sqlite
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
Identify recent policy changes for your ECS services. This query is useful for monitoring and managing your autoscaling configuration, allowing you to track changes made within the last month.

```sql+postgres
select
  resource_id,
  policy_type
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs'
  and creation_time > now() - interval '30 days';
```

```sql+sqlite
select
  resource_id,
  policy_type
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs'
  and creation_time > datetime('now','-30 days');
```

### Get the CloudWatch alarms associated with the Auto Scaling policy
Determine the areas in which CloudWatch alarms are linked to an Auto Scaling policy. This can be beneficial in understanding the alarm triggers and managing resources within the Elastic Container Service (ECS).

```sql+postgres
select
  resource_id,
  policy_type,
  jsonb_array_elements(alarms) -> 'AlarmName' as alarm_name
from
  aws_appautoscaling_policy
where
  service_namespace = 'ecs';
```

```sql+sqlite
select
  resource_id,
  policy_type,
  json_extract(json_each.value, '$.AlarmName') as alarm_name
from
  aws_appautoscaling_policy,
  json_each(alarms)
where
  service_namespace = 'ecs';
```

### Get the configuration for Step scaling type policies
Explore the setup of step scaling policies within the ECS service namespace to understand how application auto scaling is configured.

```sql+postgres
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

```sql+sqlite
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
