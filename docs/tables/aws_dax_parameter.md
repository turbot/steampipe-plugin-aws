---
title: "Table: aws_dax_parameter - Query AWS DAX Parameter Groups using SQL"
description: "Allows users to query AWS DAX Parameter Groups to retrieve information about their configuration settings."
---

# Table: aws_dax_parameter - Query AWS DAX Parameter Groups using SQL

The `aws_dax_parameter` table in Steampipe provides information about Parameter Groups within AWS DynamoDB Accelerator (DAX). This table allows DevOps engineers to query parameter group-specific details, including parameter names, types, values, and whether they are modifiable. Users can utilize this table to gather insights on parameter groups, such as understanding the configurations that control the behavior of their DAX clusters, and to verify if the parameters are set as per their requirements. The schema outlines the various attributes of the DAX parameter group, including the parameter name, value, source, data type, and whether it's modifiable or not.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dax_parameter` table, you can use the `.inspect aws_dax_parameter` command in Steampipe.

Key columns:

- `parameter_group_name`: The name of the DAX parameter group. This is a key attribute to join this table with the `aws_dax_parameter_group` table.
- `parameter_name`: The name of the parameter. This provides insights on the specific configuration setting of the DAX parameter group.
- `is_modifiable`: Indicates whether the parameter can be modified. This is useful for understanding which parameters can be adjusted to tune the performance of the DAX clusters.

## Examples

### Basic info

```sql
select
  parameter_name,
  parameter_group_name,
  parameter_value,
  data_type,
  parameter_type
from
  aws_dax_parameter;
```

### Count parameters by parameter group

```sql
select
  parameter_group_name,
  region,
  count(parameter_name) as number_of_parameters
from
  aws_dax_parameter
group by
  parameter_group_name, 
  region;
```

### List modifiable parameters

```sql
select
  parameter_name,
  parameter_group_name,
  parameter_value,
  data_type,
  parameter_type,
  is_modifiable
from
  aws_dax_parameter
where
  is_modifiable = 'TRUE';
```

### List parameters that are not user defined

```sql
select
  parameter_name,
  change_type,
  parameter_group_name,
  parameter_value,
  data_type,
  parameter_type,
  source
from
  aws_dax_parameter
where
  source <> 'user';
  ```