---
title: "Table: aws_ec2_transit_gateway - Query AWS EC2 Transit Gateway using SQL"
description: "Allows users to query AWS EC2 Transit Gateway resources for detailed information on configuration, status, and associations."
---

# Table: aws_ec2_transit_gateway - Query AWS EC2 Transit Gateway using SQL

The `aws_ec2_transit_gateway` table in Steampipe provides information about Transit Gateways within Amazon Elastic Compute Cloud (EC2). This table allows DevOps engineers to query Transit Gateway-specific details, including its configuration, state, and associations. Users can utilize this table to gather insights on Transit Gateways, such as its attached VPCs, VPN connections, Direct Connect gateways, and more. The schema outlines the various attributes of the Transit Gateway, including the transit gateway ID, creation time, state, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_transit_gateway` table, you can use the `.inspect aws_ec2_transit_gateway` command in Steampipe.

Key columns:

- `transit_gateway_id`: The ID of the transit gateway. This can be used to join with other tables that reference the transit gateway ID, such as `aws_ec2_transit_gateway_route_table'.
- `state`: The state of the transit gateway. This is useful for filtering transit gateways based on their state (available, deleted, etc.).
- `creation_time`: The time the transit gateway was created. This can be useful for auditing or tracking the lifecycle of transit gateways.

## Examples

### Basic Transit Gateway info

```sql
select
  transit_gateway_id,
  state,
  owner_id,
  creation_time
from
  aws_ec2_transit_gateway;
```


### List transit gateways which automatically accepts shared account attachment

```sql
select
  transit_gateway_id,
  auto_accept_shared_attachments
from
  aws_ec2_transit_gateway
where
  auto_accept_shared_attachments = 'enable';
```


### Find the number of transit gateways by default route table id

```sql
select
  association_default_route_table_id,
  count(transit_gateway_id) as transit_gateway
from
  aws_ec2_transit_gateway
group by
  association_default_route_table_id;
```


### Map all transit gateways to the application to which they belong with an application tag

```sql
select
  transit_gateway_id,
  tags
from
  aws_ec2_transit_gateway
where
  not tags :: JSONB ? 'application';
```