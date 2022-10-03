# Table: aws_kms_key

AWS Key Management Service (KMS) is an Amazon Web Services product that allows administrators to create, delete and control keys that encrypt data stored in AWS databases and products.

## Examples

### Basic info

```sql
select
  id,
  title,
  arn,
  key_manager,
  creation_date
from
  aws_kms_key;
```

### List of KMS keys where key rotation is not enabled

```sql
select
  id,
  key_rotation_enabled
from
  aws_kms_key
where
  not key_rotation_enabled;
```


### List of KMS Customer Managed keys scheduled for deletion

```sql
select
  id,
  key_state,
  deletion_date
from
  aws_kms_key
where
  key_state = 'PendingDeletion';
```


### List of unused Customer Managed Keys

```sql
select
  id,
  enabled as key_enabled
from
  aws_kms_key
where
  not enabled;
```


### Count of AWS KMS keys by Key manager

```sql
select
  key_manager,
  count(key_manager) as count
from
  aws_kms_key
group by
  key_manager;
```