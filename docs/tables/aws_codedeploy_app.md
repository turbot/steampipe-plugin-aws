---
title: "Steampipe Table: aws_codedeploy_app - Query AWS CodeDeploy Applications using SQL"
description: "Allows users to query AWS CodeDeploy Applications to return detailed information about each application, including application name, ID, and associated deployment groups."
folder: "CodeDeploy"
---

# Table: aws_codedeploy_app - Query AWS CodeDeploy Applications using SQL

The AWS CodeDeploy service automates code deployments to any instance, including Amazon EC2 instances and instances running on-premises. An Application in AWS CodeDeploy is a name that uniquely identifies the application you want to deploy. AWS CodeDeploy uses this name, which functions like a container, to ensure the correct combination of revision, deployment configuration, and deployment group are referenced during a deployment.

## Table Usage Guide

The `aws_codedeploy_app` table in Steampipe provides you with information about applications within AWS CodeDeploy. This table allows you, as a DevOps engineer, to query application-specific details, including application name, compute platform, and linked deployment groups. You can utilize this table to gather insights on applications, such as their deployment configurations, linked deployment groups, and compute platforms. The schema outlines the various attributes of the CodeDeploy application for you, including the application name, application ID, and the linked deployment groups.

## Examples

### Basic info
Explore the deployment applications in your AWS environment to understand their creation time and associated computing platform. This is beneficial for tracking the history and configuration of your applications across different regions.

```sql+postgres
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

```sql+sqlite
select
  arn,
  application_id,
  application_name,
  compute_platform,
  create_time,
  region
from
  aws_codedeploy_app;
```

### Get total applications deployed on each platform
Explore the distribution of applications across various platforms to better understand your deployment strategy. This can assist in identifying platforms that are heavily utilized for deploying applications, aiding in resource allocation and management decisions.

```sql+postgres
select
  count(arn) as application_count,
  compute_platform
from
  aws_codedeploy_app
group by
  compute_platform;
```

```sql+sqlite
select
  count(arn) as application_count,
  compute_platform
from
  aws_codedeploy_app
group by
  compute_platform;
```

### List applications linked to GitHub
Identify instances where applications are linked to GitHub within the AWS CodeDeploy service. This is useful for gaining insights into the integration between your applications and GitHub, which can help in managing and troubleshooting your deployment processes.

```sql+postgres
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

```sql+sqlite
select
  arn,
  application_id,
  compute_platform,
  create_time,
  github_account_name
from
  aws_codedeploy_app
where
  linked_to_github = 1;
```