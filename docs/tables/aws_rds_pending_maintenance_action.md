---
title: "Steampipe Table: aws_rds_pending_maintenance_action - Query AWS RDS DB Maintenance Actions using SQL"
description: "Allows users to query AWS RDS DB Maintenance Actions and retrieve valuable information about pending maintenance tasks for RDS DB instances and clusters."
folder: "RDS"
---

# Table: aws_rds_pending_maintenance_action - Query AWS RDS DB Maintenance Actions using SQL

The AWS RDS Pending Maintenance Action is a feature of Amazon Relational Database Service (RDS) that allows users to manage and schedule maintenance tasks for their RDS DB instances and clusters. These maintenance actions can include updates, patches, and other necessary changes to ensure optimal performance and security of the database services.

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

```sql+sqlite
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

### Determine if a maintenance action is for clusters
Check whether a specific maintenance action is associated with a DB cluster.

```sql+postgres
select
  resource_identifier,
  case
    when resource_identifier like '%:cluster:%' then true
    else false
  end as is_cluster
from
  aws_rds_pending_maintenance_action;
```

```sql+sqlite
select
  resource_identifier,
  case
    when resource_identifier like '%:cluster:%' then 1
    else 0
  end as is_cluster
from
  aws_rds_pending_maintenance_action;
```

### List DB clusters pending maintenance actions
Identify DB clusters pending maintenance actions to plan and prioritize maintenance schedules effectively.

```sql+postgres
select
  a.db_cluster_identifier,
  b.action,
  a.status,
  b.opt_in_status,
  b.forced_apply_date,
  b.current_apply_date,
  b.auto_applied_after_date
from 
  aws_rds_db_instance as a
  join aws_rds_pending_maintenance_action as b on b.resource_identifier = a.arn;
```

```sql+sqlite
select
  a.db_cluster_identifier,
  b.action,
  a.status,
  b.opt_in_status,
  b.forced_apply_date,
  b.current_apply_date,
  b.auto_applied_after_date
from 
  aws_rds_db_instance as a
  join aws_rds_pending_maintenance_action as b on b.resource_identifier = a.arn;
```

### Retrieve pending maintenance actions for a specific cluster
Fetch pending maintenance actions for a specific DB cluster using its Amazon Resource Name (ARN).

```sql+postgres
select
  resource_identifier,
  action,
  current_apply_date,
  forced_apply_date
from
  aws_rds_pending_maintenance_action
where
  resource_identifier = 'arn:aws:rds:us-east-1:123456789012:cluster:my-aurora-cluster-1';
``` 

```sql+sqlite
select
  resource_identifier,
  action,
  current_apply_date,
  forced_apply_date
from
  aws_rds_pending_maintenance_action
where
  resource_identifier = 'arn:aws:rds:us-east-1:123456789012:cluster:my-aurora-cluster-1';
``` 