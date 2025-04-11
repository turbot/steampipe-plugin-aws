---
title: "Steampipe Table: aws_elastic_beanstalk_environment - Query AWS Elastic Beanstalk Environments using SQL"
description: "Allows users to query AWS Elastic Beanstalk Environments to gain insights into their configuration, status, health, related applications, and other metadata."
folder: "Elastic Beanstalk"
---

# Table: aws_elastic_beanstalk_environment - Query AWS Elastic Beanstalk Environments using SQL

The AWS Elastic Beanstalk Environment is a part of the AWS Elastic Beanstalk service that allows developers to deploy and manage applications in the AWS cloud without worrying about the infrastructure that runs those applications. This service automatically handles the capacity provisioning, load balancing, scaling, and application health monitoring. It supports applications developed in Java, .NET, PHP, Node.js, Python, Ruby, Go, and Docker.

## Table Usage Guide

The `aws_elastic_beanstalk_environment` table in Steampipe provides you with information about environments within AWS Elastic Beanstalk. This table allows you as a DevOps engineer to query environment-specific details, including configuration settings, environment health, related applications, and associated metadata. You can utilize this table to gather insights on environments, such as environments with specific configurations, health status, associated applications, and more. The schema outlines the various attributes of the Elastic Beanstalk environment for you, including the environment name, ID, application name, status, health, and associated tags.

## Examples

### Basic info
Explore the configuration of your AWS Elastic Beanstalk environments to understand their applications and tiers. This is useful for reviewing the setup and organization of your cloud applications.

```sql+postgres
select
  environment_id,
  environment_name,
  application_name,
  arn,
  tier
from
  aws_elastic_beanstalk_environment;
```

```sql+sqlite
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
Identify instances where configuration updates and application version deployments are currently in progress. This can be useful in managing and tracking ongoing operations within your environment.

```sql+postgres
select
  environment_name,
  abortable_operation_in_progress
from
  aws_elastic_beanstalk_environment
where
  abortable_operation_in_progress = 'true';
```

```sql+sqlite
select
  environment_name,
  abortable_operation_in_progress
from
  aws_elastic_beanstalk_environment
where
  abortable_operation_in_progress = 'true';
```

### List unhealthy environments
Determine the areas in which AWS Elastic Beanstalk environments are unhealthy. This query is useful for identifying and addressing problematic environments to ensure optimal application performance.

```sql+postgres
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

```sql+sqlite
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
Identify instances where health monitoring has been suspended in certain environments to understand potential vulnerabilities and ensure optimal performance.

```sql+postgres
select
  environment_name,
  health_status
from
  aws_elastic_beanstalk_environment
where
  health_status = 'Suspended';
```

```sql+sqlite
select
  environment_name,
  health_status
from
  aws_elastic_beanstalk_environment
where
  health_status = 'Suspended';
```

### List managed actions for each environment
Identify the managed actions associated with each environment in the AWS Elastic Beanstalk service. This can help in monitoring the status and type of actions, providing insights for better management and optimization of your environments.

```sql+postgres
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

```sql+sqlite
select
  environment_name,
  json_extract(a.value, '$.ActionDescription') as action_description,
  json_extract(a.value, '$.ActionId') as action_id,
  json_extract(a.value, '$.ActionType') as action_type,
  json_extract(a.value, '$.Status') as action_status,
  json_extract(a.value, '$.WindowStartTime') as action_window_start_time
from
  aws_elastic_beanstalk_environment,
  json_each(managed_actions) as a;
```

### list the configuration settings for each environment
Determine the areas in which configuration settings for various environments are tracked and updated. This can be used to keep track of deployment status, platform details, and other critical factors in your AWS Elastic Beanstalk environments.

```sql+postgres
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

```sql+sqlite
select
  environment_name,
  application_name,
  json_extract(c.value, '$.DateCreated') as date_created,
  json_extract(c.value, '$.DateUpdated') as date_updated,
  json_extract(c.value, '$.DeploymentStatus') as deployment_status,
  json_extract(c.value, '$.Description') as description,
  json_extract(c.value, '$.OptionSettings.Namespace') as option_settings_namespace,
  json_extract(c.value, '$.OptionSettings.OptionName') as option_name,
  json_extract(c.value, '$.OptionSettings.ResourceName') as option_resource_name,
  json_extract(c.value, '$.OptionSettings.Value') as option_value,
  json_extract(c.value, '$.PlatformArn') as platform_arn,
  json_extract(c.value, '$.SolutionStackName') as solution_stack_name,
  json_extract(c.value, '$.TemplateName') as template_name
from
  aws_elastic_beanstalk_environment,
  json_each(configuration_settings) as c;
```