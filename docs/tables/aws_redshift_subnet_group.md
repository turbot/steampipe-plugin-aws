---
title: "Table: aws_redshift_subnet_group - Query AWS Redshift Subnet Groups using SQL"
description: "Allows users to query AWS Redshift Subnet Groups and get detailed information about each subnet group, including its name, description, VPC ID, subnet IDs, and status."
---

# Table: aws_redshift_subnet_group - Query AWS Redshift Subnet Groups using SQL

The `aws_redshift_subnet_group` table in Steampipe provides information about subnet groups within AWS Redshift. This table allows DevOps engineers to query subnet group-specific details, including their names, descriptions, VPC IDs, subnet IDs, and their status. Users can utilize this table to gather insights on subnet groups, such as which subnet groups are available, their associated VPCs and subnets, and their current status. The schema outlines the various attributes of the subnet group, including the subnet group name, VPC ID, subnet IDs, and status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_redshift_subnet_group` table, you can use the `.inspect aws_redshift_subnet_group` command in Steampipe.

### Key columns:

- `subnet_group_name`: The name of the subnet group. This can be used to join with other tables that contain subnet group names.
- `vpc_id`: The ID of the VPC the subnet group is associated with. This can be used to join with other tables that contain VPC IDs.
- `subnet_ids`: The IDs of the subnets in the subnet group. This can be used to join with other tables that contain subnet IDs.

## Examples

### Basic info

```sql
select
  cluster_subnet_group_name,
  description,
  subnet_group_status,
  vpc_id
from
  aws_redshift_subnet_group;
```


### Get the subnet info for each subnet group

```sql
select
  cluster_subnet_group_name,
  subnet -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  subnet -> 'SubnetAvailabilityZone' ->> 'SupportedPlatforms' as supported_platforms,
  subnet ->> 'SubnetIdentifier' as subnet_identifier,
  subnet ->> 'SubnetStatus' as subnet_status
from
  aws_redshift_subnet_group,
  jsonb_array_elements(subnets) as subnet;
```


### List subnet groups without the application tag key

```sql
select
  cluster_subnet_group_name,
  tags
from
  aws_redshift_subnet_group
where
  not tags :: JSONB ? 'application';
```
