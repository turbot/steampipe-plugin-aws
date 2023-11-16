---
title: "Table: aws_wafv2_ip_set - Query AWS WAFv2 IPSet using SQL"
description: "Allows users to query AWS WAFv2 IPSet information, including IP addresses, IP address version, and associated metadata."
---

# Table: aws_wafv2_ip_set - Query AWS WAFv2 IPSet using SQL

The `aws_wafv2_ip_set` table in Steampipe provides information about IPSet within AWS WAFv2. This table allows DevOps engineers to query IPSet-specific details, including IP addresses, IP address version, and associated metadata. Users can utilize this table to gather insights on IPSet, such as the IP addresses that AWS WAF is inspecting for web requests, the IP address version (IPv4 or IPv6), and more. The schema outlines the various attributes of the IPSet, including the IPSet ID, IPSet name, IPSet ARN, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wafv2_ip_set` table, you can use the `.inspect aws_wafv2_ip_set` command in Steampipe.

**Key columns**:

- `id`: The unique identifier for the IPSet. This can be used to join this table with other tables.
- `name`: The name of the IPSet. This can be useful for querying specific IPSet.
- `arn`: The Amazon Resource Name (ARN) of the IPSet. This can be used to join this table with other tables where ARN is a common attribute.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

```sql
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
