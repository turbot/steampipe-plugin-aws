# Table: aws_vpc_route_table

A route table contains a set of rules, called routes, that are used to determine where network traffic from your subnet or gateway is directed.

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
