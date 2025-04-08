---
title: "Steampipe Table: aws_vpc_nat_gateway - Query AWS VPC NAT Gateways using SQL"
description: "Allows users to query NAT Gateways within Amazon Virtual Private Cloud (VPC). The `aws_vpc_nat_gateway` table in Steampipe provides information about each NAT Gateway within a VPC. This table can be used to gather insights on NAT Gateways, such as their state, subnet association, and associated Elastic IP addresses."
folder: "VPC"
---

# Table: aws_vpc_nat_gateway - Query AWS VPC NAT Gateways using SQL

An AWS VPC NAT Gateway is a highly available, managed Network Address Translation (NAT) service for your resources in a private subnet to access the internet. It is designed to automatically scale up to the bandwidth you need, and you only pay for the amount of traffic processed. The NAT gateway handles all traffic leaving your VPC and routes it to the internet.

## Table Usage Guide

The `aws_vpc_nat_gateway` table in Steampipe provides you with information about each NAT Gateway within Amazon Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, to query NAT Gateway-specific details, including its current state, the subnet it is associated with, and any associated Elastic IP addresses. You can utilize this table to verify the configuration and status of NAT Gateways, ensuring they are properly connected and functioning within your VPC. The schema outlines the various attributes of the NAT Gateway for you, including the NAT Gateway ID, creation time, state, subnet ID, and associated IP addresses.

## Examples

### IP address details of the NAT gateway
Determine the private and public IP addresses associated with your NAT gateway to manage network traffic and enhance security. This can also help in identifying the network interface and allocation IDs for better resource management.

```sql+postgres
select
  nat_gateway_id,
  address ->> 'PrivateIp' as private_ip,
  address ->> 'PublicIp' as public_ip,
  address ->> 'NetworkInterfaceId' as nic_id,
  address ->> 'AllocationId' as allocation_id
from
  aws_vpc_nat_gateway
  cross join jsonb_array_elements(nat_gateway_addresses) as address;
```

```sql+sqlite
select
  nat_gateway_id,
  json_extract(address.value, '$.PrivateIp') as private_ip,
  json_extract(address.value, '$.PublicIp') as public_ip,
  json_extract(address.value, '$.NetworkInterfaceId') as nic_id,
  json_extract(address.value, '$.AllocationId') as allocation_id
from
  aws_vpc_nat_gateway,
  json_each(nat_gateway_addresses) as address;
```

### VPC details associated with the NAT gateway
Explore the relationship between your NAT gateway and associated VPC details to understand the network architecture better. This can be particularly useful in managing and optimizing your cloud resources.

```sql+postgres
select
  nat_gateway_id,
  vpc_id,
  subnet_id
from
  aws_vpc_nat_gateway;
```

```sql+sqlite
select
  nat_gateway_id,
  vpc_id,
  subnet_id
from
  aws_vpc_nat_gateway;
```

### List NAT gateways without application tags key
Discover the segments of your network that lack application tags on their NAT gateways. This can help ensure comprehensive tagging, improving network management and cost allocation.

```sql+postgres
select
  nat_gateway_id,
  tags
from
  aws_vpc_nat_gateway
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  nat_gateway_id,
  tags
from
  aws_vpc_nat_gateway
where
  json_extract(tags, '$.application') IS NULL;
```


### Count of NAT gateways by VPC Id
Determine the number of Network Address Translation (NAT) gateways associated with each Virtual Private Cloud (VPC) to better manage and optimize your network resources.

```sql+postgres
select
  vpc_id,
  count(nat_gateway_id) as nat_gateway_id
from
  aws_vpc_nat_gateway
group by
  vpc_id;
```

```sql+sqlite
select
  vpc_id,
  count(nat_gateway_id) as nat_gateway_id
from
  aws_vpc_nat_gateway
group by
  vpc_id;
```