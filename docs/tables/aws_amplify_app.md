---
title: "Steampipe Table: aws_amplify_app - Query AWS Amplify Apps using SQL"
description: "Allows users to query AWS Amplify Apps to retrieve detailed information about each application, including its name, ARN, creation date, default domain, and more."
folder: "Amplify"
---

# Table: aws_amplify_app - Query AWS Amplify Apps using SQL

The AWS Amplify App is a part of AWS Amplify, a set of tools and services that enables developers to build secure, scalable, full-stack applications. These applications can be built with integrated backend services like authentication, analytics, and content delivery, with capabilities such as real-time data syncing. AWS Amplify Apps allow for the creation, configuration, and management of continuous deployment workflows for web apps in the AWS Amplify Console.

## Table Usage Guide

The `aws_amplify_app` table in Steampipe provides you with information about apps within AWS Amplify. This table allows you, as a DevOps engineer, to query app-specific details, including the name, ARN, creation date, last update date, default domain, and associated metadata. You can utilize this table to gather insights on Amplify Apps, such as the apps' status, platform, repository, and more. The schema outlines the various attributes of the Amplify App for you, including the app ID, app ARN, platform, repository, production branch, and associated tags.

## Examples

### Basic info
Explore the fundamental details of your AWS Amplify applications to gain insights into their creation time, platform, and build specifications. This aids in understanding the overall structure and configuration of your applications for better management and optimization.The query provides an overview of applications in AWS Amplify, including their identification details, creation time, and platform information. This can be useful for auditing purposes or to gain a quick understanding of the applications' configuration and status.


```sql+postgres
select
  app_id,
  name,
  description,
  arn,
  platform,
  create_time,
  build_spec
from
  aws_amplify_app;
```

```sql+sqlite
select
  app_id,
  name,
  description,
  arn,
  platform,
  create_time,
  build_spec
from
  aws_amplify_app;
```

### List apps created within the last 90 days
Discover the segments that have recently been developed within the past 3 months. This could be beneficial for understanding the evolution of your app portfolio and identifying any new trends or patterns.The query identifies recently created applications on AWS Amplify by filtering those that have been established within the last 90 days. This could be useful for monitoring new app development or tracking changes over a quarterly period.


```sql+postgres
select
  name,
  app_id,
  create_time
from
  aws_amplify_app
where
  create_time >= (now() - interval '90' day)
order by
  create_time;
```

```sql+sqlite
select
  name,
  app_id,
  create_time
from
  aws_amplify_app
where
  create_time >= datetime('now', '-90 day')
order by
  create_time;
```

### List apps updated within the last hour
Explore which applications have been updated in the last hour to stay informed about the most recent changes. This is useful for closely monitoring application updates and ensuring they are functioning as expected after each update.The query identifies and organizes applications that have been updated within the past hour. This could be useful for monitoring recent changes to applications, particularly in a large or rapidly evolving system.


```sql+postgres
select
  name,
  app_id,
  update_time
from
  aws_amplify_app
where
  update_time >= (now() - interval '1' hour)
order by
  update_time;
```

```sql+sqlite
select
  name,
  app_id,
  update_time
from
  aws_amplify_app
where
  update_time >= datetime('now', '-1 hour')
order by
  update_time;
```

### Describe information about the production branch for an app
This query is used to gain insights into the status of a specific application's production branch. It provides crucial information such as the last deployment time and branch status, which can be useful for monitoring and managing app development and deployment.The query provides a snapshot of the production branch status for a specific application, including when it was last deployed. This is useful for tracking the progress and status of application updates.


```sql+postgres
select
  production_branch ->> 'BranchName' as branch_name,
  production_branch ->> 'LastDeployTime' as last_deploy_time,
  production_branch ->> 'Status' as status
from
  aws_amplify_app
where
  name = 'amplify_app_name';
```

```sql+sqlite
select
  json_extract(production_branch, '$.BranchName') as branch_name,
  json_extract(production_branch, '$.LastDeployTime') as last_deploy_time,
  json_extract(production_branch, '$.Status') as status
from
  aws_amplify_app
where
  name = 'amplify_app_name';
```

### List information about the build spec for an app
Explore the build specifications for a specific application in AWS Amplify, including the backend, frontend, test, and environment settings. This can help fine-tune development and testing processes, and ensure optimal environment configurations.The query provides specific details about the construction specifications for both the backend and frontend of an application, including testing protocols and environmental settings. This information can be useful for developers looking to understand the structure and configuration of an app for troubleshooting or replication purposes.


```sql+postgres
select
  name,
  app_id,
  build_spec ->> 'backend' as build_backend_spec,
  build_spec ->> 'frontend' as build_frontend_spec,
  build_spec ->> 'test' as build_test_spec,
  build_spec ->> 'env' as build_env_settings
from
  aws_amplify_app
where
  name = 'amplify_app_name';
```

```sql+sqlite
select
  name,
  app_id,
  json_extract(build_spec, '$.backend') as build_backend_spec,
  json_extract(build_spec, '$.frontend') as build_frontend_spec,
  json_extract(build_spec, '$.test') as build_test_spec,
  json_extract(build_spec, '$.env') as build_env_settings
from
  aws_amplify_app
where
  name = 'amplify_app_name';
```

### List information on rewrite(200) redirect settings for an app
This example allows you to identify instances where an app is using rewrite (200) redirect settings. It's particularly useful for understanding the conditions, sources, and targets associated with these redirects, helping to optimize app navigation and user experience.The query provides an overview of the 200 status code redirect settings for a specific application, allowing users to assess and manage how traffic is being redirected within the app. This can be useful for troubleshooting or optimizing application performance.


```sql+postgres
select
  name,
  redirects_array ->> 'Condition' as country_code,
  redirects_array ->> 'Source' as source_address,
  redirects_array ->> 'Status' as redirect_type,
  redirects_array ->> 'Target' as destination_address
from
  aws_amplify_app,
  jsonb_array_elements(custom_rules) as redirects_array
where
  redirects_array ->> 'Status' = '200'
  and name = 'amplify_app_name';
```

```sql+sqlite
select
  name,
  json_extract(redirects_array.value, '$.Condition') as country_code,
  json_extract(redirects_array.value, '$.Source') as source_address,
  json_extract(redirects_array.value, '$.Status') as redirect_type,
  json_extract(redirects_array.value, '$.Target') as destination_address
from
  aws_amplify_app,
  json_each(custom_rules) as redirects_array
where
  json_extract(redirects_array.value, '$.Status') = '200'
  and name = 'amplify_app_name';
```

### List all apps that have branch auto build enabled
Determine the areas in which automatic build feature is enabled for applications. This is useful for identifying applications that are set to automatically build and deploy code changes from connected branches, thereby facilitating continuous integration and delivery.The query provides a list of all applications that have automatic branch building enabled. This is useful for developers who want to identify which apps are set to automatically build and deploy when code is pushed to a connected repository, allowing for efficient monitoring and management.


```sql+postgres
select
  app_id,
  name,
  description,
  arn
from
  aws_amplify_app
where
  enable_branch_auto_build = true;
```

```sql+sqlite
select
  app_id,
  name,
  description,
  arn
from
  aws_amplify_app
where
  enable_branch_auto_build = 1;
```