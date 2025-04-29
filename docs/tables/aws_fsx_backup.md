---
title: "Steampipe Table: aws_fsx_backup - Query AWS FSx Backups using SQL"
description: "Allows users to query AWS FSx Backups for detailed information about each backup, including its ID, type, lifecycle status, and associated metadata."
folder: "FSx"
---

# Table: aws_fsx_backup - Query AWS FSx Backups using SQL

The AWS FSx Backup is a feature of Amazon FSx that enables you to create point-in-time backups of your file systems. These backups can be used to restore your file system to a previous state, protect against data loss, and meet compliance requirements. AWS FSx automatically creates daily backups of your file systems, and you can also create manual backups as needed.

## Table Usage Guide

The `aws_fsx_backup` table in Steampipe provides you with information about AWS FSx backups. This table allows you as a DevOps engineer to query backup-specific details, including the backup ID, type, lifecycle status, creation time, and associated metadata. You can utilize this table to gather insights on backups, such as backup status, progress, and any failures that occurred during backup creation. The schema outlines the various attributes of the AWS FSx backup for you, including the backup ARN, file system ID, KMS key ID, and associated tags.

## Examples

### Basic info
Explore the basic information about your FSx backups, including their IDs, types, and lifecycle status. This can help you understand the current state of your backups and identify any that might need attention.

```sql+postgres
select
  backup_id,
  backup_type,
  lifecycle,
  creation_time,
  file_system_id
from
  aws_fsx_backup;
```

```sql+sqlite
select
  backup_id,
  backup_type,
  lifecycle,
  creation_time,
  file_system_id
from
  aws_fsx_backup;
```

### List failed backups
Identify backups that have failed, along with their failure details. This can help you troubleshoot issues and ensure your backup strategy is working effectively.

```sql+postgres
select
  backup_id,
  file_system_id,
  lifecycle,
  failure_details
from
  aws_fsx_backup
where
  lifecycle = 'FAILED';
```

```sql+sqlite
select
  backup_id,
  file_system_id,
  lifecycle,
  failure_details
from
  aws_fsx_backup
where
  lifecycle = 'FAILED';
```

### List backups for a specific file system
Find all backups associated with a particular file system. This can help you track the backup history and ensure you have sufficient backups for disaster recovery.

```sql+postgres
select
  backup_id,
  backup_type,
  lifecycle,
  creation_time
from
  aws_fsx_backup
where
  file_system_id = 'fs-1234567890abcdef0';
```

```sql+sqlite
select
  backup_id,
  backup_type,
  lifecycle,
  creation_time
from
  aws_fsx_backup
where
  file_system_id = 'fs-1234567890abcdef0';
```

### List backups with specific tags
Find backups that have specific tags associated with them. This can help you organize and manage your backups based on custom criteria.

```sql+postgres
select
  backup_id,
  file_system_id,
  tags
from
  aws_fsx_backup
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  backup_id,
  file_system_id,
  tags
from
  aws_fsx_backup
where
  json_extract(tags, '$.Environment') = 'Production';
```

### List backups without encryption
Identify backups that are not encrypted using KMS keys. This can help ensure compliance with security requirements.

```sql+postgres
select
  backup_id,
  file_system_id,
  kms_key_id,
  lifecycle
from
  aws_fsx_backup
where
  kms_key_id is null;
```

```sql+sqlite
select
  backup_id,
  file_system_id,
  kms_key_id,
  lifecycle
from
  aws_fsx_backup
where
  kms_key_id is null;
```

### List backups older than 30 days
Find backups that are older than 30 days to help manage backup retention and storage costs.

```sql+postgres
select
  backup_id,
  file_system_id,
  creation_time,
  lifecycle
from
  aws_fsx_backup
where
  creation_time < now() - interval '30 days';
```

```sql+sqlite
select
  backup_id,
  file_system_id,
  creation_time,
  lifecycle
from
  aws_fsx_backup
where
  creation_time < datetime('now', '-30 days');
```

### List in-progress backups
Monitor backups that are currently in progress to track backup operations.

```sql+postgres
select
  backup_id,
  file_system_id,
  progress_percent,
  creation_time
from
  aws_fsx_backup
where
  lifecycle = 'CREATING';
```

```sql+sqlite
select
  backup_id,
  file_system_id,
  progress_percent,
  creation_time
from
  aws_fsx_backup
where
  lifecycle = 'CREATING';
```

### List backups by type and status
Analyze backup distribution by type and status to understand backup patterns and identify potential issues.

```sql+postgres
select
  backup_type,
  lifecycle,
  count(*) as backup_count
from
  aws_fsx_backup
group by
  backup_type,
  lifecycle
order by
  backup_count desc;
```

```sql+sqlite
select
  backup_type,
  lifecycle,
  count(*) as backup_count
from
  aws_fsx_backup
group by
  backup_type,
  lifecycle
order by
  backup_count desc;
```

### List backups with detailed failure information
Get detailed information about failed backups to help with troubleshooting.

```sql+postgres
select
  backup_id,
  file_system_id,
  lifecycle,
  failure_details ->> 'Message' as failure_message,
  failure_details ->> 'Type' as failure_type
from
  aws_fsx_backup
where
  failure_details is not null;
```

```sql+sqlite
select
  backup_id,
  file_system_id,
  lifecycle,
  json_extract(failure_details, '$.Message') as failure_message,
  json_extract(failure_details, '$.Type') as failure_type
from
  aws_fsx_backup
where
  failure_details is not null;
``` 