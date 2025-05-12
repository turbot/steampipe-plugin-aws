---
title: "Steampipe Table: aws_s3_object - Query AWS S3 Object using SQL"
description: "Allows users to query AWS S3 Objects and retrieve metadata and details about each object stored in S3 buckets."
folder: "S3"
---

# Table: aws_s3_object - Query AWS S3 Object using SQL

The AWS S3 Object is a fundamental entity within the Amazon S3 service. It is a simple storage unit, which users can store and retrieve, ranging from 0 bytes to 5 terabytes of data. It offers robust, secure and scalable storage for data backup, archival and analytics.

## Table Usage Guide

The `aws_s3_object` table in Steampipe provides you with information about objects within AWS Simple Storage Service (S3). This table enables you, as a DevOps engineer, to query object-specific details, including its size, last modified date, storage class, and associated metadata. You can utilize this table to gather insights on objects, such as objects' storage utilization, retrieval of object metadata, verification of object encryption status, and more. The schema outlines the various attributes of the S3 object, including the bucket name, key, size, storage class, and associated tags.

**Important Notes**
- You must specify a `bucket_name` in a where or join clause in order to use this table.
- It's recommended that you specify the `prefix` column when querying buckets with a large number of objects to reduce the query time.
- The `body` column returns the raw bytes of the object data as a string. If the bytes entirely consist of valid UTF8 runes, e.g., `.txt files`, an UTF8 data will be set as column value and you will be able to query the object body ([refer example below](#get-data-details-of-a-particular-object-in-a-bucket)). However, for the invalid UTF8 runes, e.g., `.png files`, the bas64 encoding of the bytes will be set as column value and you will not be able to query the object body for those objects.
- Using this table adds to the cost of your monthly bill from AWS. Optimizations have been put in place to minimize the impact as much as possible. You should refer to AWS S3 Pricing to understand the cost implications.
- If you encrypt an object by using server-side encryption with customer-provided encryption keys (SSE-C) when you store the object in Amazon S3, then when you GET the object, you must use the following query parameter:
  - `sse_customer_algorithm`
  - `sse_customer_key_md5`
  - `sse_customer_key`

## Examples

### Basic info
Explore which items in your specified AWS S3 bucket have been recently modified or have a specific storage class. This is beneficial in managing storage costs and ensuring data integrity.

```sql+postgres
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test';
```

```sql+sqlite
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test';
```

### List all objects with a `prefix` in a bucket
This query is useful for pinpointing specific locations within a storage bucket where objects have a certain prefix. This can assist in organizing and managing data in a more efficient manner by focusing on a specific subset of objects.

```sql+postgres
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and prefix = 'test/logs/2021/03/01/12';
```

```sql+sqlite
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and prefix = 'test/logs/2021/03/01/12';
```

### Get object with a `key` in a bucket
Discover the details of a specific object within a designated bucket, which can be useful for tracking changes, assessing storage class, or identifying versions. This can be particularly beneficial for maintaining data integrity and managing storage within an AWS S3 bucket.

```sql+postgres
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and prefix = 'test/logs/2021/03/01/12/abc.txt';
```

```sql+sqlite
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and prefix = 'test/logs/2021/03/01/12/abc.txt';
```

### List all objects which are encrypted with CMK in a bucket
Explore which objects within a specific S3 bucket have been encrypted using a Customer Managed Key (CMK). This is particularly useful for auditing security measures and ensuring compliance with data protection regulations.

```sql+postgres
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and sse_kms_key_id is not null;
```

```sql+sqlite
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and sse_kms_key_id is not null;
```

### List all objects which were not modified in the last 3 months in a bucket
Determine the areas in your storage where objects have not been updated in the last three months. This can assist in identifying outdated or unused data, optimizing storage use and managing costs.

```sql+postgres
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and last_modified < current_date - interval '3 months';
```

```sql+sqlite
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class,
  version_id
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and last_modified < date('now','-3 months');
```

### List all objects in a bucket where any user other than the `OWNER` has `FULL_CONTROL`
Identify instances where any user, other than the owner, has full control over all objects in a specific bucket. This query is useful for assessing security risks and ensuring proper access management in your AWS S3 buckets.

```sql+postgres
select
  s.key,
  s.bucket_name,
  s.owner,
  acl_grant -> 'Grantee' as grantee,
  acl_grant ->> 'Permission' as permission
from
  aws_s3_object AS s
  cross join lateral jsonb_array_elements(s.acl -> 'Grants') as acl_grant
where
  s.bucket_name = 'steampipe-test'
  and acl_grant ->> 'Permission' = 'FULL_CONTROL'
  and acl_grant -> 'Grantee' ->> 'ID' != s.owner ->> 'ID';
```

```sql+sqlite
select
  s.key,
  s.bucket_name,
  s.owner,
  json_extract(acl_grant, '$.Grantee') as grantee,
  json_extract(acl_grant, '$.Permission') as permission
from
  aws_s3_object as s,
  json_each(json_extract(aws_s3_object.acl, '$.Grants')) as acl_grant
where
  bucket_name = 'steampipe-test'
  and json_extract(acl_grant, '$.Permission') = 'FULL_CONTROL'
  and json_extract(json_extract(acl_grant, '$.Grantee'), '$.ID') != json_extract(json_extract(aws_s3_object.owner, '$.ID'));
```

### List all objects in a bucket where the legal hold is enabled
Identify instances where legal holds are active within a specific storage bucket. This can be useful for maintaining compliance and managing data retention policies.

```sql+postgres
select
  key,
  bucket_name,
  object_lock_legal_hold_status
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and object_lock_legal_hold_status = 'ON';
```

```sql+sqlite
select
  key,
  bucket_name,
  object_lock_legal_hold_status
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and object_lock_legal_hold_status = 'ON';
```

### List all objects in a bucket with governance lock mode enabled
Discover the segments that have a governance lock mode enabled within a specific bucket. This is useful for maintaining compliance and ensuring data immutability in regulated industries.

```sql+postgres
select
  key,
  bucket_name,
  object_lock_retain_until_date,
  object_lock_mode,
  object_lock_legal_hold_status
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and object_lock_mode = 'GOVERNANCE';
```

```sql+sqlite
select
  key,
  bucket_name,
  object_lock_retain_until_date,
  object_lock_mode,
  object_lock_legal_hold_status
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and object_lock_mode = 'GOVERNANCE';
```

### List all objects in a bucket which are set to be retained for more than 1 year from now
Discover the objects within a specific storage area that are scheduled to be kept for over a year from the current date. This query can be useful in understanding long-term data retention policies and identifying potential areas for storage optimization.

```sql+postgres
select
  s.key,
  s.bucket_name,
  s.object_lock_retain_until_date,
  s.object_lock_mode,
  s.object_lock_legal_hold_status
from
  aws_s3_object as s
where
  bucket_name = 'steampipe-test'
  and object_lock_retain_until_date > current_date + interval '1 year';
```

```sql+sqlite
select
  key,
  bucket_name,
  object_lock_retain_until_date,
  object_lock_mode,
  object_lock_legal_hold_status
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and date(object_lock_retain_until_date) > date('now','+1 year');
```

### List objects without the 'application' tags key
Discover the segments that have the 'application' tag key in the 'steampipe-test' bucket. This can be particularly useful when trying to identify specific objects within a large bucket for management or organization purposes.

```sql+postgres
select
  key,
  bucket_name,
  jsonb_pretty(tags) as tags
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and tags ->> 'application' is not null;
```

```sql+sqlite
select
  key,
  bucket_name,
  tags
from
  aws_s3_object
where
  bucket_name = 'steampipe-test'
  and json_extract(tags, '$.application') is not null;
```

### List all objects where bucket key is disabled
Determine the areas in which the bucket key is disabled in your AWS S3 objects and buckets. This information can be useful for identifying potential security risks or configuration issues.

```sql+postgres
select
  key,
  o.arn as object_arn,
  bucket_name,
  last_modified,
  bucket_key_enabled
from
  aws_s3_object as o,
  aws_s3_bucket as b
where
  o.bucket_name = b.name
  and not bucket_key_enabled;
```

```sql+sqlite
select
  key,
  o.arn as object_arn,
  bucket_name,
  last_modified,
  bucket_key_enabled
from
  aws_s3_object as o
join
  aws_s3_bucket as b
on
  o.bucket_name = b.name
where
  not bucket_key_enabled;
```

### List all objects where buckets do not block public access
Determine the areas in which your stored objects are potentially exposed to public access due to their hosting buckets' security settings. This is useful for identifying potential security risks and ensuring proper access controls are in place.

```sql+postgres
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class
from
  aws_s3_object
where
  bucket_name in
  (
    select
      name
    from
      aws_s3_bucket
    where
      not block_public_acls
      or not block_public_policy
      or not ignore_public_acls
      or not restrict_public_buckets
  );
```

```sql+sqlite
select
  key,
  arn,
  bucket_name,
  last_modified,
  storage_class
from
  aws_s3_object
where
  bucket_name in
  (
    select
      name
    from
      aws_s3_bucket
    where
      block_public_acls = 0
      or block_public_policy = 0
      or ignore_public_acls = 0
      or restrict_public_buckets = 0
  );
```

### Get data details of a particular object in a bucket
Discover the segments that provide specific information about an object within a particular storage area. This can be useful for monitoring changes over time or verifying the integrity of the object.

```sql+postgres
select
  s.key,
  b ->> 'awsAccountId' as account_id,
  b ->> 'digestEndTime' as digest_end_time,
  b ->> 'digestPublicKeyFingerprint' as digest_public_key_fingerprint,
  b ->> 'digestS3Bucket' as digest_s3_bucket,
  b ->> 'digestStartTime' as digest_start_time
from
  aws_s3_object as s,
  jsonb_array_elements(body::jsonb) as b
where
  bucket_name = 'steampipe-test'
  and prefix = 'test1/log_text.txt';
```

```sql+sqlite
select
  s.key,
  json_extract(b.value, '$.awsAccountId') as account_id,
  json_extract(b.value, '$.digestEndTime') as digest_end_time,
  json_extract(b.value, '$.digestPublicKeyFingerprint') as digest_public_key_fingerprint,
  json_extract(b.value, '$.digestS3Bucket') as digest_s3_bucket,
  json_extract(b.value, '$.digestStartTime') as digest_start_time
from
  aws_s3_object as s,
  json_each(json(body)) as b
where
  bucket_name = 'steampipe-test'
  and prefix = 'test1/log_text.txt';
```

### Retrieve object details encrypted with customer-provided encryption keys (SSE-C)
This query retrieves details of objects stored in an S3 bucket that are encrypted using Server-Side Encryption with Customer-Provided Keys (SSE-C). SSE-C allows you to provide your own encryption keys when storing objects in S3, ensuring that AWS does not manage the encryption keys.

```sql+postgres
select 
  key,
  bucket_name,
  content_language,
  content_length,
  content_type
from
  aws_s3_object 
where 
  bucket_name = 'tes-encryption-31' 
  and sse_customer_algorithm = 'AES256' 
  and sse_customer_key = '/J03dxHHdcPTDNi97Aq7mYxBjnxOX0kV6UzSHVOh8es=' 
  and sse_customer_key_md5 = 'gaWCs7+kcAeTCCLlbVdTXA==';
```

```sql+sqlite
select 
  key,
  bucket_name,
  content_language,
  content_length,
  content_type
from
  aws_s3_object 
where 
  bucket_name = 'tes-encryption-31' 
  and sse_customer_algorithm = 'AES256' 
  and sse_customer_key = '/J03dxHHdcPTDNi97Aq7mYxBjnxOX0kV6UzSHVOh8es=' 
  and sse_customer_key_md5 = 'gaWCs7+kcAeTCCLlbVdTXA==';
```