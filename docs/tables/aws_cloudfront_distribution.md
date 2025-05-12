---
title: "Steampipe Table: aws_cloudfront_distribution - Query AWS CloudFront Distributions using SQL"
description: "Allows users to query AWS CloudFront Distributions to gain insights into their configuration, status, and associated metadata."
folder: "CloudFront"
---

# Table: aws_cloudfront_distribution - Query AWS CloudFront Distributions using SQL

The AWS CloudFront Distributions is a part of Amazon's content delivery network (CDN) services. It speeds up the distribution of your static and dynamic web content, such as .html, .css, .js, and image files, to your users. CloudFront delivers your content through a worldwide network of data centers called edge locations and ensures that end-user requests are served by the closest edge location.

## Table Usage Guide

The `aws_cloudfront_distribution` table in Steampipe provides you with information about distributions within AWS CloudFront. This table allows you, as a DevOps engineer, to query distribution-specific details, including distribution configuration, status, and associated metadata. You can utilize this table to gather insights on distributions, such as viewing all distributions, checking if logging is enabled, verifying if a distribution is configured to use a custom SSL certificate, and more. The schema outlines the various attributes of the CloudFront distribution for you, including the ARN, domain name, status, and associated tags.

## Examples

### Basic info
Analyze the settings of your AWS Cloudfront distributions to understand their current status and configuration. This can help you to identify potential issues or areas for improvement, such as outdated HTTP versions or disabled IPv6.

```sql+postgres
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

```sql+sqlite
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

### List distributions with logging disabled
Determine the areas in your AWS Cloudfront distribution settings where logging is disabled. This is useful for identifying potential gaps in your logging strategy, which could impact security and troubleshooting capabilities.

```sql+postgres
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

```sql+sqlite
select
  id,
  json_extract(logging, '$.Bucket') as bucket,
  json_extract(logging, '$.Enabled') as logging_enabled,
  json_extract(logging, '$.IncludeCookies') as include_cookies
from
  aws_cloudfront_distribution
where
  json_extract(logging, '$.Enabled') = 'false';
```


### List distributions with IPv6 DNS requests not enabled
Identify instances where IPv6 DNS requests are not enabled within your AWS CloudFront distributions. This can help in improving network performance and future-proofing your system as IPv6 becomes more prevalent.

```sql+postgres
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

```sql+sqlite
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

### List distributions that enforce field-level encryption
Determine the areas in which field-level encryption is enforced within your distributions. This can be handy for improving security by ensuring sensitive data fields are encrypted.

```sql+postgres
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

```sql+sqlite
select
  id,
  arn,
  json_extract(default_cache_behavior, '$.FieldLevelEncryptionId') as field_level_encryption_id,
  json_extract(default_cache_behavior, '$.DefaultTTL') as default_ttl
from
  aws_cloudfront_distribution
where
  json_extract(default_cache_behavior, '$.FieldLevelEncryptionId') <> '';
```


### List distributions whose origins use encrypted traffic
Determine the areas in which your AWS Cloudfront distributions are utilizing encrypted traffic. This can be beneficial to ensure data security and compliance with industry standards and regulations.

```sql+postgres
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

```sql+sqlite
select
  'id',
  arn,
  json_extract(p.value, '$.CustomOriginConfig.HTTPPort') as http_port,
  json_extract(p.value, '$.CustomOriginConfig.HTTPSPort') as https_port,
  json_extract(p.value, '$.CustomOriginConfig.OriginKeepaliveTimeout') as origin_keepalive_timeout,
  json_extract(p.value, '$.CustomOriginConfig.OriginProtocolPolicy') as origin_protocol_policy
from
  aws_cloudfront_distribution,
  json_each(origins) as p
where
  json_extract(p.value, '$.CustomOriginConfig.OriginProtocolPolicy') = 'https-only';
```


### List distributions whose origins use insecure SSL protocols
Discover the segments of your Cloudfront distributions where origins are using insecure SSL protocols. This is useful for identifying potential security vulnerabilities in your network.

```sql+postgres
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

```sql+sqlite
select
  'id',
  arn,
  json_extract(p.value, '$.CustomOriginConfig.OriginSslProtocols.Items') as items,
  json_extract(p.value, '$.CustomOriginConfig.OriginSslProtocols.Quantity') as quantity
from
  aws_cloudfront_distribution,
  json_each(origins) as p
where
  json_extract(p.value, '$.CustomOriginConfig.OriginSslProtocols.Items') LIKE '%SSLv3%';
```