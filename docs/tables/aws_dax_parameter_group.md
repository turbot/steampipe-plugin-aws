---
title: "Steampipe Table: aws_dax_parameter_group - Query AWS DAX Parameter Groups using SQL"
description: "Allows users to query AWS DynamoDB Accelerator (DAX) Parameter Groups, providing details such as parameter group name, ARN, description, and parameter settings."
folder: "DAX"
---

# Table: aws_dax_parameter_group - Query AWS DAX Parameter Groups using SQL

The AWS DAX Parameter Group is a resource that provides a container for database engine parameter values that can be applied to one or more DAX clusters. These parameters act as a means to manage the behavior of the DAX instances within the cluster. In essence, it allows you to establish configurations and settings for your DAX databases, providing customization and control over the DAX environment.

## Table Usage Guide

The `aws_dax_parameter_group` table in Steampipe provides you with information about Parameter Groups within AWS DynamoDB Accelerator (DAX). This table enables you, as a DevOps engineer, to query Parameter Group-specific details, including the group name, ARN, description, and parameter settings. You can utilize this table to gather insights on Parameter Groups, such as their configurations, associated parameters, and more. The schema outlines the various attributes of the DAX Parameter Group for you, including the parameter group name, ARN, description, and associated parameters.

## Examples

### Basic info
Gain insights into the regions and their associated descriptions within your AWS DAX parameter groups. This can be useful to understand the geographical distribution and purpose of your parameter groups.

```sql+postgres
select
  parameter_group_name,
  description,
  region
from
  aws_dax_parameter_group;
```

```sql+sqlite
select
  parameter_group_name,
  description,
  region
from
  aws_dax_parameter_group;
```

### Get cluster details associated with the parameter group
Discover the segments that are linked to a specific parameter group in your DAX clusters. This is useful for assessing the configuration of your clusters and understanding their current state.

```sql+postgres
select
  p.parameter_group_name,
  c.cluster_name,
  c.node_type,
  c.status
from
  aws_dax_parameter_group as p,
  aws_dax_cluster as c
where
  c.parameter_group ->> 'ParameterGroupName' = p.parameter_group_name;
```

```sql+sqlite
select
  p.parameter_group_name,
  c.cluster_name,
  c.node_type,
  c.status
from
  aws_dax_parameter_group as p,
  aws_dax_cluster as c
where
  json_extract(c.parameter_group, '$.ParameterGroupName') = p.parameter_group_name;
```