# Table: aws_amplify_app

Represents the different branches of a repository for building, deploying, and hosting an Amplify app.

## Examples

### Basic info

```sql
select
  app_id,
  arn,
  auto_branch_creation_config,
  auto_branch_creation_patterns,
  basic_auth_credentials,
  build_spec_json,
  create_time,
  custom_headers,
  custom_rules,
  default_domain,
  description,
  enable_auto_branch_creation,
  enable_basic_auth,
  enable_branch_auto_build,
  enable_branch_auto_deletion,
  environment_variables,
  iam_service_role_arn,
  name,
  platform,
  production_branch,
  repository,
  repository_clone_method,
  update_time,
  title,
  tags
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

````sql
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
````

### List information about the build spec for an app

```sql
select
  name,
  app_id,
  build_spec_json ->> 'backend' as build_backend_spec,
  build_spec_json ->> 'frontend' as build_frontend_spec,
  build_spec_json ->> 'test' as build_test_spec,
  build_spec_json ->> 'env' as build_env_settings
from
  aws_amplify_app
where
  name = 'amplify_app_name';
```

### List information about the redirect settings for an app

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
