---
title: "Table: aws_vpc_route - Query AWS VPC Routes using SQL"
description: "Allows users to query AWS VPC Routes to retrieve detailed information about each route in a route table within a VPC."
---

# Table: aws_vpc_route - Query AWS VPC Routes using SQL

The `aws_vpc_route` table in Steampipe provides information about each route in a route table within a VPC. This table allows DevOps engineers to query route-specific details, including the destination CIDR block, the ID of the route table the route is in, and the type of target (e.g., internet gateway, virtual private gateway, etc.). Users can utilize this table to gather insights on routes, such as verifying route configurations, checking route targets, and examining route propagation. The schema outlines the various attributes of the route, including the destination CIDR block, route table ID, and associated targets.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_route` table, you can use the `.inspect aws_vpc_route` command in Steampipe.

### Key columns:

- `destination_cidr_block`: This is the IPv4 network range for the destination of the route. It is a key column as it is used to identify the network range that the route applies to.
- `route_table_id`: This is the ID of the route table which contains the route. It is important as it can be used to join this table with the aws_vpc_route_table.
- `target_type`: This column specifies the type of target for the route (e.g., internet gateway, virtual private gateway, etc.). It is useful for understanding where the traffic is being routed.


## Examples

### List of route tables whose routes are directed to the internet

```sql
select
  route_table_id,
  gateway_id
from
  aws_vpc_route
where
  gateway_id ilike 'igw%'
  and destination_cidr_block = '0.0.0.0/0';
```


### List of route tables whose route target is not available

```sql
select
  route_table_id,
  state
from
  aws_vpc_route
where
  state = 'blackhole';
```


### Routing details for each route table

```sql
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