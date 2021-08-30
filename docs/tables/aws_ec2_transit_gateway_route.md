# Table: aws_ec2_transit_gateway_route

Transit gateway routes IPv4 and IPv6 packets between attachments using transit gateway route tables. You can configure these route tables to propagate routes from the route tables for the attached VPCs, VPN connections, and Direct Connect gateways.

## Examples

### Basic transit gateway route info

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