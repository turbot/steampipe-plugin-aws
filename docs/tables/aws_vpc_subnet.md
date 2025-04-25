---
title: "Steampipe Table: aws_vpc_subnet - Query AWS VPC Subnets using SQL"
description: "Allows users to query AWS VPC Subnets and obtain detailed information about each subnet, including its configuration, associated VPC, availability zone, and CIDR block."
folder: "VPC"
---

# Table: aws_vpc_subnet - Query AWS VPC Subnets using SQL

An AWS VPC (Virtual Private Cloud) Subnet is a range of IP addresses in your VPC. It allows you to launch AWS resources into a specified subnet, providing logical separation of resources based on security and operational needs. Subnets can be public, private, or VPN-only, providing flexible networking architecture.

## Table Usage Guide

The `aws_vpc_subnet` table in Steampipe provides you with information about subnets within Amazon Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, to query subnet-specific details, including its configuration, associated VPC, availability zone, and CIDR block. You can utilize this table to gather insights on subnets, such as subnet size, associated route tables, network ACLs, and more. The schema outlines the various attributes of the subnet, including the subnet ID, VPC ID, state, CIDR block, and associated tags for you.

## Examples

### Basic VPC subnet IP address info
Determine the areas in which IP addresses are assigned within your AWS VPC subnets. This allows you to understand how IP address assignment and mapping is configured, which is crucial for managing network accessibility and security.

```sql+postgres
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

```sql+sqlite
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
Explore which subnets are located in each availability zone within a Virtual Private Cloud (VPC). This is useful for understanding your network layout and ensuring a balanced distribution across availability zones for resilience and high availability.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which there are available IP addresses within each subnet. This is useful for understanding where there is capacity for new devices or services.

```sql+postgres
select
  subnet_id,
  cidr_block,
  available_ip_address_count,
  power(2, 32 - masklen(cidr_block :: cidr)) -1 as raw_size
from
  aws_vpc_subnet;
```

```sql+sqlite
Error: SQLite does not support CIDR operations.
```


### Route table associated with each subnet
Explore which route tables are linked to each subnet in your AWS VPC. This can help you understand the routing of network traffic within your virtual private cloud.

```sql+postgres
select
  associations_detail ->> 'SubnetId' as subnet_id,
  route_table_id
from
  aws_vpc_route_table as rt
  cross join jsonb_array_elements(associations) as associations_detail
  join aws_vpc_subnet as sub on sub.subnet_id = associations_detail ->> 'SubnetId';
```

```sql+sqlite
select
  json_extract(associations_detail.value, '$.SubnetId') as subnet_id,
  route_table_id
from
  aws_vpc_route_table as rt
  join json_each(rt.associations) as associations_detail
  join aws_vpc_subnet as sub on sub.subnet_id = json_extract(associations_detail.value, '$.SubnetId');
```


### Subnet count by VPC ID
Assess the distribution of subnets across various VPCs to understand the network segmentation in your AWS environment.

```sql+postgres
select
  vpc_id,
  count(subnet_id) as subnet_count
from
  aws_vpc_subnet
group by
  vpc_id;
```

```sql+sqlite
select
  vpc_id,
  count(subnet_id) as subnet_count
from
  aws_vpc_subnet
group by
  vpc_id;
```