---
title: "Table: aws_rds_db_option_group - Query AWS RDS DB Option Groups using SQL"
description: "Allows users to query AWS RDS DB Option Groups and provides information about the option groups within Amazon Relational Database Service (RDS)."
---

# Table: aws_rds_db_option_group - Query AWS RDS DB Option Groups using SQL

The `aws_rds_db_option_group` table in Steampipe provides information about the option groups within Amazon Relational Database Service (RDS). This table allows database administrators and developers to query option group-specific details, including the options and parameters associated with the group, the engine name, and the major engine version. Users can utilize this table to gather insights on option groups, such as identifying the configurations of specific databases, verifying the parameters of option groups, and more. The schema outlines the various attributes of the RDS DB Option Group, including the name, ARN, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_option_group` table, you can use the `.inspect aws_rds_db_option_group` command in Steampipe.

### Key columns:

- `name`: The name of the option group. This is the primary identifier of the option group and can be used to join with other tables where option group name is required.
- `arn`: The Amazon Resource Name (ARN) of the option group. ARNs are unique identifiers for AWS resources and can be used for more complex queries involving multiple AWS services.
- `engine_name`: The name of the engine that this option group can be applied to. This can be useful for filtering option groups based on specific database engines.

## Examples

### Basic parameter group info

```sql
select
  name,
  description,
  engine_name,
  major_engine_version,
  vpc_id
from
  aws_rds_db_option_group;
```


### List of option groups which can be applied to both VPC and non-VPC instances

```sql
select
  name,
  description,
  engine_name,
  allows_vpc_and_non_vpc_instance_memberships
from
  aws_rds_db_option_group
where
  allows_vpc_and_non_vpc_instance_memberships;
```


### Option details of each option group

```sql
select
  name,
  option ->> 'OptionName' as option_name,
  option -> 'Permanent' as Permanent,
  option -> 'Persistent' as Persistent,
  option -> 'VpcSecurityGroupMemberships' as vpc_security_group_membership,
  option -> 'Port' as Port
from
  aws_rds_db_option_group
  cross join jsonb_array_elements(options) as option;
```