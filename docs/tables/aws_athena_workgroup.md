---
title: "Steampipe Table: aws_athena_workgroup - Query AWS Athena Workgroup using SQL"
description: "Allows users to query AWS Athena Workgroup details such as workgroup name, state, description, creation time, and more."
folder: "Athena"
---

# Table: aws_athena_workgroup - Query AWS Athena Workgroup using SQL

An AWS Athena Workgroup is a resource that acts as a primary server for running queries. It provides a means of managing query execution across multiple users and teams within an organization. This allows for better control over costs, performance, and security when querying data with Athena.

## Table Usage Guide

The `aws_athena_workgroup` table in Steampipe provides you with information about workgroups within AWS Athena. This table allows you as a DevOps engineer to query workgroup-specific details, including workgroup name, state, description, creation time, and more. You can utilize this table to gather insights on workgroups, such as workgroup configurations, encryption configurations, and enforcement settings. The schema outlines the various attributes of the Athena workgroup for you, including the workgroup ARN, state, tags, and configuration details.

## Examples

### List all workgroups with basic information
Explore the various workgroups within your AWS Athena service to gain insights into their basic details such as name, description, and creation time. This can be useful for understanding your workgroup configuration and identifying any potential areas for optimization or reorganization.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which workgroups are utilizing a specific version of the Athena engine. This is useful for assessing upgrade needs or understanding the distribution of engine versions across your workgroups.

```sql+postgres
select 
  name, 
  description 
from 
  aws_athena_workgroup 
where 
  effective_engine_version = 'Athena engine version 3';
```

```sql+sqlite
select 
  name, 
  description 
from 
  aws_athena_workgroup 
where 
  effective_engine_version = 'Athena engine version 3';
```

### Count workgroups in each region
Assess the distribution of workgroups across different regions to understand workload allocation and capacity planning. This can assist in identifying regions that may be under or over-utilized.

```sql+postgres
select 
  region, 
  count(*) 
from 
  aws_athena_workgroup 
group by 
  region;
```

```sql+sqlite
select 
  region, 
  count(*) 
from 
  aws_athena_workgroup 
group by 
  region;
```

### List disabled workgroups
Determine the areas in which workgroups are inactive, providing insights into resource usage and potential areas for optimization or re-allocation.

```sql+postgres
select 
  name, 
  description, 
  creation_time
from 
  aws_athena_workgroup 
where
  state = 'DISABLED';
```

```sql+sqlite
select 
  name, 
  description, 
  creation_time
from 
  aws_athena_workgroup 
where
  state = 'DISABLED';
```