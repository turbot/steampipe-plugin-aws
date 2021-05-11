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
  logging ->> 'Enabled' as logging_enable,
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