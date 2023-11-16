---
title: "Table: aws_rds_db_subnet_group - Query AWS RDS DB Subnet Groups using SQL"
description: "Allows users to query AWS RDS DB Subnet Groups to retrieve information about each DB subnet group configured in an AWS account."
---

# Table: aws_rds_db_subnet_group - Query AWS RDS DB Subnet Groups using SQL

The `aws_rds_db_subnet_group` table in Steampipe provides information about DB subnet groups within Amazon Relational Database Service (RDS). This table allows database administrators, developers, and other technical professionals to query details about DB subnet groups, including the subnet group name, description, VPC ID, and associated subnets. Users can utilize this table to gather insights on DB subnet groups, such as subnet group configurations, associated VPCs, and more. The schema outlines the various attributes of the DB subnet group, including the subnet group status, ARN, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_subnet_group` table, you can use the `.inspect aws_rds_db_subnet_group` command in Steampipe.

### Key columns:

- `db_subnet_group_name`: The name of the DB subnet group. This is a unique key and can be used to join with other tables that reference DB subnet groups.
- `vpc_id`: The ID of the VPC that the DB subnet group belongs to. This can be used to join with other tables that reference VPCs.
- `arn`: The Amazon Resource Name (ARN) of the DB subnet group. This can be used to join with other tables that reference DB subnet groups by ARN.

## Examples

### DB subnet group basic info

```sql
select
  name,
  status,
  vpc_id
from
  aws_rds_db_subnet_group;
```


### Subnets info of each subnet in subnet group

```sql
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


### List of subnet group without application tag key

```sql
select
  name,
  tags
from
  aws_rds_db_subnet_group
where
  not tags :: JSONB ? 'application';
```
