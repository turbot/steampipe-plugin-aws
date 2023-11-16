---
title: "Table: aws_s3_multi_region_access_point - Query AWS S3 Multi-Region Access Point using SQL"
description: "Allows users to query AWS S3 Multi-Region Access Points to retrieve information about their configuration, status, and associated policies."
---

# Table: aws_s3_multi_region_access_point - Query AWS S3 Multi-Region Access Point using SQL

The `aws_s3_multi_region_access_point` table in Steampipe provides information about Multi-Region Access Points within Amazon Simple Storage Service (S3). This table allows DevOps engineers to query Multi-Region Access Point-specific details, including the name, ARN, status, creation time, and associated policies. Users can utilize this table to gather insights on Multi-Region Access Points, such as their current status, the buckets they are associated with, and the policies applied to them. The schema outlines the various attributes of the Multi-Region Access Point, including the ARN, alias, home region, and associated bucket details.

Amazon S3 Multi-Region Access Point provides a global endpoint that applications can use to fulfill requests from S3 buckets located in multiple AWS Regions. You can use Multi-Region Access Points to build multi-region applications with the same architecture that's used in a single region, and then run those applications anywhere in the world. Instead of sending requests over the congested public internet, Multi-Region Access Points provide built-in network resilience with the acceleration of internet-based requests to Amazon S3.

You must grant the s3:ListAllMyBuckets permission to the user, role, or an IAM entity that makes a request to manage a Multi-Region Access Point.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_s3_multi_region_access_point` table, you can use the `.inspect aws_s3_multi_region_access_point` command in Steampipe.

### Key columns:

- `name`: The name of the Multi-Region Access Point. This can be used to join this table with other tables that contain information about specific Multi-Region Access Points.
- `arn`: The Amazon Resource Name (ARN) of the Multi-Region Access Point. This can be used to join this table with other tables that contain ARN-specific information.
- `home_region`: The home region of the Multi-Region Access Point. This can be useful when joining this table with other tables that contain region-specific information.

## Examples

### Basic info

```sql
select
  alias,
  status,
  created_at
from
  aws_s3_multi_region_access_point;
```

### List multi-region access points that do not block public access

```sql
select
  name,
  public_access_block ->> 'BlockPublicAcls' as block_public_acls,
  public_access_block ->> 'BlockPublicPolicy' as block_public_policy,
  public_access_block ->> 'IgnorePublicAcls' as ignore_public_acls,
  public_access_block ->> 'RestrictPublicBuckets' as restrict_public_buckets 
from
  aws_s3_multi_region_access_point 
where
  public_access_block ->> 'BlockPublicAcls'::text = 'true' 
  and public_access_block ->> 'BlockPublicPolicy'::text = 'true' 
  and public_access_block ->> 'IgnorePublicAcls'::text = 'true' 
  and public_access_block ->> 'RestrictPublicBuckets'::text = 'true';
```

### Get policy details of each multi-region access point

```sql
select
  name,
  policy -> 'Established' -> 'Policy' as established_policy,
  policy -> 'Proposed' -> 'Policy' as proposed_policy
from
  aws_s3_multi_region_access_point;
```

### Count the number of multi-region access points per bucket

```sql
select
  r ->> 'Bucket' as bucket_name,
  count(name) access_point_count
from
  aws_s3_multi_region_access_point,
  jsonb_array_elements(regions) as r
group by
  bucket_name;
```

## Get bucket details of each multi-region access point

```sql
select
  a.name,
  r ->> 'Bucket' as bucket_name,
  b.creation_date as bucket_creation_date,
  b.bucket_policy_is_public,
  b.versioning_enabled
from
  aws_s3_multi_region_access_point as a,
  jsonb_array_elements(a.regions) as r,
  aws_s3_bucket as b
where
  b.name = r ->> 'Bucket';
```
