---
title: "Steampipe Table: aws_rds_pending_maintenance_action - Query AWS RDS DB Maintenance Actions using SQL"
description: "Allows users to query AWS RDS DB Maintenance Actions and retrieve valuable information about pending maintenance tasks for RDS DB instances and clusters."
---

# Table: aws_rds_pending_maintenance_action - Query AWS RDS DB Maintenance Actions using SQL

The AWS RDS DB Maintenance Action is a feature of Amazon Relational Database Service (RDS) that allows users to manage and schedule maintenance tasks for their RDS DB instances and clusters. These maintenance actions can include updates, patches, and other necessary changes to ensure optimal performance and security of the database services.

## Table Usage Guide

The `aws_rds_pending_maintenance_action` table in Steampipe provides you with information about pending maintenance actions for RDS DB instances and clusters. This table allows you, as a DevOps engineer, to query details about maintenance tasks that are scheduled or required for your RDS resources. You can utilize this table to gather insights on the nature of maintenance actions, their statuses, and timelines for application. The schema outlines the various attributes of the maintenance actions, including the resource identifier, action type, and relevant dates.

## Examples

### List of pending maintenance actions for RDS DB instances
Discover the pending maintenance actions that need to be addressed for your RDS DB instances. This is crucial for ensuring that your databases are up to date and secure.

```sql+postgres
select
  resource_identifier,
  action,
  opt_in_status,
  forced_apply_date,
  current_apply_date,
  auto_applied_after_date
from
  aws_rds_pending_maintenance_action;
```

### Check if a maintenance action is for a cluster
Determine if a specific maintenance action is associated with a DB cluster.

```sql+postgres
select
  resource_identifier,
  case
    when resource_identifier like '%:cluster:%' then true
    else false
  end as is_cluster
from
  aws_rds_pending_maintenance_action
where 
  resource_identifier like '%clustername';
```

```sql+sqlite
select
  resource_identifier,
  case
    when resource_identifier like '%:cluster:%' then 1
    else 0
  end as is_cluster
from
  aws_rds_pending_maintenance_action
where 
  resource_identifier like '%clustername';
```

### List DB cluster pending maintenance actions
Discover the segments that require pending maintenance actions in your database clusters. This is useful in planning and prioritizing maintenance schedules, by understanding which actions are due and their respective timelines.

```sql+postgres
select
  a.db_cluster_identifier,
  action,
  a.status,
  opt_in_status,
  forced_apply_date,
  current_apply_date,
  auto_applied_after_date
from 
  aws_rds_db_cluster as a
join 
  aws_rds_pending_maintenance_action as b 
on 
  a.arn = b.resource_identifier;
```