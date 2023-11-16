---
title: "Table: aws_cloudfront_cache_policy - Query AWS CloudFront Cache Policies using SQL"
description: "Allows users to query AWS CloudFront Cache Policies for details about their configuration, status, and associated metadata."
---

# Table: aws_cloudfront_cache_policy - Query AWS CloudFront Cache Policies using SQL

The `aws_cloudfront_cache_policy` table in Steampipe provides information about Cache Policies within AWS CloudFront. This table allows DevOps engineers to query policy-specific details, including the configuration, status, and associated metadata. Users can utilize this table to gather insights on cache policies, such as their identifiers, comment descriptions, the default time to live (TTL), maximum and minimum TTL, and more. The schema outlines the various attributes of the cache policy, including the policy ARN, creation time, last modified time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudfront_cache_policy` table, you can use the `.inspect aws_cloudfront_cache_policy` command in Steampipe.

Key columns:

- `id`: The unique identifier of the cache policy. This can be used to join this table with others that reference cache policies by their ID.
- `arn`: The Amazon Resource Number of the cache policy. It is a globally unique identifier that is useful for joining with other tables that reference cache policies by their ARN.
- `last_modified_time`: The time when the cache policy was last modified. This is useful for tracking changes to cache policies over time.

## Examples

### Basic info

```sql
select
  id,
  name,
  comment,
  min_ttl,
  etag,
  last_modified_time
from
  aws_cloudfront_cache_policy;
```

### List cache policies where Gzip compression format is not enabled

```sql
select
  id,
  name,
  parameters_in_cache_key_and_forwarded_to_origin ->> 'EnableAcceptEncodingGzip' as enable_gzip
from
  aws_cloudfront_cache_policy
where
  parameters_in_cache_key_and_forwarded_to_origin ->> 'EnableAcceptEncodingGzip' <> 'true';
```

### List cache policies where Brotli compression format is not enabled

```sql
select
  id,
  name,
  parameters_in_cache_key_and_forwarded_to_origin ->> 'EnableAcceptEncodingBrotli' as enable_brotli
from
  aws_cloudfront_cache_policy
where
  parameters_in_cache_key_and_forwarded_to_origin ->> 'EnableAcceptEncodingBrotli' <> 'true';
```
