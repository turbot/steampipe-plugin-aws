---
title: "Table: aws_s3_object - Query AWS S3 Object using SQL"
description: "Allows users to query AWS S3 Objects and retrieve metadata and details about each object stored in S3 buckets."
---

# Table: aws_s3_object - Query AWS S3 Object using SQL

The `aws_s3_object` table in Steampipe provides information about objects within AWS Simple Storage Service (S3). This table allows DevOps engineers to query object-specific details, including its size, last modified date, storage class, and associated metadata. Users can utilize this table to gather insights on objects, such as objects' storage utilization, retrieval of object metadata, verification of object encryption status, and more. The schema outlines the various attributes of the S3 object, including the bucket name, key, size, storage class, and associated tags.

You **_must_** specify a `bucket_name` in a where or join clause in order to use this table.

We recommend specifying the `prefix` column when querying buckets with a large number of objects to reduce the query time.

The `body` column returns the raw bytes of the object data as a string. if the bytes entirely consists of valid UTF8 runes, e.g., `.txt files`, an UTF8 data will be set as column value and we will be able to query the object body
([refer example below](#get-data-details-of-a-particular-object-in-a-bucket)) otherwise for the invalid UTF8 runes, e.g., `.png files`, the bas64 encoding of the bytes will be set as column value and we will not be able to query the object body for those objects.

Note: Using this table adds to cost to your monthly bill from AWS. Optimizations have been put in place to minimize the impact as much as possible. Please refer to AWS S3 Pricing to understand the cost implications.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_s3_object` table, you can use the `.inspect aws_s3_object` command in Steampipe.

### Key columns:

- `bucket_name`: The name of the bucket containing the object. This can be used to join with the `aws_s3_bucket` table.
- `key`: The key for the object. This is unique for each object within a bucket and can be used to identify specific objects.
- `storage_class`: The storage class used for the object. This can be used to identify how the object is stored and retrieved, and to manage costs.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

```sql
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

```sql
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

### List all objects in a bucket where any user other than the `OWNER` has `FULL_CONTROL`

```sql
select
  key,
  bucket_name,
  owner,
  acl_grant -> 'Grantee' as grantee,
  acl_grant ->> 'Permission' as permission
from
  aws_s3_object,
  jsonb_array_elements(aws_s3_object.acl -> 'Grants') as acl_grant
where
  bucket_name = 'steampipe-test'
  and acl_grant ->> 'Permission' = 'FULL_CONTROL'
  and acl_grant -> 'Grantee' ->> 'ID' != aws_s3_object.owner ->> 'ID';
```

### List all objects in a bucket where the legal hold is enabled

```sql
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

```sql
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

```sql
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
  and object_lock_retain_until_date > current_date + interval '1 year';
```

### List objects without the 'application' tags key

```sql
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

### List all objects where bucket key is disabled

```sql
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

### List all objects where buckets do not block public access

```sql
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

### Get data details of a particular object in a bucket

```sql
select
  key,
  b ->> 'awsAccountId' as account_id,
  b ->> 'digestEndTime' as digest_end_time,
  b ->> 'digestPublicKeyFingerprint' as digest_public_key_fingerprint,
  b ->> 'digestS3Bucket' as digest_s3_bucket,
  b ->> 'digestStartTime' as digest_start_time
from
  aws_s3_object,
  jsonb_array_elements(body::jsonb) as b
where
  bucket_name = 'steampipe-test'
  and prefix = 'test1/log_text.txt';
```
