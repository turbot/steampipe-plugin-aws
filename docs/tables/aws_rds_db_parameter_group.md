---
title: "Steampipe Table: aws_rds_db_parameter_group - Query AWS RDS DB Parameter Groups using SQL"
description: "Allows users to query AWS RDS DB Parameter Groups, providing information about the configurations that control the behavior of the databases that they are associated with."
folder: "Resource Access Manager"
---

# Table: aws_rds_db_parameter_group - Query AWS RDS DB Parameter Groups using SQL

The AWS RDS DB Parameter Group is a feature of Amazon Relational Database Service (RDS) that allows you to manage database engine configuration settings. These groups act as containers for engine parameter values that can be applied to one or more DB instances. The parameters in a DB parameter group enable you to control the runtime settings of a DB instance.

## Table Usage Guide

The `aws_rds_db_parameter_group` table in Steampipe provides you with information about DB Parameter Groups within AWS Relational Database Service (RDS). This table allows you, as a DevOps engineer, to query parameter group-specific details, including associated databases, parameter settings, and associated metadata. You can utilize this table to gather insights on parameter groups, such as understanding the configurations that control the behavior of the databases they are associated with, ensuring appropriate settings for optimal database performance, and more. The schema outlines the various attributes of the DB Parameter Group for you, including the parameter group name, family, description, and associated tags.

## Examples

### List of DB parameter group and corresponding parameter group family
Discover the correlation between database parameter groups and their corresponding families within the AWS RDS service. This is particularly useful when managing or optimizing the configuration of your relational database system.

```sql+postgres
select
  name,
  description,
  db_parameter_group_family
from
  aws_rds_db_parameter_group;
```

```sql+sqlite
select
  name,
  description,
  db_parameter_group_family
from
  aws_rds_db_parameter_group;
```


### Parameters info of each parameter group.
Determine the areas in which different parameters within a database parameter group can be modified, and understand the specific settings that apply to each. This can help to optimize database performance and security by ensuring parameters are correctly configured.

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
  aws_rds_db_parameter_group
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
  aws_rds_db_parameter_group,
  json_each(parameters) as pg;
```