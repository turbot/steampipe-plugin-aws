# Table: aws_s3_object

Amazon S3 objects are stored in one or more Amazon S3 buckets, and each object can be up to 5 TB in size.

You **_must_** specify a `bucket_name` in a where or join clause in order to use this table.

We recommend specifying the `prefix` and `key` columns when querying buckets with a large number of objects to reduce the query time.

> Note: Using this table adds to cost to your monthly bill from AWS. Optimizations have been put in place to minimize the impact as much as possible. Please refer to AWS S3 Pricing to understand the cost implications.

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

### List all objects with a fixed `prefix` in a bucket

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

### Get object with a fixed `key` in a bucket

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
  and key = 'test/logs/2021/03/01/12/abc.txt';
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

### List all objects in a bucket where legal hold is on

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

### List all objects in a bucket with governance lock mode

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
