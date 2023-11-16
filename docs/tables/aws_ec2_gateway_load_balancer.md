---
title: "Table: aws_ec2_gateway_load_balancer - Query AWS EC2 Gateway Load Balancer using SQL"
description: "Allows users to query AWS EC2 Gateway Load Balancer details, including its configuration, state, type, and associated tags."
---

# Table: aws_ec2_gateway_load_balancer - Query AWS EC2 Gateway Load Balancer using SQL

The `aws_ec2_gateway_load_balancer` table in Steampipe provides information about Gateway Load Balancers within Amazon Elastic Compute Cloud (EC2). This table allows DevOps engineers to query load balancer-specific details, including its configuration, state, type, and associated tags. Users can utilize this table to gather insights on load balancers, such as their availability zones, subnets, and security groups. The schema outlines the various attributes of the Gateway Load Balancer, including the load balancer ARN, creation date, DNS name, scheme, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_gateway_load_balancer` table, you can use the `.inspect aws_ec2_gateway_load_balancer` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the load balancer. This can be used to join with other tables that contain load balancer ARNs.
- `load_balancer_name`: The name of the load balancer. This is useful for joining with tables that contain load balancer names.
- `vpc_id`: The ID of the Virtual Private Cloud (VPC) for the load balancer. This is useful for joining with tables that contain VPC IDs.

## Examples

### Basic gateway load balancer info

```sql
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

```sql
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

### List of gateway load balancers whose availability zone count is less than 2

```sql
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

### List of gateway load balancers whose deletion protection is not enabled

```sql
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

### List of gateway load balancers whose load balancing cross zone is enabled

```sql
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

### Security group attached to the gateway load balancers

```sql
select
  name,
  jsonb_array_elements_text(security_groups) as attached_security_group
from
  aws_ec2_gateway_load_balancer;
```

### List of gateway load balancer with state other than active

```sql
select
  name,
  state_code
from
  aws_ec2_gateway_load_balancer
where
 state_code <> 'active';
```
