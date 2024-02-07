---
title: "Steampipe Table: aws_iot_greengrass_deployment - Query AWS IoT Greengrass Deployments using SQL"
description: "Allows users to query AWS IoT Greengrass Deployments. This table provides information about Greengrass Deployments within AWS IoT Greengrass, enabling users to gather insights on deployments such as deployment name, ID, creation timestamp, and more."
---

# Table: aws_iot_greengrass_deployment - Query AWS IoT Greengrass Deployments using SQL

The AWS IoT Greengrass Deployment is a feature within the AWS IoT Greengrass service that manages the deployment of software to IoT devices. Each core device runs the components of the deployments for that device. A new deployment to the same target overwrites the previous deployment to the target.

When revising a deployment for a target, it replaces the components from the previous revision with the components in the new revision. For instance, if you initially deploy certain components to a group and later create another deployment for the same group with different components, the core devices in that group will run the new set of components.

## Table Usage Guide

The `aws_iot_greengrass_deployment` table in Steampipe enables you, as an IoT engineer or DevOps professional, to query deployment-specific details within AWS IoT Greengrass. This includes information like deployment name, ID, whether it's the latest for the target, creation timestamp, and associated target ARNs. You can use this table to monitor and manage deployments, ensuring that your IoT devices are running the desired software configurations and versions.

## Examples

### Basic info
Query basic information about Greengrass deployments, such as their names, IDs, and associated target ARNs. This is useful for a quick overview of the deployments in your environment.

```sql+postgres
select
  deployment_name,
  deployment_id,
  is_latest_for_target,
  creation_timestamp,
  parent_target_arn,
  target_arn
from
  aws_iot_greengrass_deployment;
```

```sql+sqlite
select
  deployment_name,
  deployment_id,
  is_latest_for_target,
  creation_timestamp,
  parent_target_arn,
  target_arn
from
  aws_iot_greengrass_deployment;
``

### List deployments created in the last 30 days
Identify recent deployments to manage and monitor new changes in your IoT environment.

```sql+postgres
select
  deployment_name,
  deployment_id,
  is_latest_for_target,
  creation_timestamp,
  target_arn
from
  aws_iot_greengrass_deployment
where
  creation_timestamp >= now() - interval '30' day;
```

```sql+sqlite
select
  deployment_name,
  deployment_id,
  is_latest_for_target,
  creation_timestamp,
  target_arn
from
  aws_iot_greengrass_deployment
where
  creation_timestamp >= datetime('now', '-30 day');
```

### List deployments that are the latest revision for its target
Focus on the most current deployments for each target to ensure up-to-date configuration and software.

```sql+postgres
select
  deployment_name,
  deployment_id,
  is_latest_for_target,
  revision_id,
  iot_job_arn,
  iot_job_id
from
  aws_iot_greengrass_deployment
where
  is_latest_for_target;
```

```sql+sqlite
select
  deployment_name,
  deployment_id,
  is_latest_for_target,
  revision_id,
  iot_job_arn,
  iot_job_id
from
  aws_iot_greengrass_deployment
where
  is_latest_for_target;
```

### Get deployment policy details of the deployments
Examine the specific policies associated with each deployment to understand how they are configured and managed.

```sql+postgres
select
  deployment_name,
  deployment_policies -> 'ComponentUpdatePolicy' as component_update_policy,
  deployment_policies -> 'ConfigurationValidationPolicy' as configuration_validation_policy,
  deployment_policies -> 'FailureHandlingPolicy' as failure_handling_policy
from
  aws_iot_greengrass_deployment
where
  is_latest_for_target;
```

```sql+sqlite
select
  deployment_name,
  json_extract(deployment_policies, '$.ComponentUpdatePolicy') as component_update_policy,
  json_extract(deployment_policies, '$.ConfigurationValidationPolicy') as configuration_validation_policy,
  json_extract(deployment_policies, '$.FailureHandlingPolicy') as failure_handling_policy
from
  aws_iot_greengrass_deployment
where
  is_latest_for_target;
```

### Get IoT job configuration details of deployments
Review the configurations related to IoT jobs associated with the deployments for detailed insight into their execution and management.

```sql+postgres
select
  deployment_name,
  iot_job_configuration -> 'AbortConfig' as abort_config,
  iot_job_configuration -> 'JobExecutionsRolloutConfig' as job_executions_rollout_config,
  iot_job_configuration -> 'TimeoutConfig' as timeout_config
from
  aws_iot_greengrass_deployment
where
  is_latest_for_target
```

```sql+sqlite
select
  deployment_name,
  json_extract(iot_job_configuration, '$.AbortConfig') as abort_config,
  json_extract(iot_job_configuration, '$.JobExecutionsRolloutConfig') as job_executions_rollout_config,
  json_extract(iot_job_configuration, '$.TimeoutConfig') as timeout_config
from
  aws_iot_greengrass_deployment
where
  is_latest_for_target;
```

### Get IoT thing group details of deployments
Link deployments to specific IoT thing groups for a comprehensive understanding of which groups are affected by which deployments.

```sql+postgres
select
  d.deployment_name,
  d.target_arn as thing_group_arn,
  d.parent_target_arn,
  g.group_name,
  g.status as thing_group_status,
  g.query_string,
  g.parent_group_name
from
  aws_iot_greengrass_deployment as d,
  aws_iot_thing_group as g
where
  g.arn = d.target_arn;
```

```sql+sqlite
select
  d.deployment_name,
  d.target_arn as thing_group_arn,
  d.parent_target_arn,
  g.group_name,
  g.status as thing_group_status,
  g.query_string,
  g.parent_group_name
from
  aws_iot_greengrass_deployment as d
join
  aws_iot_thing_group as g on g.arn = d.target_arn;
```