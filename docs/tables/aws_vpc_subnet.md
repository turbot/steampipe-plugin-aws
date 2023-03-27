# Table: aws_vpc_subnet

AWS VPC Subnet is a logical subdivision of an IP network. It enables dividing a network into two or more networks.

## Examples

### Basic VPC subnet IP address info

```sql
select
  vpc_id,
  subnet_id,
  cidr_block,
  assign_ipv6_address_on_creation,
  map_customer_owned_ip_on_launch,
  map_public_ip_on_launch,
  ipv6_cidr_block_association_set
from
  aws_vpc_subnet;
```


### Availability zone info for each subnet in a VPC

```sql
select
  vpc_id,
  subnet_id,
  availability_zone,
  availability_zone_id
from
  aws_vpc_subnet
order by
  vpc_id,
  availability_zone;
```


### Find the number of available IP address in each subnet

```sql
select
  subnet_id,
  cidr_block,
  available_ip_address_count,
  power(2, 32 - masklen(cidr_block :: cidr)) -1 as raw_size
from
  aws_vpc_subnet;
```


### Route table associated with each subnet

```sql
select
  associations_detail ->> 'SubnetId' as subnet_id,
  route_table_id
from
  aws_vpc_route_table as rt
  cross join jsonb_array_elements(associations) as associations_detail
  join aws_vpc_subnet as sub on sub.subnet_id = associations_detail ->> 'SubnetId';
```


### Subnet count by VPC ID

```sql
select
  vpc_id,
  count(subnet_id) as subnet_count
from
  aws_vpc_subnet
group by
  vpc_id;
```
