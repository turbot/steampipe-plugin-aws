---
title: "Table: aws_backup_legal_hold - Query AWS Backup Legal Hold using SQL"
description: "Allows users to query AWS Backup Legal Hold to obtain information about the legal hold settings of AWS backup resources."
---

# Table: aws_backup_legal_hold - Query AWS Backup Legal Hold using SQL

The `aws_backup_legal_hold` table in Steampipe provides information about legal hold settings for AWS Backup resources. This table allows DevOps engineers to query legal hold-specific details, including the backup resource ARN, the legal hold status, and the last update time. Users can utilize this table to review and monitor the legal hold status of backup resources, ensuring compliance with data retention policies and legal requirements. The schema outlines the various attributes of the legal hold, including the backup resource ARN, the legal hold status, and the last update time.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_backup_legal_hold` table, you can use the `.inspect aws_backup_legal_hold` command in Steampipe.

### Key columns:

- `arn`: The ARN of the backup resource. This can be used to join with other tables that contain AWS backup resource information.
- `legal_hold`: The legal hold status. This is crucial for understanding whether a backup resource is under legal hold or not.
- `last_update_time`: The time when the legal hold status was last updated. This can be useful for tracking changes over time.

## Examples

### Basic Info

```sql
select
  legal_hold_id,
  arn,
  creation_date,
  cancellation_date
from
  aws_backup_legal_hold;
```

### List legal holds older than 10 days

```sql
select
  legal_hold_id,
  arn,
  creation_date,
  creation_date,
  retain_record_until
from
  aws_backup_legal_hold
where
  creation_date <= current_date - interval '10' day
order by
  creation_date;
```

### Get recovery point selection details for each legal hold

```sql
select
  title,
  legal_hold_id,
  recovery_point_selection -> 'DateRange' ->> 'ToDate' as to_date,
  recovery_point_selection -> 'DateRange' ->> 'FromDate' as from_date,
  recovery_point_selection -> 'VaultNames' as vault_names,
  recovery_point_selection ->> 'ResourceIdentifiers' as resource_identifiers
from
  aws_backup_legal_hold;
```
