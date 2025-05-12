---
title: "Steampipe Table: aws_s3_access_point - Query AWS S3 Access Point using SQL"
description: "Allows users to query AWS S3 Access Point details such as name, bucket, network origin, policy status, creation time, and more."
folder: "S3"
---

# Table: aws_s3_access_point - Query AWS S3 Access Point using SQL

The AWS S3 Access Point is a feature of the AWS S3 service that simplifies managing data access at scale for applications using shared data sets on S3. Access Points are unique hostnames with dedicated access policies that describe how data can be accessed using that endpoint. They offer a way to easily manage access to shared datasets by creating separate access points for specific users or roles.

## Table Usage Guide

The `aws_s3_access_point` table in Steampipe provides you with information about Access Points within AWS Simple Storage Service (S3). This table enables you, as a DevOps engineer, developer, or data analyst, to query Access Point-specific details, including the Access Point's name, associated bucket, network origin, policy status, and creation time. You can utilize this table to gather insights on Access Points, such as their permissions, associated buckets, and more. The schema outlines the various attributes of the S3 Access Point for you, including the ARN, bucket name, creation date, and associated tags.

## Examples

### Basic info
Discover the segments that have been granted access to your S3 buckets. This query is useful in identifying and managing the access points to your AWS S3 resources, thereby enhancing your data security.

```sql+postgres
select
  name,
  access_point_arn,
  bucket_name
from
  aws_s3_access_point;
```

```sql+sqlite
select
  name,
  access_point_arn,
  bucket_name
from
  aws_s3_access_point;
```


### List access points that only accept requests from a VPC
Discover the segments that are restricted to only accept requests from a Virtual Private Cloud (VPC), allowing for increased security and control over your AWS S3 access points. This is particularly useful for organizations that want to limit their access points to specific network resources.

```sql+postgres
select
  name,
  access_point_arn,
  vpc_id
from
  aws_s3_access_point
where
  vpc_id is not null;
```

```sql+sqlite
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
Determine the areas in which your AWS S3 access points may be allowing public access. This is useful for identifying potential security vulnerabilities and ensuring your data is adequately protected.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which public access is permitted through policy settings. This query is useful for identifying potential security risks and ensuring proper data protection measures are in place.

```sql+postgres
select
  name,
  access_point_policy_is_public
from
  aws_s3_access_point
where
  access_point_policy_is_public;
```

```sql+sqlite
select
  name,
  access_point_policy_is_public
from
  aws_s3_access_point
where
  access_point_policy_is_public = 1;
```


### Count the number of access points per bucket
Discover the segments that are using various access points by counting them per storage bucket. This is beneficial in managing resources and understanding usage patterns.

```sql+postgres
select
  bucket_name,
  count(name) access_point_count
from
  aws_s3_access_point
group by
  bucket_name;
```

```sql+sqlite
select
  bucket_name,
  count(name) as access_point_count
from
  aws_s3_access_point
group by
  bucket_name;
```