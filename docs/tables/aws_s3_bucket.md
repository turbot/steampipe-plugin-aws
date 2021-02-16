# Table: aws_s3_bucket

An Amazon S3 bucket is a public cloud storage resource available in Amazon Web Services' (AWS) Simple Storage Service (S3), an object storage offering.

## Examples

### List of buckets where versioning is not enabled

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


### List of buckets which do not have default encryption enabled

```sql
select
  name,
  server_side_encryption_configuration
from
  aws_s3_bucket
where
  server_side_encryption_configuration is null;
```


### List of buckets where public access to S3 bucket is not blocked

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


### List of buckets which have public access enabled in bucket policies

```sql
select
  name,
  bucket_policy_is_public
from
  aws_s3_bucket
where
  bucket_policy_is_public;
```


### List of buckets where server access logging destination is same as the source bucket

```sql
select
  name,
  logging ->> 'TargetBucket' as target_bucket
from
  aws_s3_bucket
where
  logging ->> 'TargetBucket' = name;
```


### List of buckets where tags key 'fizz' is not null

```sql
select
  name,
  tags ->> 'fizz' as fizz
from
  aws_s3_bucket
where
  tags ->> 'fizz' is not null;
```

### List of buckets that enforce encryption in transit

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


### List of buckets that DO NOT enforce encryption in transit

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


### List of Bucket policy statements that grant external access
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