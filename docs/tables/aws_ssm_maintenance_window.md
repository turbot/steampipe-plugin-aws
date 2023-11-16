---
title: "Table: aws_ssm_maintenance_window - Query AWS Systems Manager Maintenance Windows using SQL"
description: "Allows users to query AWS Systems Manager Maintenance Windows to retrieve details about scheduled maintenance tasks for AWS resources."
---

# Table: aws_ssm_maintenance_window - Query AWS Systems Manager Maintenance Windows using SQL

The `aws_ssm_maintenance_window` table in Steampipe provides information about Maintenance Windows within AWS Systems Manager. This table allows DevOps engineers to query details about scheduled maintenance tasks for AWS resources, including the maintenance window ID, name, description, and schedule. Users can utilize this table to gather insights on maintenance windows, such as their duration, cut-off time, and whether they are enabled or not. The schema outlines the various attributes of the maintenance window, including the window ID, ARN, owner, enabled status, priority, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_maintenance_window` table, you can use the `.inspect aws_ssm_maintenance_window` command in Steampipe.

**Key columns**:

- `window_id`: The ID of the maintenance window. This can be used to join with other tables that reference the maintenance window ID.
- `name`: The name of the maintenance window. This can be useful for identifying specific maintenance windows.
- `enabled`: Indicates whether the maintenance window is enabled. This can be useful for filtering active vs inactive maintenance windows.

## Examples

### Basic info

```sql
select
  name,
  window_id,
  enabled,
  schedule,
  tags_src,
  region
from
  aws_ssm_maintenance_window;
```


### Get target details for each maintenance window

```sql
select
  name,
  p ->> 'WindowTargetId' as window_target_id,
  p ->> 'ResourceType' as resource_type,
  p ->> 'Name' as target_name
from
  aws_ssm_maintenance_window,
  jsonb_array_elements(targets) as p;
```


### Get tasks details for each maintenance window

```sql
select
  name,
  p ->> 'WindowTaskId' as window_task_id,
  p ->> 'ServiceRoleArn' as service_role_arn,
  p ->> 'Name' as task_name
from
  aws_ssm_maintenance_window,
  jsonb_array_elements(tasks) as p;
```


### List maintenance windows that are enabled

```sql
select
  name,
  window_id,
  enabled
from
  aws_ssm_maintenance_window
where
  enabled;
```
