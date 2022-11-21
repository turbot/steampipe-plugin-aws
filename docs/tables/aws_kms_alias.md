# Table: aws_kms_key

The AWS Key Management Service (KMS) Alias resource specifies a display name for a KMS key. You can use an alias to identify a KMS key in the AWS KMS console, in the DescribeKey operation, and in cryptographic operations, such as Decrypt and GenerateDataKey.

## Examples

### Basic info

```sql
select
  alias_name,
  title,
  arn,
  target_key_id,
  creation_date
from
  aws_kms_alias;
```

### List of KMS keys where key rotation disabled for the alias

```sql
select
  k.id as key_id,
  k.key_rotation_enabled as key_rotation_enabled,
  a.alias_name as alias_name,
  a.arn as alias_arn
from
  aws_kms_key as k,
  aws_kms_alias as a
where
  k.id = a.target_key_id and not key_rotation_enabled;
```


### List of KMS Customer Managed keys scheduled for deletion

```sql
select
  k.id as key_id,
  k.key_state as key_state,
  k.deletion_date as key_deletion_date,
  a.alias_name as alias_name
from
  aws_kms_key as k,
  aws_kms_alias as a
where
  k.id = a.target_key_id and key_state = 'PendingDeletion';
```

### Count of alias by Key id

```sql
select
  k.id as key_id,
  count(a.alias_name) as count
from
  aws_kms_key as k
  left join aws_kms_alias as a
  on k.id = a.target_key_id
group by
  key_id;
```