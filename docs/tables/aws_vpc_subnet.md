---
title: "Table: aws_vpc_subnet - Query AWS VPC Subnets using SQL"
description: "Allows users to query AWS VPC Subnets and obtain detailed information about each subnet, including its configuration, associated VPC, availability zone, and CIDR block."
---

# Table: aws_vpc_subnet - Query AWS VPC Subnets using SQL

The `aws_vpc_subnet` table in Steampipe provides information about subnets within Amazon Virtual Private Cloud (VPC). This table allows DevOps engineers to query subnet-specific details, including its configuration, associated VPC, availability zone, and CIDR block. Users can utilize this table to gather insights on subnets, such as subnet size, associated route tables, network ACLs, and more. The schema outlines the various attributes of the subnet, including the subnet ID, VPC ID, state, CIDR block, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_subnet` table, you can use the `.inspect aws_vpc_subnet` command in Steampipe.

**Key columns**:

- `subnet_id`: This is the unique identifier for the subnet. It is useful for joining with other tables that reference subnets.
- `vpc_id`: This column contains the ID of the VPC the subnet is a part of. It can be used to join with the `aws_vpc` table to get more information about the VPC.
- `cidr_block`: This column contains the IPv4 network range for the subnet, which can be useful for identifying overlapping subnets or planning network architecture.

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
