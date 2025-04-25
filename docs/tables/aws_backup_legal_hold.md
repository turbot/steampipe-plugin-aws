---
title: "Steampipe Table: aws_backup_legal_hold - Query AWS Backup Legal Hold using SQL"
description: "Allows users to query AWS Backup Legal Hold to obtain information about the legal hold settings of AWS backup resources."
folder: "Backup"
---

# Table: aws_backup_legal_hold - Query AWS Backup Legal Hold using SQL

The AWS Backup Legal Hold is a feature within the AWS Backup service that helps in preserving your backup recovery points, preventing them from being accidentally or intentionally deleted. It provides an additional layer of data protection by enabling you to enforce a legal hold on backup recovery points, regardless of the retention period. This feature is particularly useful in legal and compliance scenarios where data retention is of utmost importance.

## Table Usage Guide

The `aws_backup_legal_hold` table in Steampipe provides you with information about legal hold settings for AWS Backup resources. This table allows you, as a DevOps engineer, to query legal hold-specific details, including the backup resource ARN, the legal hold status, and the last update time. You can utilize this table to review and monitor the legal hold status of backup resources, ensuring compliance with data retention policies and legal requirements. The schema outlines for you the various attributes of the legal hold, including the backup resource ARN, the legal hold status, and the last update time.

## Examples

### Basic Info
Explore the instances in your AWS backup where a legal hold has been applied. This can help you understand when and where these holds were created, as well as when they were cancelled, providing valuable insights for audit or compliance purposes.

```sql+postgres
select
  legal_hold_id,
  arn,
  creation_date,
  cancellation_date
from
  aws_backup_legal_hold;
```

```sql+sqlite
select
  legal_hold_id,
  arn,
  creation_date,
  cancellation_date
from
  aws_backup_legal_hold;
```

### List legal holds older than 10 days
Determine the areas in which legal holds on your AWS backup have been in place for more than 10 days. This can help you manage your resources more effectively by identifying outdated holds that may no longer be necessary.

```sql+postgres
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

```sql+sqlite
select
  legal_hold_id,
  arn,
  creation_date,
  creation_date,
  retain_record_until
from
  aws_backup_legal_hold
where
  creation_date <= date('now','-10 day')
order by
  creation_date;
```

### Get recovery point selection details for each legal hold
Explore the specific periods and resources associated with each legal hold in your AWS backup system. This can be useful to understand the scope and duration of your data recovery points under legal holds.

```sql+postgres
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

```sql+sqlite
select
  title,
  legal_hold_id,
  json_extract(recovery_point_selection, '$.DateRange.ToDate') as to_date,
  json_extract(recovery_point_selection, '$.DateRange.FromDate') as from_date,
  json_extract(recovery_point_selection, '$.VaultNames') as vault_names,
  json_extract(recovery_point_selection, '$.ResourceIdentifiers') as resource_identifiers
from
  aws_backup_legal_hold;
```