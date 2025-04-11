---
title: "Steampipe Table: aws_s3_bucket - Query AWS S3 Buckets using SQL"
description: "Allows users to query AWS S3 buckets for detailed information about their configuration, policies, and permissions."
folder: "S3"
---

# Table: aws_s3_bucket - Query AWS S3 Buckets using SQL

An AWS S3 Bucket is a public cloud storage resource available in Amazon Web Services' (AWS) Simple Storage Service (S3). It is used to store objects, which consist of data and its descriptive metadata. S3 makes it possible to store and retrieve varying amounts of data, at any time, from anywhere on the web.

## Table Usage Guide

The `aws_s3_bucket` table in Steampipe provides you with information about S3 buckets within Amazon Simple Storage Service (S3). This table allows you, as a DevOps engineer, to query bucket-specific details, including configuration, policies, and permissions. You can utilize this table to gather insights on buckets, such as bucket policies, access controls, versioning status, and more. The schema outlines for you the various attributes of the S3 bucket, including the bucket name, creation date, region, and associated tags.

## Examples

### Basic info
Explore which AWS S3 buckets are set as public in different regions to enhance your data security by identifying potential vulnerabilities. This allows you to manage your AWS resources more effectively by pinpointing specific locations where public access may be a concern.

```sql+postgres
select
  name,
  region,
  account_id,
  bucket_policy_is_public
from
  aws_s3_bucket;
```

```sql+sqlite
select
  name,
  region,
  account_id,
  bucket_policy_is_public
from
  aws_s3_bucket;
```

### List buckets with versioning disabled
Discover the segments that have versioning disabled in your Amazon S3 buckets. This could be useful in identifying potential risks or compliance issues related to data version control.

```sql+postgres
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

```sql+sqlite
select
  name,
  region,
  account_id,
  versioning_enabled
from
  aws_s3_bucket
where
  versioning_enabled is not 1;
```

### List buckets with default encryption disabled
Uncover the details of S3 buckets that lack default encryption, a potential security risk, to enhance data protection measures in your AWS environment.

```sql+postgres
select
  name,
  server_side_encryption_configuration
from
  aws_s3_bucket
where
  server_side_encryption_configuration is null;
```

```sql+sqlite
select
  name,
  server_side_encryption_configuration
from
  aws_s3_bucket
where
  server_side_encryption_configuration is null;
```

### List buckets that do not block public access
Identify instances where AWS S3 buckets may be vulnerable due to not blocking public access. This query is useful for assessing potential security risks associated with unrestricted public access to your data.

```sql+postgres
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

```sql+sqlite
select
  name,
  block_public_acls,
  block_public_policy,
  ignore_public_acls,
  restrict_public_buckets
from
  aws_s3_bucket
where
  block_public_acls = 0
  or block_public_policy = 0
  or ignore_public_acls = 0
  or restrict_public_buckets = 0;
```

### List buckets that block public access through bucket policies
Identify instances where certain storage buckets have implemented measures to block public access, enhancing data security and privacy. This could be useful in ensuring compliance with privacy regulations and preventing unauthorized data access.

```sql+postgres
select
  name,
  bucket_policy_is_public
from
  aws_s3_bucket
where
  bucket_policy_is_public;
```

```sql+sqlite
select
  name,
  bucket_policy_is_public
from
  aws_s3_bucket
where
  bucket_policy_is_public = 1;
```

### List buckets where the server access logging destination is the same as the source bucket
Identify instances where the destination for server access logging is the same as the source bucket in AWS S3. This can help in understanding potential security risks or misconfigurations in your logging setup.

```sql+postgres
select
  name,
  logging ->> 'TargetBucket' as target_bucket
from
  aws_s3_bucket
where
  logging ->> 'TargetBucket' = name;
```

```sql+sqlite
select
  name,
  json_extract(logging, '$.TargetBucket') as target_bucket
from
  aws_s3_bucket
where
  json_extract(logging, '$.TargetBucket') = name;
```

### List buckets without the 'application' tags key
Discover the buckets that have not been tagged specifically for application purposes. This is particularly useful for managing and organizing your buckets based on their intended application usage.

```sql+postgres
select
  name,
  tags ->> 'fizz' as fizz
from
  aws_s3_bucket
where
  tags ->> 'application' is null;
```

```sql+sqlite
select
  name,
  json_extract(tags, '$.fizz') as fizz
from
  aws_s3_bucket
where
  json_extract(tags, '$.application') is null;
```

### List buckets that enforce encryption in transit
Determine the areas in which AWS S3 buckets have encryption in transit enforced. This query is useful in identifying potential security vulnerabilities, as buckets without this feature may be at risk of data interception during transfer.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(p.value, '$') as principal,
  json_extract(a.value, '$') as action,
  json_extract(s.value, '$.Effect') as effect,
  json_extract(s.value, '$.Condition') as conditions,
  json_extract(ssl.value, '$') as ssl
from
  aws_s3_bucket,
  json_each(json_extract(policy_std, '$.Statement')) as s,
  json_each(json_extract(s.value, '$.Principal.AWS')) as p,
  json_each(json_extract(s.value, '$.Action')) as a,
  json_each(json_extract(s.value, '$.Condition.Bool."aws:securetransport"')) as ssl
where
  json_extract(p.value, '$') = '*'
  and json_extract(s.value, '$.Effect') = 'Deny'
  and json_extract(ssl.value, '$') = 'false';
```

### List buckets that do not enforce encryption in transit
Determine the areas in your AWS S3 service where encryption in transit is not enforced. This is useful for identifying potential security risks and ensuring that your data is always protected during transmission.

```sql+postgres
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

```sql+sqlite
select
  name
from
  aws_s3_bucket
where
  name not in (
    select
      aws_s3_bucket.name
    from
      aws_s3_bucket,
      json_each(json_extract(policy_std, '$.Statement')) as s,
      json_each(json_extract(s.value, '$.Principal.AWS')) as p,
      json_each(json_extract(s.value, '$.Action')) as a,
      json_each(json_extract(s.value, '$.Condition.Bool."aws:securetransport"')) as ssl
    where
      json_extract(p.value, '$') = '*'
      and json_extract(s.value, '$.Effect') = 'Deny'
      and json_extract(ssl.value, '$') = 'false'
  );
```

### List bucket policy statements that grant external access for each bucket
Determine the areas in which your S3 bucket policies may be granting external access. This is useful for identifying potential security risks and ensuring only authorized access to your buckets.

```sql+postgres
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

```sql+sqlite
Error: SQLite does not support string_to_array functions.
```

### List buckets with object lock enabled
Determine the areas in which AWS S3 buckets have the object lock feature enabled. This is useful for understanding where additional data protection measures are in place.

```sql+postgres
select
  name,
  object_lock_configuration ->> 'ObjectLockEnabled' as object_lock_enabled
from
  aws_s3_bucket
where
  object_lock_configuration ->> 'ObjectLockEnabled' = 'Enabled';
```

```sql+sqlite
select
  name,
  json_extract(object_lock_configuration, '$.ObjectLockEnabled') as object_lock_enabled
from
  aws_s3_bucket
where
  json_extract(object_lock_configuration, '$.ObjectLockEnabled') = 'Enabled';
```

### List buckets with website hosting enabled
Discover the segments that have website hosting enabled within your AWS S3 buckets. This can be useful in identifying where your web content is stored or determining which buckets are serving as websites.

```sql+postgres
select
  name,
  website_configuration -> 'IndexDocument' ->> 'Suffix' as suffix
from
  aws_s3_bucket
where
  website_configuration -> 'IndexDocument' ->> 'Suffix' is not null;
```

```sql+sqlite
select
  name,
  json_extract(website_configuration, '$.IndexDocument.Suffix') as suffix
from
  aws_s3_bucket
where
  json_extract(website_configuration, '$.IndexDocument.Suffix') is not null;
```

### List object ownership control rules of buckets
Explore which AWS S3 buckets have specific object ownership control rules. This can be useful in managing access permissions and ensuring appropriate data governance.

```sql+postgres
select
  b.name,
  r ->> 'ObjectOwnership' as object_ownership
from
  aws_s3_bucket as b,
  jsonb_array_elements(object_ownership_controls -> 'Rules') as r;
```

```sql+sqlite
select
  b.name,
  json_extract(r.value, '$.ObjectOwnership') as object_ownership
from
  aws_s3_bucket as b,
  json_each(b.object_ownership_controls, '$.Rules') as r;
```