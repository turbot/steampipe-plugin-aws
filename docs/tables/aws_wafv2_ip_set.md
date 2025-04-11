---
title: "Steampipe Table: aws_wafv2_ip_set - Query AWS WAFv2 IPSet using SQL"
description: "Allows users to query AWS WAFv2 IPSet information, including IP addresses, IP address version, and associated metadata."
folder: "WAFv2"
---

# Table: aws_wafv2_ip_set - Query AWS WAFv2 IPSet using SQL

The AWS WAFv2 IPSet is a feature of AWS Web Application Firewall (WAF) service. It allows you to specify lists of IP addresses that you want to allow or block based on the originating IP addresses of a web request. This helps in protecting your web applications from common web exploits that could affect application availability, compromise security, or consume excessive resources.

## Table Usage Guide

The `aws_wafv2_ip_set` table in Steampipe provides you with information about IPSet within AWS WAFv2. This table allows you, as a DevOps engineer, to query IPSet-specific details, including IP addresses, IP address version, and associated metadata. You can utilize this table to gather insights on IPSet, such as the IP addresses that AWS WAF is inspecting for web requests, the IP address version (IPv4 or IPv6), and more. The schema outlines the various attributes of the IPSet for you, including the IPSet ID, IPSet name, IPSet ARN, and associated tags.

## Examples

### Basic info
Explore the basic details of IP sets in your AWS WAFv2 to understand their scope, location, and associated IP addresses. This can be useful in identifying potential security vulnerabilities or areas for improvement in your network protection strategy.

```sql+postgres
select
  name,
  description,
  arn,
  id,
  scope,
  addresses,
  ip_address_version,
  region
from
  aws_wafv2_ip_set;
```

```sql+sqlite
select
  name,
  description,
  arn,
  id,
  scope,
  addresses,
  ip_address_version,
  region
from
  aws_wafv2_ip_set;
```

### List global (CLOUDFRONT) IP sets
Discover the segments that are using the global scope in AWS Web Application Firewall version 2. This is useful for understanding your network's configuration and identifying potential areas for security enhancements.

```sql+postgres
select
  name,
  description,
  arn,
  id,
  scope,
  addresses,
  ip_address_version,
  region
from
  aws_wafv2_ip_set
where
  scope = 'CLOUDFRONT';
```

```sql+sqlite
select
  name,
  description,
  arn,
  id,
  scope,
  addresses,
  ip_address_version,
  region
from
  aws_wafv2_ip_set
where
  scope = 'CLOUDFRONT';
```


### List IP sets with an IPv4 address version
Determine the areas in which IP sets are utilizing the IPv4 address version within your AWS WAFv2 configuration. This can be useful for understanding your network's current IP version usage and planning for potential upgrades or changes.

```sql+postgres
select
  name,
  description,
  arn,
  id,
  scope,
  addresses,
  ip_address_version,
  region
from
  aws_wafv2_ip_set
where
  ip_address_version = 'IPV4';
```

```sql+sqlite
select
  name,
  description,
  arn,
  id,
  scope,
  addresses,
  ip_address_version,
  region
from
  aws_wafv2_ip_set
where
  ip_address_version = 'IPV4';
```


### List IP sets having a specific IP address
Determine the areas in which a specific IP address is included in your IP sets. This can be useful for identifying potential security risks or for troubleshooting network issues.

```sql+postgres
select
  name,
  description,
  arn,
  ip_address_version,
  region,
  address
from
  aws_wafv2_ip_set,
  jsonb_array_elements_text(addresses) as address
where
  address = '1.2.3.4/32';
```

```sql+sqlite
select
  name,
  description,
  arn,
  ip_address_version,
  region,
  json_extract(address.value, '$') as address
from
  aws_wafv2_ip_set,
  json_each(addresses) as address
where
  json_extract(address.value, '$') = '1.2.3.4/32';
```