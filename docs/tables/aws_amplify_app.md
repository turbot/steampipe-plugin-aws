---
title: "Table: aws_amplify_app - Query AWS Amplify Apps using SQL"
description: "Allows users to query AWS Amplify Apps to retrieve detailed information about each application, including its name, ARN, creation date, default domain, and more."
---

# Table: aws_amplify_app - Query AWS Amplify Apps using SQL

The `aws_amplify_app` table in Steampipe provides information about apps within AWS Amplify. This table allows DevOps engineers to query app-specific details, including the name, ARN, creation date, last update date, default domain, and associated metadata. Users can utilize this table to gather insights on Amplify Apps, such as the apps' status, platform, repository, and more. The schema outlines the various attributes of the Amplify App, including the app ID, app ARN, platform, repository, production branch, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_amplify_app` table, you can use the `.inspect aws_amplify_app` command in Steampipe.

### Key columns:

- `app_id`: The unique identifier for an Amplify app. This can be used to join this table with other tables that reference Amplify apps.
- `app_arn`: The Amazon Resource Name (ARN) of the Amplify app. This column can be used to join with other tables when ARNs are available.
- `name`: The name of the Amplify app. This column is useful for queries that require the app's name.

## Examples

### Basic info

```sql
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

```sql
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

### List apps updated within the last hour

```sql
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

### Describe information about the production branch for an app

```sql
select
  production_branch ->> 'BranchName' as branch_name,
  production_branch ->> 'LastDeployTime' as last_deploy_time,
  production_branch ->> 'Status' as status
from
  aws_amplify_app
where
  name = 'amplify_app_name';
```

### List information about the build spec for an app

```sql
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

### List information on rewrite(200) redirect settings for an app

```sql
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

### List all apps that have branch auto build enabled

```sql
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
