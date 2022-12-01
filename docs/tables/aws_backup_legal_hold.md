# Table: aws_backup_legal_hold

A legal hold is an administrative tool that helps prevent backups from being deleted while under a hold. While the hold is in place, backups under a hold cannot be deleted and lifecycle policies that would alter the backup status (such as transition to a Deleted state) are delayed until the legal hold is removed.

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

### Get recovery point selection details for legal holds

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
