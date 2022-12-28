# Table: aws_securitylake_data_lake

Amazon Security Lake automatically centralizes your security data from the cloud, on-premises, and custom sources into a data lake stored in your account.

### Examples

### Basic info

```sql
select
  encryption_key,
  replication_role_arn,
  s3_bucket_arn,
  status
from
  aws_securitylake_data_lake;
```

### Get S3 bucket details of each data lake

```sql
select
  distinct b.name as bucket_name,
  l.s3_bucket_arn,
  b.creation_date,
  b.bucket_policy_is_public,
  b.versioning_enabled,
  b.block_public_acls
from
  aws_securitylake_data_lake as l,
  aws_s3_bucket as b
where
  l.s3_bucket_arn = b.arn;
```

### Get retention setting details of data lakes

```sql
select
  l.encryption_key,
  l.replication_role_arn,
  l.s3_bucket_arn,
  l.status,
  r ->> 'RetentionPeriod' as retention_period,
  r ->> 'StorageClass' as storage_class
from
  aws_securitylake_data_lake as l,
  jsonb_array_elements(retention_settings) as r;
```

### List data lakes where the configuration operation is in a pending state

```sql
select
  encryption_key,
  replication_role_arn,
  s3_bucket_arn,
  status
from
  aws_securitylake_data_lake
where
  status <> 'COMPLETED';
```