---
title: "Table: aws_globalaccelerator_accelerator - Query AWS Global Accelerator using SQL"
description: "Allows users to query AWS Global Accelerator's accelerators."
---

# Table: aws_globalaccelerator_accelerator - Query AWS Global Accelerator using SQL

The `aws_globalaccelerator_accelerator` table in Steampipe provides information about accelerators within AWS Global Accelerator. Accelerators direct traffic to optimal endpoints over the AWS global network to improve the availability and performance of your applications. This table allows DevOps engineers to query accelerator-specific details, including the accelerator's ARN, creation time, and status. Users can utilize this table to gather insights on accelerators, such as the accelerators' DNS names, IP address sets, and associated tags. The schema outlines the various attributes of the accelerator, including the accelerator ARN, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_globalaccelerator_accelerator` table, you can use the `.inspect aws_globalaccelerator_accelerator` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Name (ARN) of the accelerator. This can be used to join this table with other tables.
- `name`: The name of the accelerator. This is a unique identifier that can be used to join this table with other tables.
- `status`: The status of the accelerator. This can provide insights into the operational status of the accelerators.

## Examples

### Basic info

```sql
select
  name,
  created_time,
  dns_name,
  enabled,
  ip_address_type,
  last_modified_time,
  status
from
  aws_globalaccelerator_accelerator;
```

### List IPs used by global accelerators

```sql
 select
   name,
   created_time,
   dns_name,
   enabled,
   ip_address_type,
   last_modified_time,
   status,
   anycast_ip
from
  aws_globalaccelerator_accelerator,
  jsonb_array_elements(ip_sets -> 0 -> 'IpAddresses') as anycast_ip;
```

### List global accelerators without owner tag key

```sql
select
  name,
  tags
from
  aws_globalaccelerator_accelerator
where
  not tags::JSONB ? 'owner';
```
