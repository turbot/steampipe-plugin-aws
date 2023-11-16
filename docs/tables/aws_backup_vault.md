---
title: "Table: aws_backup_vault - Query AWS Backup Vaults using SQL"
description: "Allows users to query AWS Backup Vaults, providing detailed information about each backup vault, including its name, ARN, recovery points, and more."
---

# Table: aws_backup_vault - Query AWS Backup Vaults using SQL

The `aws_backup_vault` table in Steampipe provides information about backup vaults within AWS Backup. This table allows DevOps engineers to query vault-specific details, including the vault name, ARN, number of recovery points, and associated metadata. Users can utilize this table to gather insights on backup vaults, such as the number of recovery points for each vault, the creation date of each vault, and more. The schema outlines the various attributes of the backup vault, including the vault name, ARN, creation date, last resource backup time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_backup_vault` table, you can use the `.inspect aws_backup_vault` command in Steampipe.

**Key columns**:

- `name`: The name of the backup vault. This can be used to join with other tables that contain backup vault names.
- `arn`: The Amazon Resource Name (ARN) of the backup vault. This can be used to join with other tables that contain backup vault ARNs.
- `recovery_points`: The number of recovery points that are stored in a backup vault. This can be useful for understanding the backup history of a vault.

## Examples

### Basic Info

```sql
select
  name,
  arn,
  creation_date
from
  aws_backup_vault;
```

### List vaults older than 90 days

```sql
select
  name,
  arn,
  creation_date
from
  aws_backup_vault
where
  creation_date <= (current_date - interval '90' day)
order by
  creation_date;
```

### List vaults that do not prevent the deletion of backups in the backup vault

```sql
select
  name
from
  aws_backup_vault,
  jsonb_array_elements(policy -> 'Statement') as s
where
  s ->> 'Principal' = '*'
  and s ->> 'Effect' != 'Deny'
  and s ->> 'Action' like '%DeleteBackupVault%';
```

### List policy details for backup vaults

```sql
select
  name,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_backup_vault;
```
