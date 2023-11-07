# Table: aws_vpc_nat_gateway

NAT Gateway is a highly available AWS managed service that makes it easy to connect to the Internet from instances within a private subnet in an Amazon Virtual Private Cloud (Amazon VPC).

## Examples

### IP address details of the NAT gateway

```sql
select
  nat_gateway_id,
  address ->> 'PrivateIp' as private_ip,
  address ->> 'PublicIp' as public_ip,
  address ->> 'NetworkInterfaceId' as nic_id,
  address ->> 'AllocationId' as allocation_id
from
  aws_vpc_nat_gateway
  cross join jsonb_array_elements(nat_gateway_addresses) as address;
```


### VPC details associated with the NAT gateway

```sql
select
  nat_gateway_id,
  vpc_id,
  subnet_id
from
  aws_vpc_nat_gateway;
```


### List NAT gateways without application tags key

```sql
select
  nat_gateway_id,
  tags
from
  aws_vpc_nat_gateway
where
  not tags :: JSONB ? 'application';
```


### Count of NAT gateways by VPC Id

```sql
select
  vpc_id,
  count(nat_gateway_id) as nat_gateway_id
from
  aws_vpc_nat_gateway
group by
  vpc_id;
```
