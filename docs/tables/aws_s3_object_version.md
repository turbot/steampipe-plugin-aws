```markdown
---
title: "Steampipe Table: aws_s3_object_version - Query AWS S3 Object Versions"
description: "Allows querying information about versions of objects stored in Amazon S3 buckets. This table provides details such as bucket name, delimiter, encoding type, version ID marker, prefix, whether the results are truncated, common prefixes, delete markers, and version information."
---

# Table: aws_s3_object_version - Query AWS S3 Object Versions

The `aws_s3_object_version` table in Steampipe allows you to query information about versions of objects stored in Amazon S3 buckets. This includes details such as the bucket name, delimiter, encoding type, version ID marker, prefix, whether the results are truncated, common prefixes, delete markers, and version information.

## Table Usage Guide

The `aws_s3_object_version` table in Steampipe provides you with information about object versions within AWS Simple Storage Service (S3). This table enables you, as a DevOps engineer, to query object version specific details.

**Important Notes**
- You must specify a `bucket_name` in a where or join clause in order to use this table.
- It's recommended that you specify the `prefix` column when querying buckets with a large number of object versions to reduce the query time.
- Optionally, you can specify the column values `encoding_type`, `delimeter`, or `version_id_marker` in where clause to reduce the query time.

## Examples

### Basic Info

Query basic information about AWS S3 object versions, including the bucket name, delimiter, encoding type, version ID marker, prefix, and whether the results are truncated.

```sql+postgres
select
  bucket_name,
  delimiter,
  encoding_type,
  version_id_marker,
  prefix,
  is_truncated
from
  aws_s3_object_version
where
  bucket_name = 'testbucket';
```

```sql+sqlite
select
  bucket_name,
  delimiter,
  encoding_type,
  version_id_marker,
  prefix,
  is_truncated
from
  aws_s3_object_version
where
  bucket_name = 'testbucket';
```

### List Object Versions with Common Prefixes and Delete Markers

Retrieve object versions along with common prefixes and delete markers.

```sql+postgres
select
  bucket_name,
  delimiter,
  encoding_type,
  version_id_marker,
  prefix,
  is_truncated,
  common_prefixes,
  delete_markers
from
  aws_s3_object_version
where
  bucket_name = 'testbucket';
```

```sql+sqlite
select
  bucket_name,
  delimiter,
  encoding_type,
  version_id_marker,
  prefix,
  is_truncated,
  common_prefixes,
  delete_markers
from
  aws_s3_object_version
where
  bucket_name = 'testbucket';
```

### Get Version Information

Retrieve version information for objects stored in S3 buckets.

```sql+postgres
select
  bucket_name,
  delimiter,
  encoding_type,
  version_id_marker,
  prefix,
  is_truncated,
  version
from
  aws_s3_object_version
where
  bucket_name = 'testbucket';
```

```sql+sqlite
select
  bucket_name,
  delimiter,
  encoding_type,
  version_id_marker,
  prefix,
  is_truncated,
  version
from
  aws_s3_object_version
where
  bucket_name = 'testbucket';
```

# Get the specific version details for an object.

```sql+postgres
 select
  v.bucket_name,
  v.encoding_type,
  v.version_id_marker,
  v.is_truncated,
  v.version ->> 'VersionId' as version_id
from
  aws_s3_object_version as v,
  aws_s3_object as o
where
  v.bucket_name = 'test-delete90'
and
  o.bucket_name = 'test-delete90'
and
  version_id = o.version_id;
```

```sql+sqlite
select
  v.bucket_name,
  v.encoding_type,
  v.version_id_marker,
  v.is_truncated,
  JSON_EXTRACT(v.version, '$.VersionId') as version_id
from
  aws_s3_object_version as v
join
  aws_s3_object AS o on v.bucket_name = o.bucket_name
where
  v.bucket_name = 'test-delete90'
and
  o.bucket_name = 'test-delete90'
and
  JSON_EXTRACT(v.version, '$.VersionId') = o.version_id;
```