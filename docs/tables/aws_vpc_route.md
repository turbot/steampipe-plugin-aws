# Table: aws_vpc_route

Routes are set of rules that are used to determine where network traffic from the subnet or gateway is directed.

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