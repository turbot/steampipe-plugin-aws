---
title: "Table: aws_athena_workgroup - Query AWS Athena Workgroup using SQL"
description: "Allows users to query AWS Athena Workgroup details such as workgroup name, state, description, creation time, and more."
---

# Table: aws_athena_workgroup - Query AWS Athena Workgroup using SQL

The `aws_athena_workgroup` table in Steampipe provides information about workgroups within AWS Athena. This table allows DevOps engineers to query workgroup-specific details, including workgroup name, state, description, creation time, and more. Users can utilize this table to gather insights on workgroups, such as workgroup configurations, encryption configurations, and enforcement settings. The schema outlines the various attributes of the Athena workgroup, including the workgroup ARN, state, tags, and configuration details.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_athena_workgroup` table, you can use the `.inspect aws_athena_workgroup` command in Steampipe.

**Key columns**:

- `name`: The unique name of the workgroup. This column can be used to join with other tables that require a workgroup name.
- `arn`: The Amazon Resource Name (ARN) of the workgroup. This column can be used to join with other tables that require a workgroup ARN.
- `state`: The state of the workgroup (ENABLED or DISABLED). This column can be used to filter workgroups based on their state.


## Examples

### List all workgroups with basic information

```sql
select 
  name, 
  description, 
  effective_engine_version, 
  output_location, 
  creation_time 
from 
  aws_athena_workgroup 
order by 
  creation_time;
```

### List all workgroups using engine 3

```sql
select 
  name, 
  description 
from 
  aws_athena_workgroup 
where 
  effective_engine_version = 'Athena engine version 3';
```

### Count workgroups in each region

```sql
select 
  region, 
  count(*) 
from 
  aws_athena_workgroup 
group by 
  region;
```

### List disabled workgroups

```sql
select 
  name, 
  description, 
  creation_time
from 
  aws_athena_workgroup 
where
  state = 'DISABLED';
```
