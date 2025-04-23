---
title: "Steampipe Table: aws_dynamodb_backup - Query AWS DynamoDB Backup using SQL"
description: "Allows users to query DynamoDB Backup details such as backup ARN, backup creation date, backup size, backup status, and more."
folder: "Backup"
---

# Table: aws_dynamodb_backup - Query AWS DynamoDB Backup using SQL

The AWS DynamoDB Backup service provides on-demand and continuous backups of your DynamoDB tables, safeguarding your data for archival and disaster recovery. It enables point-in-time recovery, allowing you to restore your table data from any second in the past 35 days. This service also supports backup and restore actions through AWS Management Console, AWS CLI, and AWS SDKs.

## Table Usage Guide

The `aws_dynamodb_backup` table in Steampipe provides you with information about backups in AWS DynamoDB. This table allows you, as a DevOps engineer, to query backup-specific details, including backup ARN, backup creation date, backup size, backup status, and more. You can utilize this table to gather insights on backups, such as backup status, backup type, size in bytes, and more. The schema outlines the various attributes of the DynamoDB backup for you, including the backup ARN, backup creation date, backup size, backup status, and more.

## Examples

### List backups with their corresponding tables
Determine the areas in which backups are associated with their corresponding tables in the AWS DynamoDB service. This can be useful for understanding the relationship between backups and tables, aiding in efficient data management and disaster recovery planning.

```sql+postgres
select
  name,
  table_name,
  table_id
from
  aws_dynamodb_backup;
```

```sql+sqlite
select
  name,
  table_name,
  table_id
from
  aws_dynamodb_backup;
```


### Basic backup info
Assess the elements within your AWS DynamoDB backup, such as status, type, expiry date, and size. This allows you to manage and optimize your backup strategy effectively.

```sql+postgres
select
  name,
  backup_status,
  backup_type,
  backup_expiry_datetime,
  backup_size_bytes
from
  aws_dynamodb_backup;
```

```sql+sqlite
select
  name,
  backup_status,
  backup_type,
  backup_expiry_datetime,
  backup_size_bytes
from
  aws_dynamodb_backup;
```