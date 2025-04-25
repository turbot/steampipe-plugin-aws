---
title: "Steampipe Table: aws_rds_db_option_group - Query AWS RDS DB Option Groups using SQL"
description: "Allows users to query AWS RDS DB Option Groups and provides information about the option groups within Amazon Relational Database Service (RDS)."
folder: "RDS"
---

# Table: aws_rds_db_option_group - Query AWS RDS DB Option Groups using SQL

The AWS RDS DB Option Groups service allows you to manage and configure additional features for your Amazon RDS databases. These groups are used to specify and manage the configuration options that are available for your DB instances. By using SQL, you can query these option groups to easily manage and organize your database configurations.

## Table Usage Guide

The `aws_rds_db_option_group` table in Steampipe provides you with information about the option groups within Amazon Relational Database Service (RDS). This table allows you, as a database administrator or developer, to query option group-specific details, including the options and parameters associated with the group, the engine name, and the major engine version. You can utilize this table to gather insights on option groups, such as identifying the configurations of specific databases, verifying the parameters of option groups, and more. The schema outlines the various attributes of the RDS DB Option Group for you, including the name, ARN, description, and associated tags.

## Examples

### Basic parameter group info
Analyze the settings to understand the basic information about your AWS RDS database option groups, such as the engine used and the version. This can help in managing and optimizing your database configurations.

```sql+postgres
select
  name,
  description,
  engine_name,
  major_engine_version,
  vpc_id
from
  aws_rds_db_option_group;
```

```sql+sqlite
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
Discover the segments that can be applied to both VPC and non-VPC instances in AWS RDS. This can be beneficial in understanding the flexibility and adaptability of your database configurations.

```sql+postgres
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

```sql+sqlite
select
  name,
  description,
  engine_name,
  allows_vpc_and_non_vpc_instance_memberships
from
  aws_rds_db_option_group
where
  allows_vpc_and_non_vpc_instance_memberships = 1;
```


### Option details of each option group
Explore the various options within each option group in an AWS RDS database. This can help in understanding the settings of each option, including its permanence, persistence, associated security group memberships, and port, aiding in effective database management and security.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(option.value, '$.OptionName') as option_name,
  json_extract(option.value, '$.Permanent') as Permanent,
  json_extract(option.value, '$.Persistent') as Persistent,
  json_extract(option.value, '$.VpcSecurityGroupMemberships') as vpc_security_group_membership,
  json_extract(option.value, '$.Port') as Port
from
  aws_rds_db_option_group,
  json_each(options) as option;
```