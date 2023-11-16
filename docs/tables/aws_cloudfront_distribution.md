---
title: "Table: aws_cloudfront_distribution - Query AWS CloudFront Distributions using SQL"
description: "Allows users to query AWS CloudFront Distributions to gain insights into their configuration, status, and associated metadata."
---

# Table: aws_cloudfront_distribution - Query AWS CloudFront Distributions using SQL

The `aws_cloudfront_distribution` table in Steampipe provides information about distributions within AWS CloudFront. This table allows DevOps engineers to query distribution-specific details, including distribution configuration, status, and associated metadata. Users can utilize this table to gather insights on distributions, such as viewing all distributions, checking if logging is enabled, verifying if a distribution is configured to use a custom SSL certificate, and more. The schema outlines the various attributes of the CloudFront distribution, including the ARN, domain name, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudfront_distribution` table, you can use the `.inspect aws_cloudfront_distribution` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Name (ARN) of the distribution. This is a unique identifier that can be used to join this table with other AWS tables.
- `domain_name`: The domain name of the distribution. This can be useful for identifying specific distributions or for linking with DNS records in other tables.
- `status`: The status of the distribution (e.g., Deployed, InProgress). This can be useful for tracking the deployment status of distributions.

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


### List distributions with logging disabled

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


### List distributions that enforce field-level encryption

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


### List distributions whose origins use encrypted traffic

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


### List distributions whose origins use insecure SSL protocols

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
