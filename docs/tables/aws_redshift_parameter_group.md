---
title: "Steampipe Table: aws_redshift_parameter_group - Query Amazon Redshift Parameter Groups using SQL"
description: "Allows users to query Amazon Redshift Parameter Groups to obtain detailed information about the configuration parameters and settings for Redshift clusters. This can be useful for managing and optimizing the performance of Redshift databases."
folder: "Resource Access Manager"
---

# Table: aws_redshift_parameter_group - Query Amazon Redshift Parameter Groups using SQL

Title: Amazon Redshift Parameter Groups

Description: Amazon Redshift Parameter Groups provide customizable parameters that control the behavior of all the databases that a cluster contains. They enable more granular control over the database setup by allowing the modification of system parameters. This flexibility can enhance database security, performance, and manageability.

## Table Usage Guide

The `aws_redshift_parameter_group` table in Steampipe provides you with information about Parameter Groups within Amazon Redshift. This table allows you, as a DevOps engineer, to query Parameter Group-specific details, including parameter names, values, and whether they are modifiable. You can utilize this table to gather insights on Parameter Groups, such as understanding the configuration of Redshift clusters, optimizing database performance, and ensuring adherence to best practices. The schema outlines the various attributes of the Parameter Group for you, including the parameter group name, parameter apply status, and associated tags.

## Examples

### Basic info
Analyze the settings of your AWS Redshift parameter groups to understand their configurations and relationships. This can help in optimizing the performance and security of your Redshift databases.

```sql+postgres
select
  name,
  description,
  family
from
  aws_redshift_parameter_group;
```

```sql+sqlite
select
  name,
  description,
  family
from
  aws_redshift_parameter_group;
```


### List parameter groups that have the require_ssl parameter set to false
Determine the areas in which parameter groups are not requiring SSL. This is useful for identifying potential security vulnerabilities in your AWS Redshift Parameter Groups.

```sql+postgres
select
  name,
  p ->> 'ParameterName' as parameter_name,
  p ->> 'ParameterValue' as parameter_value,
  p ->> 'Description' as description,
  p ->> 'Source' as source,
  p ->> 'DataType' as data_type,
  p ->> 'ApplyType' as apply_type,
  p ->> 'IsModifiable' as is_modifiable,
  p ->> 'AllowedValues' as allowed_values,
  p ->> 'MinimumEngineVersion' as minimum_engine_version
from
  aws_redshift_parameter_group,
  jsonb_array_elements(parameters) as p
where
  p ->> 'ParameterName' = 'require_ssl'
  and p ->> 'ParameterValue' = 'false';
```

```sql+sqlite
select
  name,
  json_extract(p.value, '$.ParameterName') as parameter_name,
  json_extract(p.value, '$.ParameterValue') as parameter_value,
  json_extract(p.value, '$.Description') as description,
  json_extract(p.value, '$.Source') as source,
  json_extract(p.value, '$.DataType') as data_type,
  json_extract(p.value, '$.ApplyType') as apply_type,
  json_extract(p.value, '$.IsModifiable') as is_modifiable,
  json_extract(p.value, '$.AllowedValues') as allowed_values,
  json_extract(p.value, '$.MinimumEngineVersion') as minimum_engine_version
from
  aws_redshift_parameter_group,
  json_each(parameters) as p
where
  json_extract(p, '$.ParameterName') = 'require_ssl'
  and json_extract(p.value, '$.ParameterValue') = 'false';
```