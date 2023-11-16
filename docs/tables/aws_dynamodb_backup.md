---
title: "Table: aws_dynamodb_backup - Query AWS DynamoDB Backup using SQL"
description: "Allows users to query DynamoDB Backup details such as backup ARN, backup creation date, backup size, backup status, and more."
---

# Table: aws_dynamodb_backup - Query AWS DynamoDB Backup using SQL

The `aws_dynamodb_backup` table in Steampipe provides information about backups in AWS DynamoDB. This table allows DevOps engineers to query backup-specific details, including backup ARN, backup creation date, backup size, backup status, and more. Users can utilize this table to gather insights on backups, such as backup status, backup type, size in bytes, and more. The schema outlines the various attributes of the DynamoDB backup, including the backup ARN, backup creation date, backup size, backup status, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dynamodb_backup` table, you can use the `.inspect aws_dynamodb_backup` command in Steampipe.

### Key columns:

- `backup_arn`: The Amazon Resource Name (ARN) associated with the backup. This can be used to join with other tables that use backup ARN.
- `table_name`: The name of the table associated with the backup. This can be used to join with other tables that use table name.
- `backup_name`: The name of the backup. This can be used to join with other tables that use backup name.

## Examples

### List backups with their corresponding tables

```sql
select
  name,
  table_name,
  table_id
from
  aws_dynamodb_backup;
```


### Basic backup info

```sql
select
  name,
  backup_status,
  backup_type,
  backup_expiry_datetime,
  backup_size_bytes
from
  aws_dynamodb_backup;
```