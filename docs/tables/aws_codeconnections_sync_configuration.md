---
title: "Steampipe Table: aws_codeconnections_sync_configuration - Query AWS CodeConnections sync configurations using SQL"
description: "Allows users to query AWS CodeConnections sync configurations, including their repositories, branches, roles, triggers, and deployment settings."
folder: "CodeConnections"
---

# Table: aws_codeconnections_sync_configuration - Query AWS CodeConnections sync configurations using SQL

AWS CodeConnections sync configurations define how changes in a linked Git repository are synchronized with an AWS resource. CloudFormation stack sync is currently the supported sync type.

## Table Usage Guide

The `aws_codeconnections_sync_configuration` table provides Git sync settings for repository links in each configured AWS region. Querying the table enumerates repository links and returns their sync configurations.

## Examples

### Basic info

```sql+postgres
select
  resource_name,
  repository_name,
  branch,
  sync_type,
  region
from
  aws_codeconnections_sync_configuration;
```

```sql+sqlite
select
  resource_name,
  repository_name,
  branch,
  sync_type,
  region
from
  aws_codeconnections_sync_configuration;
```

### Review deployment and update behavior

```sql+postgres
select
  resource_name,
  publish_deployment_status,
  pull_request_comment,
  trigger_resource_update_on
from
  aws_codeconnections_sync_configuration;
```

```sql+sqlite
select
  resource_name,
  publish_deployment_status,
  pull_request_comment,
  trigger_resource_update_on
from
  aws_codeconnections_sync_configuration;
```

### List configurations for a repository link

```sql+postgres
select
  resource_name,
  repository_name,
  branch,
  role_arn
from
  aws_codeconnections_sync_configuration
where
  repository_link_id = '6053346f-8a33-4edb-9397-10394b695173';
```

```sql+sqlite
select
  resource_name,
  repository_name,
  branch,
  role_arn
from
  aws_codeconnections_sync_configuration
where
  repository_link_id = '6053346f-8a33-4edb-9397-10394b695173';
```
