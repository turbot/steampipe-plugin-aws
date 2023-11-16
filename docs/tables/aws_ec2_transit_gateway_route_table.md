---
title: "Table: aws_ec2_transit_gateway_route_table - Query AWS EC2 Transit Gateway Route Tables using SQL"
description: "Allows users to query AWS EC2 Transit Gateway Route Tables and retrieve detailed information about each route table, including its ID, state, transit gateway ID, and other associated metadata."
---

# Table: aws_ec2_transit_gateway_route_table - Query AWS EC2 Transit Gateway Route Tables using SQL

The `aws_ec2_transit_gateway_route_table` table in Steampipe provides information about each route table associated with a transit gateway within the Amazon Elastic Compute Cloud (EC2). This table allows DevOps engineers to query route table-specific details, including the transit gateway ID, route table ID, state, and associated tags. Users can utilize this table to gather insights on transit gateway route tables, such as their current state, associated transit gateways, and more. The schema outlines the various attributes of the transit gateway route table, including the route table ID, transit gateway ID, creation time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_transit_gateway_route_table` table, you can use the `.inspect aws_ec2_transit_gateway_route_table` command in Steampipe.

**Key columns**:

- `transit_gateway_route_table_id`: The ID of the transit gateway route table. This is a unique identifier that can be used to join this table with other tables and retrieve specific information about a route table.
- `transit_gateway_id`: The ID of the transit gateway associated with the route table. This can be useful when querying for information about a specific transit gateway.
- `state`: The state of the transit gateway route table (available, pending, deleting, etc.). This can be useful when monitoring the status of transit gateway route tables.


## Examples

### Basic transit gateway route table info

```sql
select
  transit_gateway_route_table_id,
  transit_gateway_id,
  default_association_route_table,
  default_propagation_route_table
from
  aws_ec2_transit_gateway_route_table;
```


### Count of transit gateway route table by transit gateway

```sql
select
  transit_gateway_id,
  count(transit_gateway_route_table_id) as transit_gateway_route_table_count
from
  aws_ec2_transit_gateway_route_table
group by
  transit_gateway_id;
```