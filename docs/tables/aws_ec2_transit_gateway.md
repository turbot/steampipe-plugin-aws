# Table: aws_ec2_transit_gateway

A transit gateway is a network transit hub that you can use to interconnect your virtual private clouds (VPC) and on-premises networks.

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