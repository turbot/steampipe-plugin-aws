---
title: "Steampipe Table: aws_codedeploy_deployment_group - Query AWS CodeDeploy Deployment Groups using SQL"
description: "Allows users to query AWS CodeDeploy Deployment Group details including deployment configurations, target revisions, and associated alarm configurations."
folder: "CodeDeploy"
---

# Table: aws_codedeploy_deployment_group - Query AWS CodeDeploy Deployment Groups using SQL

The AWS CodeDeploy Deployment Group is a set of individual instances, CodeDeploy Lambda deployment configuration settings, or an EC2 tag set. It is used to represent a deployment's target, be it an instance, a Lambda function, or an EC2 instance. The Deployment Group is a key component of the AWS CodeDeploy service, which automates code deployments to any instance, including Amazon EC2 instances and servers running on-premise.

## Table Usage Guide

The `aws_codedeploy_deployment_group` table in Steampipe provides you with information about deployment groups within AWS CodeDeploy. This table allows you as a DevOps engineer to query deployment group-specific details, including deployment configurations, target revisions, and associated alarm configurations. You can utilize this table to gather insights on deployment groups, such as deployment configuration names, target revisions, and alarm configurations. The schema outlines the various attributes of the deployment group for you, including the deployment group name, service role ARN, deployment configuration name, target revision, and associated alarm configurations.

## Examples

### Basic info
Explore which deployment groups are active in your AWS CodeDeploy application, including their deployment style and region. This can help identify any inconsistencies or areas for optimization in deployment strategies.

```sql+postgres
select
  arn,
  deployment_group_id,
  deployment_group_name,
  application_name,
  deployment_style,
  region
from
  aws_codedeploy_deployment_group;
```

```sql+sqlite
select
  arn,
  deployment_group_id,
  deployment_group_name,
  application_name,
  deployment_style,
  region
from
  aws_codedeploy_deployment_group;
```

### Get total deployment groups on each platform
Determine the total number of deployment groups across each computing platform. This can provide insights into the distribution of resources and help in effective resource management.

```sql+postgres
select
  count(arn) as group_count,
  compute_platform
from
  aws_codedeploy_deployment_group
group by
  compute_platform;
```

```sql+sqlite
select
  count(arn) as group_count,
  compute_platform
from
  aws_codedeploy_deployment_group
group by
  compute_platform;
```

### List the last successful deployment for each deployment group
Determine the status of your most recent successful deployments across different deployment groups. This can help you track your deployment history and identify any potential issues or bottlenecks in your deployment process.

```sql+postgres
select
  arn,
  deployment_group_id,
  last_successful_deployment
from
  aws_codedeploy_deployment_group;
```

```sql+sqlite
select
  arn,
  deployment_group_id,
  last_successful_deployment
from
  aws_codedeploy_deployment_group;
```

### Get total deployment groups based on deployment style
Analyze your deployment styles to understand the distribution of your deployment groups. This can help optimize resource allocation and improve deployment efficiency.

```sql+postgres
select
  count(arn) as group_count,
  deployment_style
from
  aws_codedeploy_deployment_group
group by
  deployment_style;
```

```sql+sqlite
select
  count(arn) as group_count,
  deployment_style
from
  aws_codedeploy_deployment_group
group by
  deployment_style;
```

### List the deployment groups having automatic rollback enabled
Determine the areas in which automatic rollback is enabled for deployment groups. This is useful to quickly identify configurations that can help prevent unintended changes or disruptions to services.

```sql+postgres
select
  arn,
  deployment_group_id,
  deployment_group_name,
  auto_rollback_configuration ->> 'Enabled' as auto_rollback_configuration_enabled
from
  aws_codedeploy_deployment_group
where
  auto_rollback_configuration ->> 'Enabled' = 'true';
```

```sql+sqlite
select
  arn,
  deployment_group_id,
  deployment_group_name,
  json_extract(auto_rollback_configuration, '$.Enabled') as auto_rollback_configuration_enabled
from
  aws_codedeploy_deployment_group
where
  json_extract(auto_rollback_configuration, '$.Enabled') = 'true';
```

### List all autoscaling groups in a particular deployment group for an application
Analyze the settings to understand the configuration of autoscaling groups within a specific deployment group for a particular application. This can be useful in managing and optimizing resource usage in your cloud environment.

```sql+postgres
select
  arn as group_arn,
  deployment_group_id,
  deployment_group_name,
  auto_scaling_groups ->> 'Hook' as auto_scaling_group_hook,
  auto_scaling_groups ->> 'Name' as auto_scaling_group_name
from
  aws_codedeploy_deployment_group
where
  application_name = 'abc'
  and deployment_group_name = 'def';
```

```sql+sqlite
select
  arn as group_arn,
  deployment_group_id,
  deployment_group_name,
  json_extract(auto_scaling_groups, '$.Hook') as auto_scaling_group_hook,
  json_extract(auto_scaling_groups, '$.Name') as auto_scaling_group_name
from
  aws_codedeploy_deployment_group
where
  application_name = 'abc'
  and deployment_group_name = 'def';
```

### List the deployment groups having automatic rollback enabled
Determine the areas in which automatic rollback is enabled in deployment groups. This is useful to identify and manage risk in software deployment processes.

```sql+postgres
select
  arn,
  deployment_group_id,
  deployment_group_name,
  alarm_configuration ->> 'Enabled' as alarm_configuration_enabled
from
  aws_codedeploy_deployment_group
where
  alarm_configuration ->> 'Enabled' = 'true';
```

```sql+sqlite
select
  arn,
  deployment_group_id,
  deployment_group_name,
  json_extract(alarm_configuration, '$.Enabled') as alarm_configuration_enabled
from
  aws_codedeploy_deployment_group
where
  json_extract(alarm_configuration, '$.Enabled') = 'true';
```