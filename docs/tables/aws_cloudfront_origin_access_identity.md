---
title: "Table: aws_cloudfront_origin_access_identity - Query AWS CloudFront Origin Access Identity using SQL"
description: "Allows users to query AWS CloudFront Origin Access Identity to fetch detailed information about each identity, including its ID, S3 canonical user ID, caller reference, and associated comment."
---

# Table: aws_cloudfront_origin_access_identity - Query AWS CloudFront Origin Access Identity using SQL

The `aws_cloudfront_origin_access_identity` table in Steampipe provides information about each origin access identity within AWS CloudFront. This table allows DevOps engineers to query identity-specific details, including the identity's ID, S3 canonical user ID, caller reference, and associated comment. Users can utilize this table to gather insights on origin access identities, such as the identity's configuration and CloudFront caller reference. The schema outlines the various attributes of the origin access identity, including the ID, S3 canonical user ID, caller reference, and comment.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudfront_origin_access_identity` table, you can use the `.inspect aws_cloudfront_origin_access_identity` command in Steampipe.

**Key columns**:

- `id`: The ID for the origin access identity. This is a unique identifier that can be used to join this table with other tables.
- `s3_canonical_user_id`: The Amazon S3 canonical user ID for the origin access identity, which is used in an S3 bucket policy for CloudFront to access an S3 bucket. This can be useful for tracking access and permissions.
- `caller_reference`: The caller reference for the origin access identity. This is a value that you provide when you create the identity, and it can be useful for identifying the origin access identity when joining with other tables.

## Examples

### Basic Info

```sql
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

```sql
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
