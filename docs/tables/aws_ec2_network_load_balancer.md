---
title: "Table: aws_ec2_network_load_balancer - Query AWS EC2 Network Load Balancer using SQL"
description: "Allows users to query AWS EC2 Network Load Balancer data including configuration, status, and other related information."
---

# Table: aws_ec2_network_load_balancer - Query AWS EC2 Network Load Balancer using SQL

The `aws_ec2_network_load_balancer` table in Steampipe provides information about Network Load Balancers within AWS Elastic Compute Cloud (EC2). This table allows cloud administrators and DevOps engineers to query load balancer-specific details, including type, state, availability zones, and associated metadata. Users can utilize this table to gather insights on load balancers, such as their current status, associated subnets, and more. The schema outlines the various attributes of the Network Load Balancer, including the load balancer name, ARN, creation date, DNS name, scheme, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_network_load_balancer` table, you can use the `.inspect aws_ec2_network_load_balancer` command in Steampipe.

### Key columns:

- `name`: The name of the load balancer. This can be used to join with other tables that contain load balancer-specific information.
- `arn`: The Amazon Resource Name (ARN) of the load balancer. This unique identifier can be used to join with other tables that reference Network Load Balancers.
- `state_code`: The current state of the load balancer. This is useful for filtering load balancers based on their operational status.

## Examples

### Count of AZs registered with network load balancers

```sql
select
  name,
  count(az ->> 'ZoneName') as zone_count
from
  aws_ec2_network_load_balancer
  cross join jsonb_array_elements(availability_zones) as az
group by
  name;
```


### List of network load balancers where Cross-Zone Load Balancing is enabled

```sql
select
  name,
  lb ->> 'Key' as cross_zone,
  lb ->> 'Value' as cross_zone_value
from
  aws_ec2_network_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'load_balancing.cross_zone.enabled'
  and lb ->> 'Value' = 'false';
```


### List of network load balancers where logging is not enabled

```sql
select
  name,
  lb ->> 'Key' as logging_key,
  lb ->> 'Value' as logging_value
from
  aws_ec2_network_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'access_logs.s3.enabled'
  and lb ->> 'Value' = 'false';
```


### List of network load balancers where deletion protection is not enabled

```sql
select
  name,
  lb ->> 'Key' as deletion_protection_key,
  lb ->> 'Value' as deletion_protection_value
from
  aws_ec2_network_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'deletion_protection.enabled'
  and lb ->> 'Value' = 'false';
```
