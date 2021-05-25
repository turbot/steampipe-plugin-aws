# Table: aws_cloudfront_cache_policy

CloudFront provides a set of managed cache policies that you can attach to any of your distribution’s cache behaviors. With a managed cache policy, you don’t need to write or maintain your own cache policy. The managed policies use settings that are optimized for specific use cases. However you can write your own cache policy also if required.

## Examples

### Basic info

```sql
select
  id,
  name,
  comment,
  min_ttl,
  e_tag,
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
