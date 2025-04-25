---
title: "Steampipe Table: aws_cloudfront_origin_access_identity - Query AWS CloudFront Origin Access Identity using SQL"
description: "Allows users to query AWS CloudFront Origin Access Identity to fetch detailed information about each identity, including its ID, S3 canonical user ID, caller reference, and associated comment."
folder: "CloudFront"
---

# Table: aws_cloudfront_origin_access_identity - Query AWS CloudFront Origin Access Identity using SQL

The AWS CloudFront Origin Access Identity is a special CloudFront feature that allows secure access to your content within an Amazon S3 bucket. It's used as a virtual identity to enable sharing of your content with CloudFront while restricting access directly to your S3 bucket. Thus, it helps in maintaining the privacy of your data by preventing direct access to S3 resources.

## Table Usage Guide

The `aws_cloudfront_origin_access_identity` table in Steampipe provides you with information about each origin access identity within AWS CloudFront. This table allows you, as a DevOps engineer, to query identity-specific details, including the identity's ID, S3 canonical user ID, caller reference, and associated comment. You can utilize this table to gather insights on origin access identities, such as the identity's configuration and CloudFront caller reference. The schema outlines the various attributes of the origin access identity for you, including the ID, S3 canonical user ID, caller reference, and comment.

## Examples

### Basic Info
Explore the foundational details of your AWS Cloudfront origin access identities to better understand your system's configuration and identify any potential areas for optimization or troubleshooting. This query is particularly useful for gaining insights into the identities' associated comments, user IDs, and unique identifiers, which can assist in system management and auditing tasks.

```sql+postgres
select
  id,
  arn,
  comment,
  s3_canonical_user_id,
  etag
from
  aws_cloudfront_origin_access_identity;
```

```sql+sqlite
select
  id,
  arn,
  comment,
  s3_canonical_user_id,
  etag
from
  aws_cloudfront_origin_access_identity;
```


### List origin access identity with comments
Discover the segments that have comments associated with their origin access identity in AWS Cloudfront. This is useful for understanding which identities have additional information or instructions provided, aiding in better resource management.

```sql+postgres
select
  id,
  arn,
  comment,
  caller_reference
from
  aws_cloudfront_origin_access_identity
where
  comment <> '';
```

```sql+sqlite
select
  id,
  arn,
  comment,
  caller_reference
from
  aws_cloudfront_origin_access_identity
where
  comment != '';
```