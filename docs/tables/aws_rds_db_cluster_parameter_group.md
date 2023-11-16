---
title: "Table: aws_rds_db_cluster_parameter_group - Query AWS RDS DB Cluster Parameter Groups using SQL"
description: "Allows users to query AWS RDS DB Cluster Parameter Groups, providing detailed information about each parameter group's configuration, including its name, family, description, and ARN. This table can be used to identify unused or misconfigured parameter groups and to ensure they comply with security and operational best practices."
---

# Table: aws_rds_db_cluster_parameter_group - Query AWS RDS DB Cluster Parameter Groups using SQL

The `aws_rds_db_cluster_parameter_group` table in Steampipe provides information about DB Cluster Parameter Groups within Amazon RDS (Relational Database Service). This table allows DevOps engineers, DBAs, and security professionals to query parameter group-specific details, including its name, family, description, and ARN. Users can utilize this table to gather insights on parameter groups, such as identifying unused or misconfigured parameter groups and ensuring they comply with security and operational best practices. The schema outlines the various attributes of the DB Cluster Parameter Group, including the parameter group name, family, description, ARN, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_cluster_parameter_group` table, you can use the `.inspect aws_rds_db_cluster_parameter_group` command in Steampipe.

### Key columns:

- `name`: The name of the DB Cluster Parameter Group. This column can be used to join with other tables that contain DB Cluster Parameter Group names.
- `arn`: The Amazon Resource Number (ARN) of the DB Cluster Parameter Group. This column can be used to join with any table that contains AWS ARNs.
- `family`: The DB parameter group family name. This column can be used to join with other tables that contain parameter group family names.

## Examples

### List of DB cluster parameter group with corresponding parameter group family

```sql
select
  name,
  description,
  db_parameter_group_family
from
  aws_rds_db_cluster_parameter_group;
```


### Parameters info of each parameter group

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
  aws_rds_db_cluster_parameter_group
  cross join jsonb_array_elements(parameters) as pg;
```
