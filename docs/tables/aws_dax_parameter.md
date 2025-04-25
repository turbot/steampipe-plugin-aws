---
title: "Steampipe Table: aws_dax_parameter - Query AWS DAX Parameter Groups using SQL"
description: "Allows users to query AWS DAX Parameter Groups to retrieve information about their configuration settings."
folder: "DAX"
---

# Table: aws_dax_parameter - Query AWS DAX Parameter Groups using SQL

AWS DAX Parameter Groups are a collection of parameters that you apply to all of the nodes in a DAX cluster. These groups make it easier to manage clusters by enabling you to customize their behavior without having to individually modify each node. They are particularly useful when you want to set consistent parameters across a large number of nodes.

## Table Usage Guide

The `aws_dax_parameter` table in Steampipe provides you with information about Parameter Groups within AWS DynamoDB Accelerator (DAX). This table allows you, as a DevOps engineer, to query parameter group-specific details, including parameter names, types, values, and whether they are modifiable or not. You can utilize this table to gather insights on parameter groups, such as understanding the configurations that control the behavior of your DAX clusters, and to verify if the parameters are set as per your requirements. The schema outlines the various attributes of the DAX parameter group for you, including the parameter name, value, source, data type, and whether it's modifiable or not.

## Examples

### Basic info
Explore which parameters are in use within your AWS DAX settings to understand their values and types. This information can assist in assessing the configuration for optimal performance and security.

```sql+postgres
select
  parameter_name,
  parameter_group_name,
  parameter_value,
  data_type,
  parameter_type
from
  aws_dax_parameter;
```

```sql+sqlite
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
Identify the distribution of parameters across different parameter groups and regions. This can help you understand how parameters are organized in your AWS DAX environment, which is useful for managing and optimizing your configurations.

```sql+postgres
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

```sql+sqlite
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
Identify instances where parameters can be modified in your AWS DAX setup. This is useful to understand which aspects of your configuration can be adjusted to optimize performance.

```sql+postgres
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

```sql+sqlite
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
Identify the parameters in your AWS DAX that are not user-defined. This can help ensure that system-defined settings are not inadvertently altered, maintaining system stability and performance.

```sql+postgres
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

  ```sql+sqlite
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
  source != 'user';
```