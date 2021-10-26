# Table: aws_backup_recovery_point

A backup, or recovery point, represents the content of a resource, such as an Amazon Elastic Block Store (Amazon EBS) volume or Amazon DynamoDB table, at a specified time. Recovery point is a term that refers generally to the different backups in AWS services, such as Amazon EBS snapshots and DynamoDB backups.

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
