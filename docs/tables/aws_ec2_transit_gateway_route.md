---
title: "Table: aws_ec2_transit_gateway_route - Query AWS EC2 Transit Gateway Routes using SQL"
description: "Allows users to query AWS EC2 Transit Gateway Routes for detailed information about each route, including the destination CIDR block, the route's current state, and the transit gateway attachments."
---

# Table: aws_ec2_transit_gateway_route - Query AWS EC2 Transit Gateway Routes using SQL

The `aws_ec2_transit_gateway_route` table in Steampipe provides information about the routes in each transit gateway within AWS EC2. This table allows DevOps engineers to query route-specific details, including the destination CIDR block, the route's current state, and the transit gateway attachments. Users can utilize this table to gather insights on routes, such as verifying the transit gateway route's state, checking the destination CIDR block, and more. The schema outlines the various attributes of the transit gateway route, including the transit gateway route ID, transit gateway route destination CIDR block, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_transit_gateway_route` table, you can use the `.inspect aws_ec2_transit_gateway_route` command in Steampipe.

**Key columns**:

- `transit_gateway_route_id`: This is the ID of the transit gateway route. It is a key column for joining with other tables as it uniquely identifies each route.
- `destination_cidr_block`: This column is used to specify the range of IP addresses for the destination network. It can be useful when filtering routes based on the destination network.
- `transit_gateway_attachment_id`: This column contains the ID of the transit gateway attachment. It is crucial for understanding which attachments are associated with each route.

## Examples

### Basic info

```sql
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

```sql
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
