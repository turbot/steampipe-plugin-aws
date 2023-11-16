---
title: "Table: aws_rds_db_parameter_group - Query AWS RDS DB Parameter Groups using SQL"
description: "Allows users to query AWS RDS DB Parameter Groups, providing information about the configurations that control the behavior of the databases that they are associated with."
---

# Table: aws_rds_db_parameter_group - Query AWS RDS DB Parameter Groups using SQL

The `aws_rds_db_parameter_group` table in Steampipe provides information about DB Parameter Groups within AWS Relational Database Service (RDS). This table allows DevOps engineers to query parameter group-specific details, including associated databases, parameter settings, and associated metadata. Users can utilize this table to gather insights on parameter groups, such as understanding the configurations that control the behavior of the databases they are associated with, ensuring appropriate settings for optimal database performance, and more. The schema outlines the various attributes of the DB Parameter Group, including the parameter group name, family, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_parameter_group` table, you can use the `.inspect aws_rds_db_parameter_group` command in Steampipe.

### Key columns:

- `name`: The name of the DB parameter group. This column can be used to join the table with other tables that contain parameter group names.
- `db_parameter_group_family`: The name of the DB parameter group family that the DB parameter group is compatible with. This column can be used to join the table with other tables that contain parameter group family names.
- `description`: The description of the DB parameter group. This column can be used to join the table with other tables that contain parameter group descriptions, providing context and additional information about the group.

## Examples

### List of DB parameter group and corresponding parameter group family

```sql
select
  name,
  description,
  db_parameter_group_family
from
  aws_rds_db_parameter_group;
```


### Parameters info of each parameter group.

```sql
select
  name,
  db_parameter_group_family,
  pg ->> 'ParameterName' as parameter_name,
  pg ->> 'ParameterValue' as parameter_value,
  pg ->> 'AllowedValues' as allowed_values,
  pg ->> 'ApplyType' as apply_type,
  pg ->> 'IsModifiable' as is_modifiable,
  pg ->> 'DataType' as data_type,
  pg ->> 'Description' as description,
  pg ->> 'MinimumEngineVersion' as minimum_engine_version
from
  aws_rds_db_parameter_group
  cross join jsonb_array_elements(parameters) as pg;
```