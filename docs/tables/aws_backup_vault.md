# Table: aws_backup_vault

AWS Backup vault is a container that you organize your backups in. You can use backup vaults to set the AWS Key Management Service (AWS KMS) encryption key that is used to encrypt backups in the backup vault and to control access to the backups in the backup vault.

If you require different encryption keys or access policies for different groups of backups, you can optionally create multiple backup vaults. Otherwise, you can have all your backups organized in the default backup vault.

**Note:** If you pass in the name of a vault in where clause that doesn't exist or you don't have access to (through the vault's access policy), you will get no rows.

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
