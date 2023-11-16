---
title: "Table: aws_route53_zone - Query AWS Route 53 Zone using SQL"
description: "Allows users to query AWS Route 53 Zone details including hosted zone ID, name, type, record set count, and associated tags."
---

# Table: aws_route53_zone - Query AWS Route 53 Zone using SQL

The `aws_route53_zone` table in Steampipe provides information about hosted zones within AWS Route 53. This table allows DevOps engineers to query zone-specific details, including the hosted zone ID, name, type, record set count, and associated tags. Users can utilize this table to gather insights on hosted zones, such as the number of record sets within each zone, the type of zone (public or private), and more. The schema outlines the various attributes of the hosted zone, including the zone ID, name, type, record set count, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_zone` table, you can use the `.inspect aws_route53_zone` command in Steampipe.

### Key columns:

- `id`: The hosted zone ID. This unique identifier can be used to join this table with other Route 53 tables.
- `name`: The name of the hosted zone. This can be useful for users who want to query specific zones by name.
- `type`: The type of the hosted zone (public or private). This can be useful for users who want to filter zones based on their type.

## Examples

### Basic Zone Info
```sql
select
  name,
  id,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone;
```

### List private zones  
```sql
select
  name,
  id,
  comment,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone
where
  private_zone;
```

### List public zones  
```sql
select
  name,
  id,
  comment,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone
where
  not private_zone;
```

### Find zones by subdomain name

```sql
select
  name,
  id,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone
where
  name like '%.turbot.com.'
```

### List VPCs associated with zones

```sql
select 
  name,
  id,
  v ->> 'VPCId' as vpc_id,
  v ->> 'VPCRegion' as vpc_region
from
  aws_route53_zone,
  jsonb_array_elements(vpcs) as v;
```

### Get VPC details associated with zones

```sql
select 
  name,
  id,
  v.vpc_id as vpc_id,
  v.cidr_block as cidr_block,
  v.is_default as is_default,
  v.dhcp_options_id as dhcp_options_id
from
  aws_route53_zone,
  jsonb_array_elements(vpcs) as p,
  aws_vpc as v
where
  p ->> 'VPCId' = v.vpc_id;
```