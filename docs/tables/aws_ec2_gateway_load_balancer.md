---
title: "Steampipe Table: aws_ec2_gateway_load_balancer - Query AWS EC2 Gateway Load Balancer using SQL"
description: "Allows users to query AWS EC2 Gateway Load Balancer details, including its configuration, state, type, and associated tags."
folder: "EC2"
---

# Table: aws_ec2_gateway_load_balancer - Query AWS EC2 Gateway Load Balancer using SQL

The AWS EC2 Gateway Load Balancer is a resource that operates at the third layer of the Open Systems Interconnection (OSI) model, the network layer. It is designed to manage, scale, and secure your network traffic in a simple and cost-effective manner. This service provides you with a single point of contact for all network traffic, regardless of the scale, and ensures that it is efficiently distributed across multiple resources.

## Table Usage Guide

The `aws_ec2_gateway_load_balancer` table in Steampipe provides you with information about Gateway Load Balancers within Amazon Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query load balancer-specific details, including its configuration, state, type, and associated tags. You can utilize this table to gather insights on load balancers, such as their availability zones, subnets, and security groups. The schema outlines the various attributes of the Gateway Load Balancer for you, including the load balancer ARN, creation date, DNS name, scheme, and associated tags.

## Examples

### Basic gateway load balancer info
Determine the areas in which your AWS EC2 gateway load balancer is deployed and its current operational state. This information can help you assess the elements within your network infrastructure and optimize for better performance.

```sql+postgres
select
  name,
  arn,
  type,
  state_code,
  vpc_id,
  availability_zones
from
  aws_ec2_gateway_load_balancer;

```

```sql+sqlite
select
  name,
  arn,
  type,
  state_code,
  vpc_id,
  availability_zones
from
  aws_ec2_gateway_load_balancer;
```

### Availability zone information of all the gateway load balancers
Determine the areas in which your gateway load balancers are located and gain insights into their specific settings. This can help you assess your load balancing strategy and optimize resource allocation.

```sql+postgres
select
  name,
  az ->> 'LoadBalancerAddresses' as load_balancer_addresses,
  az ->> 'OutpostId' as outpost_id,
  az ->> 'SubnetId' as subnet_id,
  az ->> 'ZoneName' as zone_name
from
  aws_ec2_gateway_load_balancer,
  jsonb_array_elements(availability_zones) as az;

```

```sql+sqlite
select
  name,
  json_extract(az.value, '$.LoadBalancerAddresses') as load_balancer_addresses,
  json_extract(az.value, '$.OutpostId') as outpost_id,
  json_extract(az.value, '$.SubnetId') as subnet_id,
  json_extract(az.value, '$.ZoneName') as zone_name
from
  aws_ec2_gateway_load_balancer,
  json_each(availability_zones) as az;
```

### List of gateway load balancers whose availability zone count is less than 2
Determine the areas in which gateway load balancers may be at risk of service disruption due to having less than two availability zones. This can help in proactive infrastructure planning and risk mitigation.

```sql+postgres
select
  name,
  count(az ->> 'ZoneName') as zone_count
from
  aws_ec2_gateway_load_balancer,
  jsonb_array_elements(availability_zones) as az
group by
  name
having
  count(az ->> 'ZoneName') < 2;
```

```sql+sqlite
select
  name,
  count(json_extract(az.value, '$.ZoneName')) as zone_count
from
  aws_ec2_gateway_load_balancer,
  json_each(availability_zones) as az
group by
  name
having
  count(json_extract(az.value, '$.ZoneName')) < 2;
```

### List of gateway load balancers whose deletion protection is not enabled
Identify instances where gateway load balancers do not have deletion protection enabled. This can be useful to ensure the security and longevity of your data by avoiding accidental deletion.

```sql+postgres
select
  name,
  lb ->> 'Key' as deletion_protection_key,
  lb ->> 'Value' as deletion_protection_value
from
  aws_ec2_gateway_load_balancer,
  jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'deletion_protection.enabled'
  and lb ->> 'Value' = 'false';
```

```sql+sqlite
select
  name,
  json_extract(lb.value, '$.Key') as deletion_protection_key,
  json_extract(lb.value, '$.Value') as deletion_protection_value
from
  aws_ec2_gateway_load_balancer,
  json_each(load_balancer_attributes) as lb
where
  json_extract(lb.value, '$.Key') = 'deletion_protection.enabled'
  and json_extract(lb.value, '$.Value') = 'false';
```

### List of gateway load balancers whose load balancing cross zone is enabled
Explore which gateway load balancers have the cross-zone load balancing feature enabled. This is useful in understanding the traffic distribution across multiple zones for better load balancing and increased application availability.

```sql+postgres
select
  name,
  lb ->> 'Key' as load_balancing_cross_zone_key,
  lb ->> 'Value' as load_balancing_cross_zone_value
from
  aws_ec2_gateway_load_balancer,
  jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'load_balancing.cross_zone.enabled'
  and lb ->> 'Value' = 'true';
```

```sql+sqlite
select
  name,
  json_extract(lb.value, '$.Key') as load_balancing_cross_zone_key,
  json_extract(lb.value, '$.Value') as load_balancing_cross_zone_value
from
  aws_ec2_gateway_load_balancer,
  json_each(load_balancer_attributes) as lb
where
  json_extract(lb.value, '$.Key') = 'load_balancing.cross_zone.enabled'
  and json_extract(lb.value, '$.Value') = 'true';
```

### Security group attached to the gateway load balancers
Identify instances where your security groups are linked to your gateway load balancers. This can help you assess your security setup and ensure appropriate measures are in place.

```sql+postgres
select
  name,
  jsonb_array_elements_text(security_groups) as attached_security_group
from
  aws_ec2_gateway_load_balancer;
```

```sql+sqlite
select
  name,
  json_extract(json_each.value, '$') as attached_security_group
from
  aws_ec2_gateway_load_balancer,
  json_each(security_groups);
```

### List of gateway load balancer with state other than active
Identify instances where gateway load balancers in AWS EC2 are not in an 'active' state. This is useful to pinpoint potential issues or disruptions in network traffic routing.

```sql+postgres
select
  name,
  state_code
from
  aws_ec2_gateway_load_balancer
where
 state_code <> 'active';
```

```sql+sqlite
select
  name,
  state_code
from
  aws_ec2_gateway_load_balancer
where
 state_code != 'active';
```