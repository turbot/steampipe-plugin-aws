# Table: aws_vpc

A VPC is a virtual network in Amazon AWS.

## Examples

### Find default VPCs

```sql
select
  vpc_id,
  is_default,
  cidr_block,
  state,
  account_id,
  region
from
  aws_vpc
where
  is_default;
```


### Show CIDR details

```sql
select
  vpc_id,
  cidr_block,
  host(cidr_block),
  broadcast(cidr_block),
  netmask(cidr_block),
  network(cidr_block)
from
  aws_vpc;
```


### List VPCs with public CIDR blocks

```sql
select
  vpc_id,
  cidr_block,
  state,
  region
from
  aws_vpc
where
  not cidr_block <<= '10.0.0.0/8'
  and not cidr_block <<= '192.168.0.0/16'
  and not cidr_block <<= '172.16.0.0/12';
```