---
title: "Table: aws_redshift_parameter_group - Query Amazon Redshift Parameter Groups using SQL"
description: "Allows users to query Amazon Redshift Parameter Groups to obtain detailed information about the configuration parameters and settings for Redshift clusters. This can be useful for managing and optimizing the performance of Redshift databases."
---

# Table: aws_redshift_parameter_group - Query Amazon Redshift Parameter Groups using SQL

The `aws_redshift_parameter_group` table in Steampipe provides information about Parameter Groups within Amazon Redshift. This table allows DevOps engineers to query Parameter Group-specific details, including parameter names, values, and whether they are modifiable. Users can utilize this table to gather insights on Parameter Groups, such as understanding the configuration of Redshift clusters, optimizing database performance, and ensuring adherence to best practices. The schema outlines the various attributes of the Parameter Group, including the parameter group name, parameter apply status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_redshift_parameter_group` table, you can use the `.inspect aws_redshift_parameter_group` command in Steampipe.

**Key columns**:

- `name`: The name of the parameter group. This is the primary key for the table and can be used to join with other tables that also contain Redshift parameter group names.
- `tags`: The metadata tags assigned to the parameter group. These can be used to filter and categorize parameter groups based on user-defined criteria.
- `parameter_group_family`: The family of the parameter group. This can be useful when joining with other tables that contain information about Redshift clusters, as it provides context about the type and compatibility of the parameter group.

## Examples

### Basic info

```sql
select
  name,
  description,
  family
from
  aws_redshift_parameter_group;
```


### List parameter groups that have the require_ssl parameter set to false

```sql
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
