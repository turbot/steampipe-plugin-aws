---
title: "Steampipe Table: aws_s3_multi_region_access_point - Query AWS S3 Multi-Region Access Point using SQL"
description: "Allows users to query AWS S3 Multi-Region Access Points to retrieve information about their configuration, status, and associated policies."
folder: "Region"
---

# Table: aws_s3_multi_region_access_point - Query AWS S3 Multi-Region Access Point using SQL

The AWS S3 Multi-Region Access Point is a feature of AWS S3 that simplifies data access across multiple regions. It enhances performance by providing a single global endpoint to access a data set that is replicated across multiple geographies. It also offers automatic routing of requests to the bucket in the region that delivers the lowest latency.

## Table Usage Guide

The `aws_s3_multi_region_access_point` table in Steampipe provides you with information about Multi-Region Access Points within Amazon Simple Storage Service (S3). This table allows you, as a DevOps engineer, to query Multi-Region Access Point-specific details, including the name, ARN, status, creation time, and associated policies. You can utilize this table to gather insights on Multi-Region Access Points, such as their current status, the buckets they are associated with, and the policies applied to them. The schema outlines the various attributes of the Multi-Region Access Point for you, including the ARN, alias, home region, and associated bucket details.

Amazon S3 Multi-Region Access Point provides you with a global endpoint that your applications can use to fulfill requests from S3 buckets located in multiple AWS Regions. You can use Multi-Region Access Points to build multi-region applications with the same architecture that's used in a single region, and then run those applications anywhere in the world. Instead of sending requests over the congested public internet, Multi-Region Access Points provide you with built-in network resilience with the acceleration of internet-based requests to Amazon S3.

**Important Notes**
- You must grant the s3:ListAllMyBuckets permission to yourself, your role, or an IAM entity that makes a request to manage a Multi-Region Access Point.

## Examples

### Basic info
Explore which multi-region access points in AWS S3 are active or inactive and when they were created. This can be useful for auditing your AWS S3 configuration and ensuring that only necessary access points are active.

```sql+postgres
select
  alias,
  status,
  created_at
from
  aws_s3_multi_region_access_point;
```

```sql+sqlite
select
  alias,
  status,
  created_at
from
  aws_s3_multi_region_access_point;
```

### List multi-region access points that do not block public access
Discover the segments of multi-region access points in your AWS S3 storage that do not have public access restrictions. This can help you identify potential security risks and ensure appropriate access controls are in place.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(public_access_block, '$.BlockPublicAcls') as block_public_acls,
  json_extract(public_access_block, '$.BlockPublicPolicy') as block_public_policy,
  json_extract(public_access_block, '$.IgnorePublicAcls') as ignore_public_acls,
  json_extract(public_access_block, '$.RestrictPublicBuckets') as restrict_public_buckets 
from
  aws_s3_multi_region_access_point 
where
  json_extract(public_access_block, '$.BlockPublicAcls') = 'true' 
  and json_extract(public_access_block, '$.BlockPublicPolicy') = 'true' 
  and json_extract(public_access_block, '$.IgnorePublicAcls') = 'true' 
  and json_extract(public_access_block, '$.RestrictPublicBuckets') = 'true';
```

### Count the number of multi-region access points per bucket
Explore the distribution of multi-region access points across different buckets to better understand your AWS S3 usage patterns

```sql+postgres
select
  r ->> 'Bucket' as bucket_name,
  count(name) access_point_count
from
  aws_s3_multi_region_access_point,
  jsonb_array_elements(regions) as r
group by
  bucket_name;
```

```sql+sqlite
select
  json_extract(r.value, '$.Bucket') as bucket_name,
  count(name) as access_point_count
from
  aws_s3_multi_region_access_point,
  json_each(regions) as r
group by
  bucket_name;
```

### Get bucket details of each multi-region access point
Gain insights into the details of each multi-region access point, including bucket creation date and versioning status, to enhance your understanding of your AWS S3 configuration.
```sql+postgres
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

```sql+sqlite
select
  a.name,
  json_extract(r.value, '$.Bucket') as bucket_name,
  b.creation_date as bucket_creation_date,
  b.bucket_policy_is_public,
  b.versioning_enabled
from
  aws_s3_multi_region_access_point as a,
  json_each(a.regions) as r,
  aws_s3_bucket as b
where
  b.name = json_extract(r.value, '$.Bucket');
```