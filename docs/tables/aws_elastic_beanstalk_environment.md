---
title: "Table: aws_elastic_beanstalk_environment - Query AWS Elastic Beanstalk Environments using SQL"
description: "Allows users to query AWS Elastic Beanstalk Environments to gain insights into their configuration, status, health, related applications, and other metadata."
---

# Table: aws_elastic_beanstalk_environment - Query AWS Elastic Beanstalk Environments using SQL

The `aws_elastic_beanstalk_environment` table in Steampipe provides information about environments within AWS Elastic Beanstalk. This table allows DevOps engineers to query environment-specific details, including configuration settings, environment health, related applications, and associated metadata. Users can utilize this table to gather insights on environments, such as environments with specific configurations, health status, associated applications, and more. The schema outlines the various attributes of the Elastic Beanstalk environment, including the environment name, ID, application name, status, health, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elastic_beanstalk_environment` table, you can use the `.inspect aws_elastic_beanstalk_environment` command in Steampipe.

**Key columns**:

- `environment_name`: The name of the environment. This can be used to join this table with other tables that contain environment-specific information.
- `application_name`: The name of the application associated with the environment. This can be used to join this table with application-specific tables in AWS Elastic Beanstalk.
- `environment_id`: The ID of the environment. This can be used to join this table with other tables that contain environment-specific information, especially when multiple environments have the same name.

## Examples

### Basic info

```sql
select
  environment_id,
  environment_name,
  application_name,
  arn,
  tier
from
  aws_elastic_beanstalk_environment;
```

### List environments which have configuration updates and application version deployments in progress

```sql
select
  environment_name,
  abortable_operation_in_progress
from
  aws_elastic_beanstalk_environment
where
  abortable_operation_in_progress = 'true';
```

### List unhealthy environments

```sql
select
  environment_name,
  application_name,
  environment_id,
  health
from
  aws_elastic_beanstalk_environment
where
  health = 'Red';
```

### List environments with health monitoring disabled

```sql
select
  environment_name,
  health_status
from
  aws_elastic_beanstalk_environment
where
  health_status = 'Suspended';
```

### List managed actions for each environment

```sql
select
  environment_name,
  a ->> 'ActionDescription' as action_description,
  a ->> 'ActionId' as action_id,
  a ->> 'ActionType' as action_type,
  a ->> 'Status' as action_status,
  a ->> 'WindowStartTime' as action_window_start_time
from
  aws_elastic_beanstalk_environment,
  jsonb_array_elements(managed_actions) as a;
```

### list the configuration settings for each environment

```sql
select
  environment_name,
  application_name,
  c ->> 'DateCreated' as date_created,
  c ->> 'DateUpdated' as date_updated,
  c ->> 'DeploymentStatus' as deployment_status,
  c ->> 'Description' as description,
  c -> 'OptionSettings' ->> 'Namespace' as option_settings_namespace,
  c -> 'OptionSettings' ->> 'OptionName' as option_name,
  c -> 'OptionSettings' ->> 'ResourceName' as option_resource_name,
  c -> 'OptionSettings' ->> 'Value' as option_value,
  c ->> 'PlatformArn' as platform_arn,
  c ->> 'SolutionStackName' as solution_stack_name,
  c ->> 'TemplateName' as template_name
from
  aws_elastic_beanstalk_environment,
  jsonb_array_elements(configuration_settings) as c;
```