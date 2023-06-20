# Table: aws_s3_bucket

An Amazon S3 bucket is a public cloud storage resource available in Amazon Web Services' (AWS) Simple Storage Service (S3), an object storage offering.

## Examples

### Basic info

```sql
select
  name,
  region,
  account_id,
  bucket_policy_is_public
from
  aws_s3_bucket;
```

### List buckets with versioning disabled

```sql
select
  name,
  region,
  account_id,
  versioning_enabled
from
  aws_s3_bucket
where
  not versioning_enabled;
```

### List buckets with default encryption disabled

```sql
select
  name,
  server_side_encryption_configuration
from
  aws_s3_bucket
where
  server_side_encryption_configuration is null;
```

### List buckets that do not block public access

```sql
select
  name,
  block_public_acls,
  block_public_policy,
  ignore_public_acls,
  restrict_public_buckets
from
  aws_s3_bucket
where
  not block_public_acls
  or not block_public_policy
  or not ignore_public_acls
  or not restrict_public_buckets;
```

### List buckets that block public access through bucket policies

```sql
select
  name,
  bucket_policy_is_public
from
  aws_s3_bucket
where
  bucket_policy_is_public;
```

### List buckets where the server access logging destination is the same as the source bucket

```sql
select
  name,
  logging ->> 'TargetBucket' as target_bucket
from
  aws_s3_bucket
where
  logging ->> 'TargetBucket' = name;
```

### List buckets without the 'application' tags key

```sql
select
  name,
  tags ->> 'fizz' as fizz
from
  aws_s3_bucket
where
  tags ->> 'application' is not null;
```

### List buckets that enforce encryption in transit

```sql
select
  name,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s ->> 'Condition' as conditions,
  ssl
from
  aws_s3_bucket,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a,
  jsonb_array_elements_text(
    s -> 'Condition' -> 'Bool' -> 'aws:securetransport'
  ) as ssl
where
  p = '*'
  and s ->> 'Effect' = 'Deny'
  and ssl :: bool = false;
```

### List buckets that do not enforce encryption in transit

```sql
select
  name
from
  aws_s3_bucket
where
  name not in (
    select
      name
    from
      aws_s3_bucket,
      jsonb_array_elements(policy_std -> 'Statement') as s,
      jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
      jsonb_array_elements_text(s -> 'Action') as a,
      jsonb_array_elements_text(
        s -> 'Condition' -> 'Bool' -> 'aws:securetransport'
      ) as ssl
    where
      p = '*'
      and s ->> 'Effect' = 'Deny'
      and ssl :: bool = false
  );
```

### List bucket policy statements that grant external access for each bucket

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_s3_bucket,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  string_to_array(p, ':') as pa,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and (
    pa[5] != account_id
    or p = '*'
  );
```

### List buckets with object lock enabled

```sql
select
  name,
  object_lock_configuration ->> 'ObjectLockEnabled' as object_lock_enabled
from
  aws_s3_bucket
where
  object_lock_configuration ->> 'ObjectLockEnabled' = 'Enabled';
```

### List buckets with website hosting enabled

```sql
select
  name,
  website_configuration -> 'IndexDocument' ->> 'Suffix' as suffix
from
  aws_s3_bucket
where
  website_configuration -> 'IndexDocument' ->> 'Suffix' is not null;
```

### List object ownership control rules of buckets

```sql
select
  b.name,
  r ->> 'ObjectOwnership' as object_ownership
from
  aws_s3_bucket as b,
  jsonb_array_elements(object_ownership_controls -> 'Rules') as r;
```

### List intelligent tiering configurations of buckets

```sql
select
  name,
  c -> 'Filter' as intelligent_tiering_configuration_filter,
  c ->> 'Id' as intelligent_tiering_configuration_id,
  c ->> 'Status' as intelligent_tiering_configuration_status,
  c -> 'Tierings' as intelligent_tiering_configuration_tierings
from
  aws_s3_bucket,
  jsonb_array_elements(intelligent_tiering_configuration)as c;
```
