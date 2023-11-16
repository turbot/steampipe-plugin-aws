---
title: "Table: aws_backup_plan - Query AWS Backup Plan using SQL"
description: "Allows users to query AWS Backup Plan data, providing detailed information about each backup plan created within an AWS account. Useful for DevOps engineers to monitor and manage backup strategies and ensure data recovery processes are in place."
---

# Table: aws_backup_plan - Query AWS Backup Plan using SQL

The `aws_backup_plan` table in Steampipe provides information about each backup plan within AWS Backup. This table allows DevOps engineers to query backup plan-specific details, including backup options, creation and version details, and associated metadata. Users can utilize this table to gather insights on backup plans, such as the backup frequency, backup window, lifecycle of the backup, and more. The schema outlines the various attributes of the backup plan, including the backup plan ARN, creation date, version, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_backup_plan` table, you can use the `.inspect aws_backup_plan` command in Steampipe.

**Key columns**:

- `name`: This column provides the name of the backup plan. It is a key column for joining with other tables that reference the backup plan by name.
- `arn`: The Amazon Resource Name (ARN) of the backup plan. This unique identifier is useful for joining data across AWS services.
- `version_id`: This column provides the unique version ID of the backup plan. It can be used to track changes to the backup plan over time.

## Examples

### Basic Info

```sql
select
  name,
  backup_plan_id,
  arn,
  creation_date,
  last_execution_date
from
  aws_backup_plan;
```

### List plans older than 90 days

```sql
select
  name,
  backup_plan_id,
  arn,
  creation_date,
  last_execution_date
from
  aws_backup_plan
where
  creation_date <= (current_date - interval '90' day)
order by
  creation_date;
```

### List plans that were deleted in the last 7 days

```sql
select
  name,
  arn,
  creation_date,
  deletion_date
from
  aws_backup_plan
where
  deletion_date > current_date - 7
order by
  deletion_date;
```
