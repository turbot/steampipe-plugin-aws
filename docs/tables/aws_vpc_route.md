---
title: "Steampipe Table: aws_vpc_route - Query AWS VPC Routes using SQL"
description: "Allows users to query AWS VPC Routes to retrieve detailed information about each route in a route table within a VPC."
folder: "VPC"
---

# Table: aws_vpc_route - Query AWS VPC Routes using SQL

The AWS VPC Route is a component of Amazon Virtual Private Cloud (VPC) that allows network traffic to be directed from a subnet route table to a specific network gateway or instance. It provides the ability to control the navigational path for outbound traffic. This is crucial for managing the accessibility of network interfaces and ensuring the secure transmission of data within your AWS environment.

## Table Usage Guide

The `aws_vpc_route` table in Steampipe gives you information about each route in a route table within a VPC. This table allows you, as a DevOps engineer, to query route-specific details, including the destination CIDR block, the ID of the route table the route is in, and the type of target (e.g., internet gateway, virtual private gateway, etc.). You can utilize this table to gather insights on routes, such as verifying route configurations, checking route targets, and examining route propagation. The schema outlines the various attributes of the route for you, including the destination CIDR block, route table ID, and associated targets.

## Examples

### List of route tables whose routes are directed to the internet
Discover the segments of your network that are directly connected to the internet. This is useful for identifying potential security risks and ensuring that your network configuration aligns with your company's policies.

```sql+postgres
select
  route_table_id,
  gateway_id
from
  aws_vpc_route
where
  gateway_id ilike 'igw%'
  and destination_cidr_block = '0.0.0.0/0';
```

```sql+sqlite
select
  route_table_id,
  gateway_id
from
  aws_vpc_route
where
  gateway_id like 'igw%'
  and destination_cidr_block = '0.0.0.0/0';
```


### List of route tables whose route target is not available
Determine the areas in which certain route tables are in a 'blackhole' state, indicating that their route target is not available. This query can be useful in identifying potential network connectivity issues within your AWS Virtual Private Cloud (VPC).

```sql+postgres
select
  route_table_id,
  state
from
  aws_vpc_route
where
  state = 'blackhole';
```

```sql+sqlite
select
  route_table_id,
  state
from
  aws_vpc_route
where
  state = 'blackhole';
```


### Routing details for each route table
Explore the routing configurations for each route within your network to gain insights into their status and associated destinations. This can be helpful in assessing network traffic paths and identifying any potential bottlenecks or issues.

```sql+postgres
select
  route_table_id,
  state,
  destination_cidr_block,
  destination_ipv6_cidr_block,
  carrier_gateway_id,
  destination_prefix_list_id,
  egress_only_internet_gateway_id,
  gateway_id,
  instance_id,
  nat_gateway_id,
  network_interface_id,
  transit_gateway_id,
  vpc_peering_connection_id
from
  aws_vpc_route;
```

```sql+sqlite
select
  route_table_id,
  state,
  destination_cidr_block,
  destination_ipv6_cidr_block,
  carrier_gateway_id,
  destination_prefix_list_id,
  egress_only_internet_gateway_id,
  gateway_id,
  instance_id,
  nat_gateway_id,
  network_interface_id,
  transit_gateway_id,
  vpc_peering_connection_id
from
  aws_vpc_route;
```