---
title: "Steampipe Table: aws_ec2_transit_gateway_route - Query AWS EC2 Transit Gateway Routes using SQL"
description: "Allows users to query AWS EC2 Transit Gateway Routes for detailed information about each route, including the destination CIDR block, the route's current state, and the transit gateway attachments."
folder: "EC2"
---

# Table: aws_ec2_transit_gateway_route - Query AWS EC2 Transit Gateway Routes using SQL

The AWS EC2 Transit Gateway Routes enable you to manage connectivity between multiple Virtual Private Clouds (VPCs) and on-premises networks by acting as a hub. They simplify network architecture by reducing the number of connections required to connect multiple VPCs and on-premises networks. Transit Gateway Routes also provide flexible routing policies to support various types of network architectures.

## Table Usage Guide

The `aws_ec2_transit_gateway_route` table in Steampipe provides you with information about the routes in each transit gateway within AWS EC2. This table allows you, as a DevOps engineer, to query route-specific details, including the destination CIDR block, the route's current state, and the transit gateway attachments. You can utilize this table to gather insights on routes, such as verifying the transit gateway route's state, checking the destination CIDR block, and more. The schema outlines the various attributes of the transit gateway route for you, including the transit gateway route ID, transit gateway route destination CIDR block, and associated tags.

## Examples

### Basic info
Explore the configuration of your AWS EC2 transit gateway routes to understand their current state and type. This can help you identify potential network routing issues or areas for optimization.

```sql+postgres
select
  transit_gateway_route_table_id,
  destination_cidr_block,
  prefix_list_id,
  state,
  type
from
  aws_ec2_transit_gateway_route;
```

```sql+sqlite
select
  transit_gateway_route_table_id,
  destination_cidr_block,
  prefix_list_id,
  state,
  type
from
  aws_ec2_transit_gateway_route;
```

### List active routes
Explore which transit gateway routes are currently active to manage network traffic effectively. This is useful for maintaining network efficiency and ensuring optimal route configurations.

```sql+postgres
select
  transit_gateway_route_table_id,
  destination_cidr_block,
  state,
  type
from
  aws_ec2_transit_gateway_route
where
  state = 'active';
```

```sql+sqlite
select
  transit_gateway_route_table_id,
  destination_cidr_block,
  state,
  type
from
  aws_ec2_transit_gateway_route
where
  state = 'active';
```