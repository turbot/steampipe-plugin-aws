---
title: "Steampipe Table: aws_ec2_network_interface - Query AWS EC2 Network Interfaces using SQL"
description: "Allows users to query AWS EC2 Network Interfaces and provides comprehensive details about each interface, including its associated instances, security groups, and subnet information."
folder: "EC2"
---

# Table: aws_ec2_network_interface - Query AWS EC2 Network Interfaces using SQL

An AWS EC2 Network Interface is a virtual network interface that you can attach to an instance in a VPC. Network interfaces are the point of networking for any instance that is attached to a Virtual Private Cloud (VPC). They can include a primary private IPv4 address, one or more secondary private IPv4 addresses, one Elastic IP address per private IPv4 address, one public IPv4 address, one or more IPv6 addresses, a MAC address, one or more security groups, a source/destination check flag, and a description.

## Table Usage Guide

The `aws_ec2_network_interface` table in Steampipe provides you with information about Network Interfaces within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query network interface-specific details, including the attached instances, associated security groups, subnet information, and more. You can utilize this table to gather insights on network interfaces, such as their status, type, private and public IP addresses, and the associated subnet and VPC details. The schema outlines for you the various attributes of the EC2 network interface, including the interface ID, description, owner ID, availability zone, and associated tags.

## Examples

### Basic IP address info
Determine the areas in which your AWS EC2 network interfaces are operating by exploring the type of interface, its corresponding private and public IP addresses, and its MAC address. This can be particularly useful for managing network connectivity and troubleshooting network issues within your AWS environment.

```sql+postgres
select
  network_interface_id,
  interface_type,
  description,
  private_ip_address,
  association_public_ip,
  mac_address
from
  aws_ec2_network_interface;
```

```sql+sqlite
select
  network_interface_id,
  interface_type,
  description,
  private_ip_address,
  association_public_ip,
  mac_address
from
  aws_ec2_network_interface;
```

### Find all ENIs with private IPs that are in a given subnet (10.66.0.0/16)
Discover the segments that have private IPs within a specific subnet. This is useful for identifying network interfaces within a particular subnet, which can aid in network management and security assessment.

```sql+postgres
select
  network_interface_id,
  interface_type,
  description,
  private_ip_address,
  association_public_ip,
  mac_address
from
  aws_ec2_network_interface
where
  private_ip_address :: cidr <<= '10.66.0.0/16';
```

```sql+sqlite
Error: SQLite does not support CIDR operations.
```

### Count of ENIs by interface type
Discover the segments that have the most network interfaces in your AWS EC2 environment, helping you understand your network configuration and potentially optimize resource allocation.

```sql+postgres
select
  interface_type,
  count(interface_type) as count
from
  aws_ec2_network_interface
group by
  interface_type
order by
  count desc;
```

```sql+sqlite
select
  interface_type,
  count(interface_type) as count
from
  aws_ec2_network_interface
group by
  interface_type
order by
  count desc;
```

### Security groups attached to each ENI
Determine the areas in which certain security groups are attached to each network interface within your Amazon EC2 instances. This can help in managing security and access controls effectively.

```sql+postgres
select
  network_interface_id as eni,
  sg ->> 'GroupId' as "security group id",
  sg ->> 'GroupName' as "security group name"
from
  aws_ec2_network_interface
  cross join jsonb_array_elements(groups) as sg
order by
  eni;
```

```sql+sqlite
select
  network_interface_id as eni,
  json_extract(sg, '$.GroupId') as "security group id",
  json_extract(sg, '$.GroupName') as "security group name"
from
  (
    select
      network_interface_id,
      json_each.value as sg
    from
      aws_ec2_network_interface,
      json_each(groups)
  )
order by
  eni;
```

### Get network details for each ENI
Discover the segments that are common between your network interfaces and virtual private clouds (VPCs) to better understand your network structure. This can assist in identifying areas for potential consolidation or optimization.

```sql+postgres
select
  e.network_interface_id,
  v.vpc_id,
  v.is_default,
  v.cidr_block,
  v.state,
  v.account_id,
  v.region
from
  aws_ec2_network_interface e,
  aws_vpc v
where 
  e.vpc_id = v.vpc_id;
```

```sql+sqlite
select
  e.network_interface_id,
  v.vpc_id,
  v.is_default,
  v.cidr_block,
  v.state,
  v.account_id,
  v.region
from
  aws_ec2_network_interface e
join
  aws_vpc v
on 
  e.vpc_id = v.vpc_id;
```