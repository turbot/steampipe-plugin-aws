---
title: "Table: aws_dax_subnet_group - Query AWS DAX Subnet Group using SQL"
description: "Allows users to query AWS DAX Subnet Group details, such as the subnet group name, description, VPC ID, and the subnets in the group."
---

# Table: aws_dax_subnet_group - Query AWS DAX Subnet Group using SQL

The `aws_dax_subnet_group` table in Steampipe provides information about subnet groups within Amazon DynamoDB Accelerator (DAX). This table allows DevOps engineers to query subnet group-specific details, including the subnet group name, description, VPC ID, and the subnets in the group. Users can utilize this table to gather insights on subnet groups, such as their associated VPCs, subnet IDs, and more. The schema outlines the various attributes of the DAX subnet group, including the subnet group name, VPC ID, subnet ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dax_subnet_group` table, you can use the `.inspect aws_dax_subnet_group` command in Steampipe.

### Key columns:

- `subnet_group_name`: The name of the subnet group. This can be used to join this table with other tables that contain subnet group names.
- `vpc_id`: The ID of the VPC associated with the subnet group. This can be used to join this table with other tables that contain VPC IDs.
- `subnet_ids`: The IDs of the subnets in the subnet group. This can be used to join this table with other tables that contain subnet IDs.

## Examples

### Basic info

```sql
select
  subnet_group_name,
  description,
  vpc_id,
  subnets,
  region
from
  aws_dax_subnet_group;
```

### List VPC details for each subnet group

```sql
select
  subnet_group_name,
  v.vpc_id,
  v.arn as vpc_arn,
  v.cidr_block as vpc_cidr_block,
  v.state as vpc_state,
  v.is_default as is_default_vpc,
  v.region
from
  aws_dax_subnet_group g
join aws_vpc v
  on v.vpc_id = g.vpc_id;
```

### List subnet details for each subnet group

```sql
select
  subnet_group_name,
  g.vpc_id,
  vs.subnet_arn,
  vs.cidr_block as subnet_cidr_block,
  vs.state as subnet_state,
  vs.availability_zone as subnet_availability_zone,
  vs.region
from
  aws_dax_subnet_group g,
  jsonb_array_elements(subnets) s
join aws_vpc_subnet vs
  on vs.subnet_id = s ->> 'SubnetIdentifier';
```