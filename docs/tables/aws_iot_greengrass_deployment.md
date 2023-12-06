# Table: aws_iot_greengrass_deployment

Each core device runs the components of the deployments for that device. A new deployment to the same target overwrites the previous deployment to the target. When you create a deployment, you define the components and configurations to apply to the core device's existing software.

When you revise a deployment for a target, you replace the components from the previous revision with the components in the new revision. For example, you deploy the Log manager and Secret manager components to the thing group TestGroup. Then you create another deployment for TestGroup that specifies only the secret manager component. As a result, the core devices in that group no longer run the log manager.

## Examples

### Basic info

```sql
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

### List deployments created in the last 30 days

```sql
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

### List deployments that are the latest revision for its target

```sql
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

```sql
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

### Get IoT job configuration details of deployments

```sql
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

### Get IoT thing group details of deployments

```sql
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