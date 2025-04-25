---
title: "Steampipe Table: aws_ec2_network_load_balancer - Query AWS EC2 Network Load Balancer using SQL"
description: "Allows users to query AWS EC2 Network Load Balancer data including configuration, status, and other related information."
folder: "ELB"
---

# Table: aws_ec2_network_load_balancer - Query AWS EC2 Network Load Balancer using SQL

The AWS EC2 Network Load Balancer is a high-performance load balancer that operates at the transport layer (Layer 4) and is designed to handle millions of requests per second while maintaining ultra-low latencies. It is best suited for load balancing of TCP traffic and capable of handling volatile workloads and traffic patterns. It also supports long-lived TCP connections, which are ideal for WebSocket type of applications.

## Table Usage Guide

The `aws_ec2_network_load_balancer` table in Steampipe provides you with information about Network Load Balancers within AWS Elastic Compute Cloud (EC2). This table allows you, as a cloud administrator or DevOps engineer, to query load balancer-specific details, including type, state, availability zones, and associated metadata. You can utilize this table to gather insights on load balancers, such as their current status, associated subnets, and more. The schema outlines the various attributes of the Network Load Balancer for you, including the load balancer name, ARN, creation date, DNS name, scheme, and associated tags.

## Examples

### Count of AZs registered with network load balancers
Analyze the distribution of network load balancers across various availability zones to optimize resource allocation and ensure a balanced load. This can help in enhancing the application's performance and availability.

```sql+postgres
select
  name,
  count(az ->> 'ZoneName') as zone_count
from
  aws_ec2_network_load_balancer
  cross join jsonb_array_elements(availability_zones) as az
group by
  name;
```

```sql+sqlite
select
  name,
  count(json_extract(az.value, '$.ZoneName')) as zone_count
from
  aws_ec2_network_load_balancer,
  json_each(availability_zones) as az
group by
  name;
```


### List of network load balancers where Cross-Zone Load Balancing is enabled
Determine the areas in which Cross-Zone Load Balancing is enabled for network load balancers. This can be particularly useful to identify potential areas of network inefficiency or to optimize load balancing across zones.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(lb.value, '$.Key') as cross_zone,
  json_extract(lb.value, '$.Value') as cross_zone_value
from
  aws_ec2_network_load_balancer,
  json_each(load_balancer_attributes) as lb
where
  json_extract(lb.value, '$.Key') = 'load_balancing.cross_zone.enabled'
  and json_extract(lb.value, '$.Value') = 'false';
```


### List of network load balancers where logging is not enabled
Determine the areas in your network load balancers where logging is not enabled. This is essential for identifying potential security risks and ensuring compliance with data governance policies.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(lb.value, '$.Key') as logging_key,
  json_extract(lb.value, '$.Value') as logging_value
from
  aws_ec2_network_load_balancer,
  json_each(load_balancer_attributes) as lb
where
  json_extract(lb.value, '$.Key') = 'access_logs.s3.enabled'
  and json_extract(lb.value, '$.Value') = 'false';
```


### List of network load balancers where deletion protection is not enabled
Determine the areas in your network where load balancers are potentially vulnerable due to deletion protection not being enabled. This is particularly useful for identifying potential risks and ensuring the security and stability of your network.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(lb.value, '$.Key') as deletion_protection_key,
  json_extract(lb.value, '$.Value') as deletion_protection_value
from
  aws_ec2_network_load_balancer,
  json_each(load_balancer_attributes) as lb
where
  json_extract(lb.value, '$.Key') = 'deletion_protection.enabled'
  and json_extract(lb.value, '$.Value') = 'false';
```