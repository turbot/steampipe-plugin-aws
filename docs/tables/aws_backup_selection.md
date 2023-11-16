---
title: "Table: aws_backup_selection - Query AWS Backup Selections using SQL"
description: "Allows users to query AWS Backup Selections to obtain detailed information about the backup selection resources within AWS Backup service."
---

# Table: aws_backup_selection - Query AWS Backup Selections using SQL

The `aws_backup_selection` table in Steampipe provides comprehensive information about backup selection resources within the AWS Backup service. This table allows DevOps engineers, security professionals, and system administrators to query backup selection-specific details, including the selection's ARN, backup plan ID, creation and modification dates, and associated creator request ID. Users can utilize this table to gather insights on backup selections, such as identifying backup selections associated with specific backup plans, tracking creation and modification times of backup selections, and more. The schema outlines the various attributes of the backup selection, including the backup selection ARN, backup plan ID, creation date, creator request ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_backup_selection` table, you can use the `.inspect aws_backup_selection` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the backup selection. This can be used to join with other tables that contain AWS resource ARNs.
- `backup_plan_id`: The ID of the backup plan that is associated with the backup selection. This can be used to join with the `aws_backup_plan` table.
- `creation_date`: The date and time when the backup selection was created. This can be useful for tracking the creation time of backup selections.

## Examples

### Basic Info

```sql
select
  selection_name,
  backup_plan_id,
  iam_role_arn,
  region,
  account_id
from
  aws_backup_selection;
```

### List EBS volumes that are in a backup plan

```sql
with filtered_data as (
  select
    backup_plan_id,
    jsonb_agg(r) as assigned_resource
  from
    aws_backup_selection,
    jsonb_array_elements(resources) as r
  group by backup_plan_id
)
select
  v.volume_id,
  v.region,
  v.account_id
from
  aws_ebs_volume as v
  join filtered_data t on t.assigned_resource ?| array[v.arn];
```
