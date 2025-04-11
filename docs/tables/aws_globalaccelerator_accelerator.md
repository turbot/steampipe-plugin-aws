---
title: "Steampipe Table: aws_globalaccelerator_accelerator - Query AWS Global Accelerator using SQL"
description: "Allows users to query AWS Global Accelerator's accelerators."
folder: "Global Accelerator"
---

# Table: aws_globalaccelerator_accelerator - Query AWS Global Accelerator using SQL

The AWS Global Accelerator is a networking service that improves the availability and performance of the applications you offer to your global users. It leverages the vast, congestion-free AWS global network to direct internet traffic from your users to your applications on AWS. With Global Accelerator, your users are directed to your workload based on their geographic location, application health, and routing policies that you configure.

## Table Usage Guide

The `aws_globalaccelerator_accelerator` table in Steampipe provides you with information about accelerators within AWS Global Accelerator. These accelerators direct traffic to optimal endpoints over the AWS global network to enhance the availability and performance of your applications. This table allows you, as a DevOps engineer, to query accelerator-specific details, including the accelerator's ARN, creation time, and status. You can utilize this table to gather insights on accelerators, such as the accelerators' DNS names, IP address sets, and associated tags. The schema outlines the various attributes of the accelerator for you, including the accelerator ARN, creation date, and associated tags.

## Examples

### Basic info
Explore the status and key details of your AWS Global Accelerator to understand its operational state and configuration. This can assist in troubleshooting or optimizing performance.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which global accelerators are being used by identifying their IP addresses. This can help in assessing network performance and identifying potential issues.

```sql+postgres
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

```sql+sqlite
select
  name,
  created_time,
  dns_name,
  enabled,
  ip_address_type,
  last_modified_time,
  status,
  json_extract(ip.value, '$') as anycast_ip
from
  aws_globalaccelerator_accelerator,
  json_each(json_extract(ip_sets, '$[0].IpAddresses')) as ip;
```

### List global accelerators without owner tag key
Identify instances where global accelerators in AWS do not have an assigned owner. This is useful for maintaining organizational oversight and ensuring accountability for resource usage.

```sql+postgres
select
  name,
  tags
from
  aws_globalaccelerator_accelerator
where
  not tags::JSONB ? 'owner';
```

```sql+sqlite
select
  name,
  tags
from
  aws_globalaccelerator_accelerator
where
  json_extract(tags, '$.owner') is null;
```