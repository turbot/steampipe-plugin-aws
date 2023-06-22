# Table: aws_elastic_beanstalk_environment

AWS Elastic Beanstalk makes it easy to create new environments for your application. You can create and manage separate environments for development, testing, and production use, and you can deploy any version of your application to any environment. Environments can be long-running or temporary. When you terminate an environment, you can save its configuration to recreate it later.

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