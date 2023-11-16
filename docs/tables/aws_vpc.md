---
title: "Table: aws_vpc - Query AWS VPC using SQL"
description: "Allows users to query VPCs within AWS. It provides information about each VPC's configuration, including its ID, state, CIDR block, and whether it is the default VPC."
---

# Table: aws_vpc - Query AWS VPC using SQL

The `aws_vpc` table in Steampipe provides information about Virtual Private Clouds (VPCs) within Amazon Web Services (AWS). This table allows network administrators and DevOps engineers to query VPC-specific details, including its ID, state, CIDR block, and whether it is the default VPC. Users can utilize this table to gather insights on VPCs, such as their networking configuration, security settings, and associated resources. The schema outlines the various attributes of the VPC, including the VPC ID, state, CIDR block, default VPC status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc` table, you can use the `.inspect aws_vpc` command in Steampipe.

### Key columns:

- `vpc_id`: The unique identifier for the VPC. This can be used to join with other tables that reference the VPC ID, such as the `aws_subnet` table.
- `state`: The current state of the VPC (pending or available). This can be useful to filter VPCs based on their availability for use.
- `cidr_block`: The IPv4 network range for the VPC, in CIDR notation. This can be used to understand the networking configuration of the VPC.

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