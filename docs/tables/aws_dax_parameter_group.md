---
title: "Table: aws_dax_parameter_group - Query AWS DAX Parameter Groups using SQL"
description: "Allows users to query AWS DynamoDB Accelerator (DAX) Parameter Groups, providing details such as parameter group name, ARN, description, and parameter settings."
---

# Table: aws_dax_parameter_group - Query AWS DAX Parameter Groups using SQL

The `aws_dax_parameter_group` table in Steampipe provides information about Parameter Groups within AWS DynamoDB Accelerator (DAX). This table allows DevOps engineers to query Parameter Group-specific details, including the group name, ARN, description, and parameter settings. Users can utilize this table to gather insights on Parameter Groups, such as their configurations, associated parameters, and more. The schema outlines the various attributes of the DAX Parameter Group, including the parameter group name, ARN, description, and associated parameters.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dax_parameter_group` table, you can use the `.inspect aws_dax_parameter_group` command in Steampipe.

### Key columns:

- `name`: The name of the DAX Parameter Group. It is a unique identifier and can be used to join this table with other tables that contain DAX Parameter Group information.
- `arn`: The Amazon Resource Name (ARN) of the DAX Parameter Group. This unique identifier can be used to join this table with other AWS resource tables.
- `description`: The description of the DAX Parameter Group. This can provide context and additional information when joining with other tables.

## Examples

### Basic info

```sql
select
  parameter_group_name,
  description,
  region
from
  aws_dax_parameter_group;
```

### Get cluster details associated with the parameter group

```sql
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