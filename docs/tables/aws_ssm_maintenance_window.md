---
title: "Steampipe Table: aws_ssm_maintenance_window - Query AWS Systems Manager Maintenance Windows using SQL"
description: "Allows users to query AWS Systems Manager Maintenance Windows to retrieve details about scheduled maintenance tasks for AWS resources."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_maintenance_window - Query AWS Systems Manager Maintenance Windows using SQL

The AWS Systems Manager Maintenance Windows is a feature that allows you to define a schedule for when to perform potentially disruptive actions on your instances such as patching an operating system, updating drivers, or installing software or patches. During these windows, AWS Systems Manager performs the tasks you've assigned, and you can track tasks and executions in detail. It provides a safe and consistent method to apply patches and updates to your instances.

## Table Usage Guide

The `aws_ssm_maintenance_window` table in Steampipe provides you with information about Maintenance Windows within AWS Systems Manager. This table allows you, as a DevOps engineer, to query details about scheduled maintenance tasks for AWS resources, including the maintenance window ID, name, description, and schedule. You can utilize this table to gather insights on maintenance windows, such as their duration, cut-off time, and whether they are enabled or not. The schema outlines the various attributes of the maintenance window for you, including the window ID, ARN, owner, enabled status, priority, and associated tags.

## Examples

### Basic info
Determine the areas in which AWS System Manager's Maintenance Windows are enabled and scheduled. This is useful for understanding the operational status and schedule of maintenance tasks across different regions.

```sql+postgres
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

```sql+sqlite
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
This query is useful for gaining insights into each maintenance window's target details within your AWS Simple Systems Manager (SSM). This can help manage and schedule tasks on your resources more effectively.

```sql+postgres
select
  name,
  p ->> 'WindowTargetId' as window_target_id,
  p ->> 'ResourceType' as resource_type,
  p ->> 'Name' as target_name
from
  aws_ssm_maintenance_window,
  jsonb_array_elements(targets) as p;
```

```sql+sqlite
select
  name,
  json_extract(p.value, '$.WindowTargetId') as window_target_id,
  json_extract(p.value, '$.ResourceType') as resource_type,
  json_extract(p.value, '$.Name') as target_name
from
  aws_ssm_maintenance_window,
  json_each(targets) as p;
```


### Get tasks details for each maintenance window
Explore the specifics of tasks within each maintenance window in your AWS Simple Systems Manager (SSM) to better manage system maintenance and updates.

```sql+postgres
select
  name,
  p ->> 'WindowTaskId' as window_task_id,
  p ->> 'ServiceRoleArn' as service_role_arn,
  p ->> 'Name' as task_name
from
  aws_ssm_maintenance_window,
  jsonb_array_elements(tasks) as p;
```

```sql+sqlite
select
  name,
  json_extract(p.value, '$.WindowTaskId') as window_task_id,
  json_extract(p.value, '$.ServiceRoleArn') as service_role_arn,
  json_extract(p.value, '$.Name') as task_name
from
  aws_ssm_maintenance_window,
  json_each(tasks) as p;
```


### List maintenance windows that are enabled
Identify the active maintenance windows within your AWS environment. This can help in planning system updates or troubleshooting activities without disrupting regular operations.

```sql+postgres
select
  name,
  window_id,
  enabled
from
  aws_ssm_maintenance_window
where
  enabled;
```

```sql+sqlite
select
  name,
  window_id,
  enabled
from
  aws_ssm_maintenance_window
where
  enabled = 1;
```