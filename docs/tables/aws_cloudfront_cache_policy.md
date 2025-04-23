---
title: "Steampipe Table: aws_cloudfront_cache_policy - Query AWS CloudFront Cache Policies using SQL"
description: "Allows users to query AWS CloudFront Cache Policies for details about their configuration, status, and associated metadata."
folder: "CloudFront"
---

# Table: aws_cloudfront_cache_policy - Query AWS CloudFront Cache Policies using SQL

The AWS CloudFront Cache Policy is a feature of AWS CloudFront that allows you to specify detailed cache behaviors, including how, when, and where CloudFront caches and delivers content. It provides control over the data that CloudFront uses to serve requests, including headers, cookies, and query strings. This policy aids in optimizing the cache key and improving the cache hit ratio, thereby enhancing the performance of your application.

## Table Usage Guide

The `aws_cloudfront_cache_policy` table in Steampipe provides you with information about Cache Policies within AWS CloudFront. This table allows you, as a DevOps engineer, to query policy-specific details, including the configuration, status, and associated metadata. You can utilize this table to gather insights on cache policies, such as their identifiers, comment descriptions, the default time to live (TTL), maximum and minimum TTL, and more. The schema outlines the various attributes of the cache policy for you, including the policy ARN, creation time, last modified time, and associated tags.

## Examples

### Basic info
Explore which AWS CloudFront cache policies are in place to understand their impact on content delivery and caching strategies. This can be beneficial in optimizing resource usage and reducing costs.

```sql+postgres
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

```sql+sqlite
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
Identify instances where Gzip compression format is not enabled in AWS CloudFront cache policies. This can help to optimize content delivery and improve website loading speeds.

```sql+postgres
select
  id,
  name,
  parameters_in_cache_key_and_forwarded_to_origin ->> 'EnableAcceptEncodingGzip' as enable_gzip
from
  aws_cloudfront_cache_policy
where
  parameters_in_cache_key_and_forwarded_to_origin ->> 'EnableAcceptEncodingGzip' <> 'true';
```

```sql+sqlite
select
  id,
  name,
  json_extract(parameters_in_cache_key_and_forwarded_to_origin, '$.EnableAcceptEncodingGzip') as enable_gzip
from
  aws_cloudfront_cache_policy
where
  json_extract(parameters_in_cache_key_and_forwarded_to_origin, '$.EnableAcceptEncodingGzip') <> 'true';
```

### List cache policies where Brotli compression format is not enabled
Identify instances where Brotli compression format is not enabled in cache policies. This could help improve website performance by enabling more efficient data compression.

```sql+postgres
select
  id,
  name,
  parameters_in_cache_key_and_forwarded_to_origin ->> 'EnableAcceptEncodingBrotli' as enable_brotli
from
  aws_cloudfront_cache_policy
where
  parameters_in_cache_key_and_forwarded_to_origin ->> 'EnableAcceptEncodingBrotli' <> 'true';
```

```sql+sqlite
select
  id,
  name,
  json_extract(parameters_in_cache_key_and_forwarded_to_origin, '$.EnableAcceptEncodingBrotli') as enable_brotli
from
  aws_cloudfront_cache_policy
where
  json_extract(parameters_in_cache_key_and_forwarded_to_origin, '$.EnableAcceptEncodingBrotli') <> 'true';
```