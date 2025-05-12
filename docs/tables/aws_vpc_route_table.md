---
title: "Steampipe Table: aws_vpc_route_table - Query AWS VPC Route Tables using SQL"
description: "Allows users to query AWS VPC Route Tables and obtain detailed information about each route table, including its associations, routes, and tags."
folder: "VPC"
---

# Table: aws_vpc_route_table - Query AWS VPC Route Tables using SQL

The AWS VPC Route Tables are essential components of Amazon Virtual Private Cloud (VPC) that control the routing for all subnets within a VPC. They determine where network traffic is directed, enabling you to build a variety of network architectures. Each route in a table specifies a destination CIDR and a target, such as a VPC peering connection, network interface, or a gateway.

## Table Usage Guide

The `aws_vpc_route_table` table in Steampipe provides you with information about VPC Route Tables within Amazon Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, to query route table-specific details, including its associations, routes, and tags. You can utilize this table to gather insights on route tables, such as the subnets associated with each route table, the destinations and targets of each route, and the tags associated with each route table. The schema outlines the various attributes of the VPC Route Table for you, including the route table ID, VPC ID, owner ID, and associated tags.

## Examples

### Route table count by VPC ID
Determine the number of route tables associated with each Virtual Private Cloud (VPC) to manage networking environment effectively. This can aid in understanding the complexity and structure of your network within AWS.

```sql+postgres
select
  vpc_id,
  count(route_table_id) as route_table_count
from
  aws_vpc_route_table
group by
  vpc_id;
```

```sql+sqlite
select
  vpc_id,
  count(route_table_id) as route_table_count
from
  aws_vpc_route_table
group by
  vpc_id;
```


### Subnet and Gateways associated with the route table
Explore the associations between subnets and gateways within a route table in your AWS VPC. This can help you better understand your network configuration and identify potential areas for optimization or troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  route_table_id,
  json_extract(associations_detail.value, '$.AssociationState.State') as state,
  json_extract(associations_detail.value, '$.GatewayId') as gateway_id,
  json_extract(associations_detail.value, '$.SubnetId') as subnet_id,
  json_extract(associations_detail.value, '$.RouteTableAssociationId') as route_table_association_id,
  json_extract(associations_detail.value, '$.Main') as main_route_table
from
  aws_vpc_route_table,
  json_each(associations) as associations_detail;
```


### Routing details for each route table
Analyze the settings to understand the routing details for each route table in your AWS VPC. This can help you gain insights into the configuration of different aspects like gateways, instances, and network interfaces, aiding in network troubleshooting and optimization.

```sql+postgres
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

```sql+sqlite
select
  route_table_id,
  json_extract(route_detail.value, '$.CarrierGatewayId.State') as carrier_gateway_id,
  json_extract(route_detail.value, '$.DestinationCidrBlock') as destination_CIDR_block,
  json_extract(route_detail.value, '$.DestinationIpv6CidrBlock') as destination_ipv6_CIDR_block,
  json_extract(route_detail.value, '$.EgressOnlyInternetGatewayId') as egress_only_internet_gateway,
  json_extract(route_detail.value, '$.GatewayId') as gateway_id,
  json_extract(route_detail.value, '$.InstanceId') as instance_id,
  json_extract(route_detail.value, '$.InstanceOwnerId') as instance_owner_id,
  json_extract(route_detail.value, '$.LocalGatewayId') as local_gateway_id,
  json_extract(route_detail.value, '$.NatGatewayId') as nat_gateway_id,
  json_extract(route_detail.value, '$.NetworkInterfaceId') as network_interface_id,
  json_extract(route_detail.value, '$.TransitGatewayId') as transit_gateway_id,
  json_extract(route_detail.value, '$.VpcPeeringConnectionId') as vpc_peering_connection_id
from
  aws_vpc_route_table,
  json_each(routes) as route_detail;
```