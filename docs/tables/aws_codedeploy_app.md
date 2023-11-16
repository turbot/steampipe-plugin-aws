---
title: "Table: aws_codedeploy_app - Query AWS CodeDeploy Applications using SQL"
description: "Allows users to query AWS CodeDeploy Applications to return detailed information about each application, including application name, ID, and associated deployment groups."
---

# Table: aws_codedeploy_app - Query AWS CodeDeploy Applications using SQL

The `aws_codedeploy_app` table in Steampipe provides information about applications within AWS CodeDeploy. This table allows DevOps engineers to query application-specific details, including application name, compute platform, and linked deployment groups. Users can utilize this table to gather insights on applications, such as their deployment configurations, linked deployment groups, and compute platforms. The schema outlines the various attributes of the CodeDeploy application, including the application name, application ID, and the linked deployment groups.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codedeploy_app` table, you can use the `.inspect aws_codedeploy_app` command in Steampipe.

**Key columns**:

- `application_name`: The name of the application. This can be used to join with other tables that require the application name.
- `application_id`: The unique ID of the application. This can be used to join with other tables that require the application ID.
- `compute_platform`: The platform type (e.g., ECS, Lambda, or Server) of the application. This can be used to join with other tables that require the compute platform.

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
