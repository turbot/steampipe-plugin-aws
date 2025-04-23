---
title: "Steampipe Table: aws_dax_subnet_group - Query AWS DAX Subnet Group using SQL"
description: "Allows users to query AWS DAX Subnet Group details, such as the subnet group name, description, VPC ID, and the subnets in the group."
folder: "DAX"
---

# Table: aws_dax_subnet_group - Query AWS DAX Subnet Group using SQL

The AWS DAX Subnet Group is a resource in Amazon DynamoDB Accelerator (DAX) that allows you to specify a particular subnet group when you create a DAX cluster. A subnet group is a collection of subnets (typically private) that you can designate for your clusters running in a virtual private cloud (VPC). This allows you to configure network access to your DAX clusters.

## Table Usage Guide

The `aws_dax_subnet_group` table in Steampipe provides you with information about subnet groups within Amazon DynamoDB Accelerator (DAX). This table allows you, as a DevOps engineer, to query subnet group-specific details, including the subnet group name, description, VPC ID, and the subnets in the group. You can utilize this table to gather insights on subnet groups, such as their associated VPCs, subnet IDs, and more. The schema outlines the various attributes of the DAX subnet group for you, including the subnet group name, VPC ID, subnet ID, and associated tags.

## Examples

### Basic info
Explore which AWS DAX subnet groups are in use, gaining insights into their associated VPCs and regions. This can be useful for assessing your network's configuration and understanding its geographical distribution.

```sql+postgres
select
  subnet_group_name,
  description,
  vpc_id,
  subnets,
  region
from
  aws_dax_subnet_group;
```

```sql+sqlite
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
Determine the areas in which each subnet group is associated with specific VPC details. This can be useful in understanding the configuration and state of your network for better resource management and security.

```sql+postgres
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

```sql+sqlite
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
This query is useful for gaining insights into the specific details of each subnet group within a network. Using this information, one could optimize network structure, improve resource allocation, or enhance security measures.

```sql+postgres
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

```sql+sqlite
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
  json_each(subnets) s
join aws_vpc_subnet vs
  on vs.subnet_id = json_extract(s.value, '$.SubnetIdentifier');
```