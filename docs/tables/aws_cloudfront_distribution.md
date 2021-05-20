# Table: aws_cloudfront_distribution

AWS Systems CloudFront is a web service that speeds up distribution of your static and dynamic web content, such as .html, .css, .js, and image files, to your users.

## Examples

### Basic info

```sql
select
  id,
  arn,
  status,
  domain_name,
  enabled,
  e_tag,
  http_version,
  is_ipv6_enabled
from
  aws_cloudfront_distribution;
```


### List distributions with logging not enabled

```sql
select
  id,
  logging ->> 'Bucket' as bucket,
  logging ->> 'Enabled' as logging_enabled,
  logging ->> 'IncludeCookies' as include_cookies
from
  aws_cloudfront_distribution
where
  logging ->> 'Enabled' = 'false';
```


### List distributions with IPv6 DNS requests not enabled

```sql
select
  id,
  arn,
  status,
  is_ipv6_enabled
from
  aws_cloudfront_distribution
where
  is_ipv6_enabled = 'false';
```


### CloudFront distributions enforce field-level encryption

```sql
select
  id,
  arn,
  default_cache_behavior ->> 'FieldLevelEncryptionId' as field_level_encryption_id,
  default_cache_behavior ->> 'DefaultTTL' as default_ttl
from
  aws_cloudfront_distribution
where
  default_cache_behavior ->> 'FieldLevelEncryptionId' <> '';
```


### The traffic between the CloudFront distributions and their origins is encrypted

```sql
select
  id,
  arn,
  p -> 'CustomOriginConfig' -> 'HTTPPort' as http_port,
  p -> 'CustomOriginConfig' -> 'HTTPSPort' as https_port,
  p -> 'CustomOriginConfig' -> 'OriginKeepaliveTimeout' as origin_keepalive_timeout,
  p -> 'CustomOriginConfig' -> 'OriginProtocolPolicy' as origin_protocol_policy
from
  aws_cloudfront_distribution,
  jsonb_array_elements(origins) as p
where
  p -> 'CustomOriginConfig' ->> 'OriginProtocolPolicy' = 'https-only';
```


### CloudFront distributions origins use insecure SSL protocols

```sql
select
  id,
  arn,
  p -> 'CustomOriginConfig' -> 'OriginSslProtocols' -> 'Items' as items,
  p -> 'CustomOriginConfig' -> 'OriginSslProtocols' -> 'Quantity' as quantity
from
  aws_cloudfront_distribution,
  jsonb_array_elements(origins) as p
where
  p -> 'CustomOriginConfig' -> 'OriginSslProtocols' -> 'Items' ?& array['SSLv3'];
```