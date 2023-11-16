---
title: "Table: aws_codedeploy_deployment_group - Query AWS CodeDeploy Deployment Group using SQL"
description: "Allows users to query AWS CodeDeploy Deployment Group details including deployment configurations, target revisions, and associated alarm configurations."
---

# Table: aws_codedeploy_deployment_group - Query AWS CodeDeploy Deployment Group using SQL

The `aws_codedeploy_deployment_group` table in Steampipe provides information about deployment groups within AWS CodeDeploy. This table allows DevOps engineers to query deployment group-specific details, including deployment configurations, target revisions, and associated alarm configurations. Users can utilize this table to gather insights on deployment groups, such as deployment configuration names, target revisions, and alarm configurations. The schema outlines the various attributes of the deployment group, including the deployment group name, service role ARN, deployment configuration name, target revision, and associated alarm configurations.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codedeploy_deployment_group` table, you can use the `.inspect aws_codedeploy_deployment_group` command in Steampipe.

Key columns:

- `deployment_group_name`: The name of the deployment group. It is a key identifier and can be used to join with other tables to fetch detailed information.
- `deployment_config_name`: The name of the deployment configuration associated with the deployment group. It can be used to join with the deployment configuration table.
- `service_role_arn`: The service role ARN associated with the deployment group. It can be used to join with IAM role tables for more detailed role information.

## Examples

### Basic info

```sql
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

```sql
select
  count(arn) as group_count,
  compute_platform
from
  aws_codedeploy_deployment_group
group by
  compute_platform;
```

### List the last successful deployment for each deployment group

```sql
select
  arn,
  deployment_group_id,
  last_successful_deployment
from
  aws_codedeploy_deployment_group;
```

### Get total deployment groups based on deployment style

```sql
select
  count(arn) as group_count,
  deployment_style
from
  aws_codedeploy_deployment_group
group by
  deployment_style;
```

### List the deployment groups having automatic rollback enabled

```sql
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

### List all autoscaling groups in a particular deployment group for an application

```sql
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

### List the deployment groups having automatic rollback enabled

```sql
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
