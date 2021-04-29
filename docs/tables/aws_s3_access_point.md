# Table: aws_s3_access_point

An Amazon S3 bucket is a public cloud storage resource available in Amazon Web Services' (AWS) Simple Storage Service (S3), an object storage offering.

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
