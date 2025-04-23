---
title: "Steampipe Table: aws_vpc - Query AWS VPC using SQL"
description: "Allows users to query VPCs within AWS. It provides information about each VPC's configuration, including its ID, state, CIDR block, and whether it is the default VPC."
folder: "VPC"
---

# Table: aws_vpc - Query AWS VPC using SQL

The AWS Virtual Private Cloud (VPC) allows you to launch AWS resources in a virtual network that you've defined. This virtual network closely resembles a traditional network that you'd operate in your own data center, with the benefits of using the scalable infrastructure of AWS. It provides advanced security features, such as security groups and network access control lists, to enable inbound and outbound filtering at the instance and subnet level.

## Table Usage Guide

The `aws_vpc` table in Steampipe provides you with information about Virtual Private Clouds (VPCs) within Amazon Web Services (AWS). This table allows you, as a network administrator or DevOps engineer, to query VPC-specific details, including its ID, state, CIDR block, and whether it is the default VPC. You can utilize this table to gather insights on VPCs, such as their networking configuration, security settings, and associated resources. The schema outlines the various attributes of the VPC for you, including the VPC ID, state, CIDR block, default VPC status, and associated tags.

## Examples

### Find default VPCs
Explore which Virtual Private Clouds (VPCs) are set as default within your AWS account. This is beneficial to understand your network configuration and to identify any potential security issues related to default settings.

```sql+postgres
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

```sql+sqlite
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
  is_default = 1;
```


### Show CIDR details
Explore the details of your virtual private cloud (VPC) to gain insights into its network characteristics such as host addresses, broadcast addresses, and network masks. This can be useful in understanding the structure and scope of your VPC's network for better resource allocation and network planning.

```sql+postgres
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

```sql+sqlite
Error: SQLite does not support CIDR operations.
```


### List VPCs with public CIDR blocks
Explore VPCs that are configured with public IP ranges instead of the recommended private ranges. This query can be used to identify potential security risks in your AWS environment.

```sql+postgres
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

```sql+sqlite
Error: SQLite does not support CIDR operations
```