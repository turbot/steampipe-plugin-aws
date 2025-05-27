---
title: "Steampipe Table: aws_s3_object_version - Query AWS S3 Object Versions"
description: "Allows querying information about versions of objects stored in Amazon S3 buckets. This table provides details such as bucket name, delimiter, encoding type, version ID marker, prefix, whether the results are truncated, common prefixes, delete markers, and version information."
folder: "S3"
---

# Table: aws_s3_object_version - Query AWS S3 Object Versions

The `aws_s3_object_version` table in Steampipe allows you to query information about versions of objects stored in Amazon S3 buckets. This includes details such as the bucket name, delimiter, encoding type, version ID marker, prefix, whether the results are truncated, common prefixes, delete markers, and version information.

## Table Usage Guide

The `aws_s3_object_version` table in Steampipe provides you with information about object versions within AWS Simple Storage Service (S3). This table enables you, as a DevOps engineer, to query object version specific details.

**Important Notes**

- You must specify a `bucket_name` in a where or join clause in order to use this table.
- It's recommended that you specify the `prefix` column when querying buckets with a large number of object versions to reduce the query time.

## Examples

### Basic Info

Query basic information about AWS S3 object versions, including the bucket name, size, version ID, storage class and last update.

```sql+postgres
select
  bucket_name,
  key,
  storage_class,
  version_id,
  is_latest,
  size
from
  aws_s3_object_version
where
  bucket_name = 'testbucket';
```

```sql+sqlite
select
  bucket_name,
  key,
  storage_class,
  version_id,
  is_latest,
  size
from
  aws_s3_object_version
where
  bucket_name = 'testbucket';
```

### List Object Versions for a particular object

Retrieve object versions along with common prefixes and delete markers.

```sql+postgres
select
  bucket_name,
  key,
  storage_class,
  version_id,
  is_latest,
  size,
  etag,
  owner_id
from
  aws_s3_object_version
where
  bucket_name = 'testbucket'
and
  key = 'test/template.txt';
```

```sql+sqlite
select
  bucket_name,
  key,
  storage_class,
  version_id,
  is_latest,
  size,
  etag,
  owner_id
from
  aws_s3_object_version
where
  bucket_name = 'testbucket'
and
  key = 'test/template.txt';
```

# Get the specific version details of objects
Ensure that you specify the exact version identifier for each object of interest. This process typically involves accessing a version-controlled storage system or database where each object can have multiple versions, each distinguished by a unique version ID.

```sql+postgres
 select
  v.bucket_name,
  v.key,
  v.storage_class,
  v.version_id,
  o.accept_ranges,
  o.body,
  o.content_type
from
  aws_s3_object_version as v,
  aws_s3_object as o
where
  v.bucket_name = 'test-delete90'
and
  o.bucket_name = 'test-delete90'
and
  v.version_id = o.version_id;
```

```sql+sqlite
select
  v.bucket_name,
  v.key,
  v.storage_class,
  v.version_id,
  o.accept_ranges,
  o.body,
  o.content_type
from
  aws_s3_object_version as v
join
  aws_s3_object as o on v.bucket_name = o.bucket_name
where
  v.bucket_name = 'test-delete90'
and
  o.bucket_name = 'test-delete90'
and
  v.version_id = o.version_id;
```
