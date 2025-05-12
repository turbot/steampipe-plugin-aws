---
title: "Steampipe Table: aws_redshift_subnet_group - Query AWS Redshift Subnet Groups using SQL"
description: "Allows users to query AWS Redshift Subnet Groups and get detailed information about each subnet group, including its name, description, VPC ID, subnet IDs, and status."
folder: "Redshift"
---

# Table: aws_redshift_subnet_group - Query AWS Redshift Subnet Groups using SQL

The AWS Redshift Subnet Group is a collection of subnets that you may want to designate for your Amazon Redshift clusters in a Virtual Private Cloud (VPC). This resource allows you to specify a range of IP addresses in the VPC and a set of security groups to be associated with the Amazon Redshift cluster. The subnets must be in the same VPC and you can specify up to 20 subnets in a subnet group.

## Table Usage Guide

The `aws_redshift_subnet_group` table in Steampipe provides you with information about subnet groups within AWS Redshift. This table allows you, as a DevOps engineer, to query subnet group-specific details, including their names, descriptions, VPC IDs, subnet IDs, and their status. You can utilize this table to gather insights on subnet groups, such as which subnet groups are available, their associated VPCs and subnets, and their current status. The schema outlines the various attributes of the subnet group for you, including the subnet group name, VPC ID, subnet IDs, and status.

## Examples

### Basic info
Explore the status of your Redshift subnet groups within your Virtual Private Cloud (VPC) to understand their current condition. This can help in assessing the health of your cloud resources, enabling timely interventions and maintenance.

```sql+postgres
select
  cluster_subnet_group_name,
  description,
  subnet_group_status,
  vpc_id
from
  aws_redshift_subnet_group;
```

```sql+sqlite
select
  cluster_subnet_group_name,
  description,
  subnet_group_status,
  vpc_id
from
  aws_redshift_subnet_group;
```


### Get the subnet info for each subnet group
Determine the availability and status of each subnet within a group. This query is useful for understanding the health and configuration of your network infrastructure, enabling proactive management and troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  cluster_subnet_group_name,
  json_extract(subnet.value, '$.SubnetAvailabilityZone.Name') as subnet_availability_zone,
  json_extract(subnet.value, '$.SubnetAvailabilityZone.SupportedPlatforms') as supported_platforms,
  json_extract(subnet.value, '$.SubnetIdentifier') as subnet_identifier,
  json_extract(subnet.value, '$.SubnetStatus') as subnet_status
from
  aws_redshift_subnet_group,
  json_each(subnets) as subnet;
```


### List subnet groups without the application tag key
Discover the segments that are missing the 'application' tag key within your subnet groups. This can be useful for identifying areas in your AWS Redshift environment that may not be correctly tagged, ensuring proper resource tracking and management.

```sql+postgres
select
  cluster_subnet_group_name,
  tags
from
  aws_redshift_subnet_group
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  cluster_subnet_group_name,
  tags
from
  aws_redshift_subnet_group
where
  json_extract(tags, '$.application') is null;
```