---
title: "Steampipe Table: aws_s3_multipart_upload - Query AWS S3 Multipart Uploads using SQL"
description: "Allows users to query AWS S3 Multipart Uploads, providing information about in-progress multipart uploads across S3 buckets."
folder: "S3"
---

# Table: aws_s3_multipart_upload - Query AWS S3 Multipart Uploads using SQL

AWS S3 Multipart Upload is a feature that allows you to upload large objects in parts. This is useful when you're uploading large objects over a network with unreliable connectivity, as it allows you to retry uploading individual parts if any part of the upload fails. Incomplete multipart uploads can incur storage costs, so it's important to monitor and manage them effectively.

## Table Usage Guide

The `aws_s3_multipart_upload` table in Steampipe provides you with information about in-progress multipart uploads within AWS Simple Storage Service (S3). This table enables you, as a DevOps engineer or system administrator, to query upload-specific details, including the upload ID, initiated date, storage class, and associated metadata. You can utilize this table to gather insights on multipart uploads, such as identifying incomplete uploads that might be incurring costs, monitoring upload progress, and managing storage utilization. The schema outlines the various attributes of the S3 multipart upload, including the bucket name, key, upload ID, and associated metadata.

**Important Notes**

- You must specify a `bucket_name` in a where or join clause in order to use this table.
- The table lists all in-progress multipart uploads for the specified bucket.

## Examples

### Basic info

Explore the status of multipart uploads in your S3 buckets to track ongoing file transfers and manage storage costs.

```sql+postgres
select
  bucket_name,
  key,
  upload_id,
  initiated,
  storage_class
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name';
```

```sql+sqlite
select
  bucket_name,
  key,
  upload_id,
  initiated,
  storage_class
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name';
```

### List all multipart uploads older than 7 days

Identify potentially abandoned multipart uploads that might be unnecessarily incurring storage costs.

```sql+postgres
select
  bucket_name,
  key,
  upload_id,
  initiated,
  storage_class
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name'
  and initiated < current_timestamp - interval '7 days';
```

```sql+sqlite
select
  bucket_name,
  key,
  upload_id,
  initiated,
  storage_class
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name'
  and initiated < datetime('now', '-7 days');
```

### Get multipart upload details by bucket and key prefix

Analyze the details of multipart uploads within a specific bucket and key prefix to monitor upload activities and manage storage.

```sql+postgres
select
  bucket_name,
  key,
  upload_id,
  initiated,
  storage_class,
  initiator_id,
  initiator_display_name
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name'
  and key like 'prefix/%';
```

```sql+sqlite
select
  bucket_name,
  key,
  upload_id,
  initiated,
  storage_class,
  initiator_id,
  initiator_display_name
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name'
  and key like 'prefix/%';
```

### List multipart uploads by storage class

Analyze the distribution of multipart uploads across different storage classes to optimize storage costs and performance.

```sql+postgres
select
  storage_class,
  count(*) as upload_count,
  bucket_name
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name'
group by
  storage_class,
  bucket_name
order by
  upload_count desc;
```

```sql+sqlite
select
  storage_class,
  count(*) as upload_count,
  bucket_name
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name'
group by
  storage_class,
  bucket_name
order by
  upload_count desc;
```

### Get multipart uploads with bucket details

Analyze multipart uploads alongside their corresponding bucket information to better understand storage utilization and permissions.

```sql+postgres
select
  u.bucket_name,
  u.key,
  u.upload_id,
  u.initiated,
  u.storage_class,
  b.versioning_enabled,
  b.region
from
  aws_s3_multipart_upload as u
  join aws_s3_bucket as b on u.bucket_name = b.name;
```

```sql+sqlite
select
  u.bucket_name,
  u.key,
  u.upload_id,
  u.initiated,
  u.storage_class,
  b.versioning_enabled,
  b.region
from
  aws_s3_multipart_upload as u
  join aws_s3_bucket as b on u.bucket_name = b.name;
```

### List multipart uploads by initiator

Identify and analyze multipart uploads based on who initiated them, which can be useful for access auditing and resource management.

```sql+postgres
select
  initiator_display_name,
  initiator_id,
  count(*) as upload_count,
  array_agg(distinct bucket_name) as buckets
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name'
group by
  initiator_display_name,
  initiator_id
order by
  upload_count desc;
```

```sql+sqlite
select
  initiator_display_name,
  initiator_id,
  count(*) as upload_count,
  group_concat(distinct bucket_name) as buckets
from
  aws_s3_multipart_upload
where
  bucket_name = 'my-bucket-name'
group by
  initiator_display_name,
  initiator_id
order by
  upload_count desc;
```
