---
title: "Table: aws_s3_access_point - Query AWS S3 Access Point using SQL"
description: "Allows users to query AWS S3 Access Point details such as name, bucket, network origin, policy status, creation time, and more."
---

# Table: aws_s3_access_point - Query AWS S3 Access Point using SQL

The `aws_s3_access_point` table in Steampipe provides information about Access Points within AWS Simple Storage Service (S3). This table allows DevOps engineers, developers, and data analysts to query Access Point-specific details, including the Access Point's name, associated bucket, network origin, policy status, and creation time. Users can utilize this table to gather insights on Access Points, such as their permissions, associated buckets, and more. The schema outlines the various attributes of the S3 Access Point, including the ARN, bucket name, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_s3_access_point` table, you can use the `.inspect aws_s3_access_point` command in Steampipe.

**Key columns**:

- `name`: The name of the access point. This column can be used to join with other tables that need access point information.
- `bucket`: The name of the bucket associated with the access point. This column is useful for joining with other tables that need bucket information.
- `arn`: The Amazon Resource Name (ARN) of the access point. This column is useful for joining with other tables that need ARN information.

## Examples

### Basic info

```sql
select
  name,
  access_point_arn,
  bucket_name
from
  aws_s3_access_point;
```


### List access points that only accept requests from a VPC

```sql
select
  name,
  access_point_arn,
  vpc_id
from
  aws_s3_access_point
where
  vpc_id is not null;
```


### List access points that do not block public access

```sql
select
  name,
  block_public_acls,
  block_public_policy,
  ignore_public_acls,
  restrict_public_buckets
from
  aws_s3_access_point
where
  not block_public_acls
  or not block_public_policy
  or not ignore_public_acls
  or not restrict_public_buckets;
```


### List buckets that allows public access through their policies

```sql
select
  name,
  access_point_policy_is_public
from
  aws_s3_access_point
where
  access_point_policy_is_public;
```


### Count the number of access points per bucket

```sql
select
  bucket_name,
  count(name) access_point_count
from
  aws_s3_access_point
group by
  bucket_name;
```
