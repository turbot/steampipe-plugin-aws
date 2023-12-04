# Table: aws_ebs_snapshot

Describes whether EBS encryption by default is enabled for your account in the current Region.

## Examples

### Basic info

```sql
select
  ebs_encryption_by_default,
  ebs_default_kms_key_id,
  region,
  account_id
from
  aws_ebs_encryption_by_default;
```

### Get regions with default encryption enabled

```sql
select
  ebs_encryption_by_default,
  ebs_default_kms_key_id,
  region
from
  aws_ebs_encryption_by_default
where ebs_encryption_by_default;
```