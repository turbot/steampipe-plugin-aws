---
title: "Steampipe Table: aws_ec2_transit_gateway_route_table - Query AWS EC2 Transit Gateway Route Tables using SQL"
description: "Allows users to query AWS EC2 Transit Gateway Route Tables and retrieve detailed information about each route table, including its ID, state, transit gateway ID, and other associated metadata."
folder: "EC2"
---

# Table: aws_ec2_transit_gateway_route_table - Query AWS EC2 Transit Gateway Route Tables using SQL

The AWS EC2 Transit Gateway Route Table is a component of Amazon's Elastic Compute Cloud (EC2) service that allows you to manage routing for your Transit Gateways. It facilitates the control of traffic between different networks within your cloud environment. Using this resource, you can define rules that determine the path network traffic takes to reach a specific destination.

## Table Usage Guide

The `aws_ec2_transit_gateway_route_table` table in Steampipe provides you with information about each route table associated with a transit gateway within your Amazon Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query route table-specific details, including the transit gateway ID, route table ID, state, and associated tags. You can utilize this table to gather insights on transit gateway route tables, such as their current state, associated transit gateways, and more. The schema outlines the various attributes of the transit gateway route table for you, including the route table ID, transit gateway ID, creation time, and associated tags.

## Examples

### Basic transit gateway route table info
Explore the fundamental characteristics of your transit gateway route tables in AWS EC2. This query is useful in understanding the default associations and propagations within your route tables, aiding in efficient network management.

```sql+postgres
select
  transit_gateway_route_table_id,
  transit_gateway_id,
  default_association_route_table,
  default_propagation_route_table
from
  aws_ec2_transit_gateway_route_table;
```

```sql+sqlite
select
  transit_gateway_route_table_id,
  transit_gateway_id,
  default_association_route_table,
  default_propagation_route_table
from
  aws_ec2_transit_gateway_route_table;
```


### Count of transit gateway route table by transit gateway
Explore which transit gateways are associated with numerous route tables in your AWS EC2 service. This can be useful for optimizing network routing paths and managing network resources effectively.

```sql+postgres
select
  transit_gateway_id,
  count(transit_gateway_route_table_id) as transit_gateway_route_table_count
from
  aws_ec2_transit_gateway_route_table
group by
  transit_gateway_id;
```

```sql+sqlite
select
  transit_gateway_id,
  count(transit_gateway_route_table_id) as transit_gateway_route_table_count
from
  aws_ec2_transit_gateway_route_table
group by
  transit_gateway_id;
```