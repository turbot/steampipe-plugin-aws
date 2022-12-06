# Table: aws_s3_multi_region_access_point

Amazon S3 Multi-Region Access Points provide a global endpoint that applications can use to fulfill requests from S3 buckets located in multiple AWS Regions. You can use Multi-Region Access Points to build multi-Region applications with the same architecture that's used in a single region, and then run those applications anywhere in the world. Instead of sending requests over the congested public internet, Multi-Region Access Points provide built-in network resilience with acceleration of internet-based requests to Amazon S3.

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

### List multi region access points that do not block public access

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
and
  public_access_block ->> 'BlockPublicPolicy':: text = 'true'
and
  public_access_block ->> 'IgnorePublicAcls':: text = 'true'
and
  public_access_block ->> 'RestrictPublicBuckets':: text = 'true';
```

### Get policy details for each multi region access point

```sql
select
  name,
  policy -> 'Established' -> 'Policy' as established_policy,
  policy -> 'Proposed' -> 'Policy' as proposed_policy
from
  aws_s3_multi_region_access_point;
```

### Count the number of multi region access points per bucket

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

## Get bucket details for each multi region access point

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