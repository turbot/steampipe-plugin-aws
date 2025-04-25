---
title: "Steampipe Table: aws_backup_selection - Query AWS Backup Selections using SQL"
description: "Allows users to query AWS Backup Selections to obtain detailed information about the backup selection resources within AWS Backup service."
folder: "Backup"
---

# Table: aws_backup_selection - Query AWS Backup Selections using SQL

The AWS Backup Selection is a component of AWS Backup, a fully managed backup service that simplifies the backup of data across AWS services. It allows you to automate and centrally manage backups, enforcing policies and monitoring backup activities for AWS resources. The selection includes a list of resources to be backed up, identified by an array of ARNs, as well as a backup plan to specify how AWS Backup handles backup and restore operations.

## Table Usage Guide

The `aws_backup_selection` table in Steampipe provides you with comprehensive information about backup selection resources within the AWS Backup service. This table allows you, as a DevOps engineer, security professional, or system administrator, to query backup selection-specific details, including the selection's ARN, backup plan ID, creation and modification dates, and associated creator request ID. You can utilize this table to gather insights on backup selections, such as identifying backup selections associated with specific backup plans, tracking creation and modification times of backup selections, and more. The schema outlines the various attributes of the backup selection for you, including the backup selection ARN, backup plan ID, creation date, creator request ID, and associated tags.

## Examples

### Basic Info
Explore which AWS backup plans are associated with specific IAM roles and regions. This can be useful for auditing and managing your AWS resources efficiently.

```sql+postgres
select
  selection_name,
  backup_plan_id,
  iam_role_arn,
  region,
  account_id
from
  aws_backup_selection;
```

```sql+sqlite
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
Identify the EBS volumes included in a backup plan to ensure crucial data is secured and maintained. This is essential for data recovery planning and to minimize potential data loss.

```sql+postgres
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

```sql+sqlite
Error: SQLite does not support the ?| operator used in array operations.
```