# Table: aws_codedeploy_app

An application is a name that uniquely identifies the application you want to deploy. CodeDeploy uses this name, which functions as a container, to ensure the correct combination of revision, deployment configuration, and deployment group are referenced during a deployment.

## Examples

### Basic info

```sql
select
  arn,
  application_id,
  application_name
  compute_platform,
  create_time,
  region
from
  aws_codedeploy_app;
```

### Get total applications deployed on each platform

```sql
select
  count(arn) as application_count,
  compute_platform
from
  aws_codedeploy_app
group by
  compute_platform;
```

### List applications linked to GitHub

```sql
select
  arn,
  application_id,
  compute_platform,
  create_time,
  github_account_name
from
  aws_codedeploy_app
where
  linked_to_github;
```
