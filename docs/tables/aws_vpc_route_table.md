---
title: "Table: aws_vpc_route_table - Query AWS VPC Route Tables using SQL"
description: "Allows users to query AWS VPC Route Tables and obtain detailed information about each route table, including its associations, routes, and tags."
---

# Table: aws_vpc_route_table - Query AWS VPC Route Tables using SQL

The `aws_vpc_route_table` table in Steampipe provides information about VPC Route Tables within Amazon Virtual Private Cloud (VPC). This table allows DevOps engineers to query route table-specific details, including its associations, routes, and tags. Users can utilize this table to gather insights on route tables, such as the subnets associated with each route table, the destinations and targets of each route, and the tags associated with each route table. The schema outlines the various attributes of the VPC Route Table, including the route table ID, VPC ID, owner ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_route_table` table, you can use the `.inspect aws_vpc_route_table` command in Steampipe.

**Key columns**:

- `route_table_id`: The ID of the route table. This column is a unique identifier and can be used to join this table with other tables.
- `vpc_id`: The ID of the VPC associated with the route table. This column is useful for querying all route tables within a specific VPC.
- `owner_id`: The AWS account ID of the owner of the route table. This column is useful for querying all route tables owned by a specific AWS account.

## Examples

### Route table count by VPC ID

```sql
select
  vpc_id,
  count(route_table_id) as route_table_count
from
  aws_vpc_route_table
group by
  vpc_id;
```


### Subnet and Gateways associated with the route table

```sql
select
  route_table_id,
  associations_detail -> 'AssociationState' ->> 'State' as state,
  associations_detail -> 'GatewayId' as gateway_id,
  associations_detail -> 'SubnetId' as subnet_id,
  associations_detail -> 'RouteTableAssociationId' as route_table_association_id,
  associations_detail -> 'Main' as main_route_table
from
  aws_vpc_route_table
  cross join jsonb_array_elements(associations) as associations_detail;
```


### Routing details for each route table

```sql
select
  route_table_id,
  route_detail -> 'CarrierGatewayId' ->> 'State' as carrier_gateway_id,
  route_detail -> 'DestinationCidrBlock' as destination_CIDR_block,
  route_detail -> 'DestinationIpv6CidrBlock' as destination_ipv6_CIDR_block,
  route_detail -> 'EgressOnlyInternetGatewayId' as egress_only_internet_gateway,
  route_detail -> 'GatewayId' as gateway_id,
  route_detail -> 'InstanceId' as instance_id,
  route_detail -> 'InstanceOwnerId' as instance_owner_id,
  route_detail -> 'LocalGatewayId' as local_gateway_id,
  route_detail -> 'NatGatewayId' as nat_gateway_id,
  route_detail -> 'NetworkInterfaceId' as network_interface_id,
  route_detail -> 'TransitGatewayId' as transit_gateway_id,
  route_detail -> 'VpcPeeringConnectionId' as vpc_peering_connection_id
from
  aws_vpc_route_table
  cross join jsonb_array_elements(routes) as route_detail;
```
