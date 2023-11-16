---
title: "Table: aws_backup_recovery_point - Query AWS Backup Recovery Points using SQL"
description: "Allows users to query AWS Backup Recovery Points to gather comprehensive information about each recovery point within an AWS Backup vault."
---

# Table: aws_backup_recovery_point - Query AWS Backup Recovery Points using SQL

The `aws_backup_recovery_point` table in Steampipe provides information about each recovery point within an AWS Backup vault. This table allows DevOps engineers and system administrators to query recovery point-specific details, including the backup vault where the recovery point is stored, the source of the backup, the state of the recovery point, and associated metadata. Users can utilize this table to gather insights on recovery points, such as identifying unencrypted recovery points, verifying backup completion status, and more. The schema outlines the various attributes of the recovery point, including the recovery point ARN, creation date, backup size, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_backup_recovery_point` table, you can use the `.inspect aws_backup_recovery_point` command in Steampipe.

Key columns:

- `recovery_point_arn`: The Amazon Resource Name (ARN) that uniquely identifies the recovery point. This can be used to join with other tables that require a recovery point ARN.
- `backup_vault_name`: The name of the backup vault where the recovery point is stored. This can be used to join with other tables that require a backup vault name.
- `source_backup_vault_arn`: The ARN of the backup vault where the source backup was created. This can be used to join with other tables that require a source backup vault ARN.

## Examples

### Basic Info

```sql
select
  backup_vault_name,
  recovery_point_arn,
  resource_type,
  status
from
  aws_backup_recovery_point;
```

### List encrypted recovery points

```sql
select
  backup_vault_name,
  recovery_point_arn,
  resource_type,
  status,
  is_encrypted
from
  aws_backup_recovery_point
where
  is_encrypted;
```
