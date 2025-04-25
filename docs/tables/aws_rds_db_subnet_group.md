---
title: "Steampipe Table: aws_rds_db_subnet_group - Query AWS RDS DB Subnet Groups using SQL"
description: "Allows users to query AWS RDS DB Subnet Groups to retrieve information about each DB subnet group configured in an AWS account."
folder: "RDS"
---

# Table: aws_rds_db_subnet_group - Query AWS RDS DB Subnet Groups using SQL

The AWS RDS DB Subnet Group is a collection of subnets that you can designate for your Amazon RDS DB instances in a Virtual Private Cloud (VPC). This grouping allows you to configure network settings for your RDS instances, enhancing data security and connectivity. It serves as a bridge between an Amazon RDS and the related network, enabling the database service to operate seamlessly within the VPC.

## Table Usage Guide

The `aws_rds_db_subnet_group` table in Steampipe provides you with information about DB subnet groups within Amazon Relational Database Service (RDS). This table allows you, as a database administrator, developer, or other technical professional, to query details about DB subnet groups, including the subnet group name, description, VPC ID, and associated subnets. You can utilize this table to gather insights on DB subnet groups, such as subnet group configurations, associated VPCs, and more. The schema outlines the various attributes of the DB subnet group for you, including the subnet group status, ARN, and associated tags.

## Examples

### DB subnet group basic info
Analyze the status of your database subnet groups within your Virtual Private Cloud (VPC) to understand their operational state. This can be beneficial for maintaining the health and performance of your AWS RDS instances.

```sql+postgres
select
  name,
  status,
  vpc_id
from
  aws_rds_db_subnet_group;
```

```sql+sqlite
select
  name,
  status,
  vpc_id
from
  aws_rds_db_subnet_group;
```

### Subnets info of each subnet in subnet group
Determine the status and location details of each subnet within a subnet group in your AWS RDS, to understand their availability and configuration. This information can be crucial for managing your database's network performance and security.

```sql+postgres
select
  name,
  subnet -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  subnet ->> 'SubnetIdentifier' as subnet_identifier,
  subnet -> 'SubnetOutpost' ->> 'Arn' as subnet_outpost,
  subnet ->> 'SubnetStatus' as subnet_status
from
  aws_rds_db_subnet_group
  cross join jsonb_array_elements(subnets) as subnet;
```

```sql+sqlite
select
  name,
  json_extract(subnet.value, '$.SubnetAvailabilityZone.Name') as subnet_availability_zone,
  json_extract(subnet.value, '$.SubnetIdentifier') as subnet_identifier,
  json_extract(subnet.value, '$.SubnetOutpost.Arn') as subnet_outpost,
  json_extract(subnet.value, '$.SubnetStatus') as subnet_status
from
  aws_rds_db_subnet_group,
  json_each(subnets) as subnet;
```

### List of subnet group without application tag key
Discover the segments that lack the 'application' tag key in your AWS RDS subnet groups. This can be useful in identifying potential areas for better resource tagging and management.

```sql+postgres
select
  name,
  tags
from
  aws_rds_db_subnet_group
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  name,
  tags
from
  aws_rds_db_subnet_group
where
  json_extract(tags, '$.application') is null;
```