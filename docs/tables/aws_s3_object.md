# Table: aws_s3_object

Amazon S3 objects are stored in one or more Amazon S3 buckets, and each object can be up to 5 TB in size.

To list objects, you must mention the `bucket` which contains the objects.

## Examples

### Basic info

```sql
select
  key,
  etag,
  size
from
  aws_s3_object
where
  bucket = 'cloudflare_logs_2021_03_01';
```

### List all objects with a fixed `prefix`

```sql
select
  key,
  etag,
  size,
  prefix
from
  aws_s3_object
where
  bucket = 'cloudflare_logs_2021_03_01'
  and prefix = '/logs/2021/03/01/12';
```

### List all objects with a fixed `key`

```sql
select
  key,
  etag,
  size,
  prefix
from
  aws_s3_object
where
  bucket = 'cloudflare_logs_2021_03_01'
  and key = '/logs/2021/03/01/12/05/32.log';
```

### List all objects which were not modified in the last 3 months

```sql
select
  key,
  bucket,
  last_modified,
  etag,
  size
from
  aws_s3_object
where
  bucket = 'static_assets'
  and last_modified < current_date - interval '3 months';
```

### List all objects in a bucket where any user other than the `OWNER` has `FULL_CONTROL`

```sql
select
  aws_s3_object.key,
  aws_s3_object.bucket,
  aws_s3_object.acl -> 'Owner' as owner,
  acl_grant -> 'Grantee' as grantee,
  acl_grant ->> 'Permission' as permission
from
  aws_s3_object,
  jsonb_array_elements(aws_s3_object.acl -> 'Grants') as acl_grant
where
  bucket = 'sensitive_assets'
  and acl_grant ->> 'Permission' = 'FULL_CONTROL'
  and acl_grant -> 'Grantee' ->> 'ID' != aws_s3_object.acl -> 'Owner' ->> 'ID';
```

### List all objects in a bucket `legal_hold` is set

```sql
select
  aws_s3_object.key,
  aws_s3_object.bucket,
  aws_s3_object.legal_hold
from
  aws_s3_object
where
  bucket = 'sensitive_assets'
  and legal_hold is not null;
```

### List all objects and their lock information

```sql
select
  key,
  bucket,
  to_timestamp(retention ->> 'RetainUntilDate', 'YYYY-MM-DDTHH:MI:SS.FF6TZH') as retain_until,
  retention ->> 'Mode' as retention_mode,
  legal_hold 
from
  aws_s3_object 
where
  bucket = 'static_assets';
```

### List all objects in a bucket which are set to be retained for more than 1 year from now

```sql
select
  key,
  bucket,
  last_modified,
  etag,
  size,
  retention ->> 'RetainUntilDate' as retain_until
from
  aws_s3_object
where
  bucket = 'static_assets'
  and to_timestamp(retention ->> 'RetainUntilDate', 'YYYY-MM-DDTHH:MI:SS.FF6TZH') > current_date + interval '1 year';
```

### List objects without the 'application' tags key

```sql
select
  key,
  bucket,
  tags ->> 'fizz' as fizz
from
  aws_s3_object
where
  bucket = 'static_assets'
  and tags ->> 'application' is not null;
```
