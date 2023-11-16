---
title: "Table: aws_ec2_network_interface - Query AWS EC2 Network Interfaces using SQL"
description: "Allows users to query AWS EC2 Network Interfaces and provides comprehensive details about each interface, including its associated instances, security groups, and subnet information."
---

# Table: aws_ec2_network_interface - Query AWS EC2 Network Interfaces using SQL

The `aws_ec2_network_interface` table in Steampipe provides information about Network Interfaces within AWS Elastic Compute Cloud (EC2). This table allows DevOps engineers to query network interface-specific details, including the attached instances, associated security groups, subnet information, and more. Users can utilize this table to gather insights on network interfaces, such as their status, type, private and public IP addresses, and the associated subnet and VPC details. The schema outlines the various attributes of the EC2 network interface, including the interface ID, description, owner ID, availability zone, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_network_interface` table, you can use the `.inspect aws_ec2_network_interface` command in Steampipe.

### Key columns:

- `network_interface_id`: The ID of the network interface. This column can be used to join this table with other tables that contain network interface information.
- `subnet_id`: The ID of the subnet for the network interface. This column can be used to join this table with other tables that contain subnet information.
- `vpc_id`: The ID of the VPC for the network interface. This column can be used to join this table with other tables that contain VPC information.

## Examples

### Basic IP address info

```sql
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

```sql
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

### Count of ENIs by interface type

```sql
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

```sql
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

### Get network details for each ENI

```sql
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