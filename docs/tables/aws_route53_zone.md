---
title: "Steampipe Table: aws_route53_zone - Query AWS Route 53 Zone using SQL"
description: "Allows users to query AWS Route 53 Zone details including hosted zone ID, name, type, record set count, and associated tags."
folder: "Route 53"
---

# Table: aws_route53_zone - Query AWS Route 53 Zone using SQL

The AWS Route 53 Zone is a component of Amazon's scalable Domain Name System (DNS) web service. It is designed to provide highly reliable and cost-effective domain registration, DNS routing, and health checking of resources within your environment. Route 53 effectively connects user requests to infrastructure running in AWS, such as Amazon EC2 instances, Elastic Load Balancing load balancers, or Amazon S3 buckets, and can also be used to route users to infrastructure outside of AWS.

## Table Usage Guide

The `aws_route53_zone` table in Steampipe provides you with information about hosted zones within AWS Route 53. This table allows you, as a DevOps engineer, to query zone-specific details, including the hosted zone ID, name, type, record set count, and associated tags. You can utilize this table to gather insights on hosted zones, such as the number of record sets within each zone, the type of zone (public or private), and more. The schema outlines the various attributes of the hosted zone for you, including the zone ID, name, type, record set count, and associated tags.

## Examples

### Basic Zone Info
Explore which zones in your AWS Route53 service are private and assess the number of resource records within each. This information can be useful in managing DNS configurations and understanding the distribution of resources within your network.
```sql+postgres
select
  name,
  id,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone;
```

```sql+sqlite
select
  name,
  id,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone;
```

### List private zones  
Discover the segments that are designated as private within the AWS Route53 service. This is particularly useful when you need to manage or review the privacy settings of your DNS zones.
```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  comment,
  private_zone,
  resource_record_set_count
from 
  aws_route53_zone
where
  private_zone = 1;
```

### List public zones  
Explore which DNS zones are public within your AWS Route53 service. This can be useful to understand your public-facing infrastructure and ensure appropriate security measures are in place.

```sql+postgres
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

```sql+sqlite
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
Explore which zones are linked to a specific subdomain to gain insights into their privacy settings and the volume of resource records they contain. This is particularly useful for managing and understanding the distribution of resources within a domain.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which virtual private clouds (VPCs) are associated with specific zones. This can be useful for understanding the geographical distribution of your VPCs and their connections to different zones.

```sql+postgres
select 
  name,
  id,
  v ->> 'VPCId' as vpc_id,
  v ->> 'VPCRegion' as vpc_region
from
  aws_route53_zone,
  jsonb_array_elements(vpcs) as v;
```

```sql+sqlite
select 
  name,
  id,
  json_extract(v.value, '$.VPCId') as vpc_id,
  json_extract(v.value, '$.VPCRegion') as vpc_region
from
  aws_route53_zone,
  json_each(vpcs) as v;
```

### Get VPC details associated with zones
Explore which Virtual Private Clouds (VPCs) are associated with specific zones in AWS Route53. This can be useful in understanding the network architecture and identifying potential security risks or configuration issues.

```sql+postgres
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

```sql+sqlite
select 
  name,
  id,
  v.vpc_id as vpc_id,
  v.cidr_block as cidr_block,
  v.is_default as is_default,
  v.dhcp_options_id as dhcp_options_id
from
  aws_route53_zone,
  aws_vpc as v
where
  json_extract(vpcs, '$.VPCId') = v.vpc_id;
```