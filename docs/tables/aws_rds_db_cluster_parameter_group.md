---
title: "Steampipe Table: aws_rds_db_cluster_parameter_group - Query AWS RDS DB Cluster Parameter Groups using SQL"
description: "Allows users to query AWS RDS DB Cluster Parameter Groups, providing detailed information about each parameter group's configuration, including its name, family, description, and ARN. This table can be used to identify unused or misconfigured parameter groups and to ensure they comply with security and operational best practices."
folder: "Resource Access Manager"
---

# Table: aws_rds_db_cluster_parameter_group - Query AWS RDS DB Cluster Parameter Groups using SQL

The AWS RDS DB Cluster Parameter Group is a component of Amazon RDS that contains database engine configuration values that are applied to a cluster. These groups enable the configuration of database settings at the cluster level, influencing the behavior of all the databases within the cluster. It allows for the fine-tuning of database instances for optimal performance and efficiency.

## Table Usage Guide

The `aws_rds_db_cluster_parameter_group` table in Steampipe provides you with information about DB Cluster Parameter Groups within Amazon RDS (Relational Database Service). This table allows you, as a DevOps engineer, DBA, or security professional, to query parameter group-specific details, including its name, family, description, and ARN. You can utilize this table to gather insights on parameter groups, such as identifying unused or misconfigured parameter groups and ensuring they comply with security and operational best practices. The schema outlines the various attributes of the DB Cluster Parameter Group for you, including the parameter group name, family, description, ARN, and associated tags.

## Examples

### List of DB cluster parameter group with corresponding parameter group family
Determine the areas in which specific database cluster parameter groups are associated with their corresponding parameter group families. This can be useful in understanding the configuration and settings of your database clusters in AWS RDS, aiding in optimization and troubleshooting.

```sql+postgres
select
  name,
  description,
  db_parameter_group_family
from
  aws_rds_db_cluster_parameter_group;
```

```sql+sqlite
select
  name,
  description,
  db_parameter_group_family
from
  aws_rds_db_cluster_parameter_group;
```


### Parameters info of each parameter group
This query is useful for understanding the configurations and settings of your database parameter groups. It allows you to examine the details of each parameter, including its value, the range of allowed values, and its modifiability, which can help in optimizing the database performance and security.

```sql+postgres
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

```sql+sqlite
select
  name,
  db_parameter_group_family,
  json_extract(pg.value, '$.ParameterName') as parameter_name,
  json_extract(pg.value, '$.ParameterValue') as parameter_value,
  json_extract(pg.value, '$.AllowedValues') as allowed_values,
  json_extract(pg.value, '$.ApplyType') as apply_type,
  json_extract(pg.value, '$.IsModifiable') as is_modifiable,
  json_extract(pg.value, '$.DataType') as data_type,
  json_extract(pg.value, '$.Description') as description,
  json_extract(pg.value, '$.MinimumEngineVersion') as minimum_engine_version
from
  aws_rds_db_cluster_parameter_group,
  json_each(parameters) as pg;
```