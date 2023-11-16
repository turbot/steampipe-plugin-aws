---
title: "Table: aws_simspaceweaver_simulation - Query AWS SimSpace Simulation using SQL"
description: "Allows users to query AWS SimSpace Simulation data. This table provides information about simulations within AWS SimSpace. Engineers can use it to query simulation-specific details, including simulation status, configuration, and associated metadata."
---

# Table: aws_simspaceweaver_simulation - Query AWS SimSpace Simulation using SQL

The `aws_simspaceweaver_simulation` table in Steampipe provides information about simulations within AWS SimSpace. This table allows engineers to query simulation-specific details, including simulation status, configuration, and associated metadata. Users can utilize this table to gather insights on simulations, such as simulation state, configuration details, and more. The schema outlines the various attributes of the simulation, including the simulation ID, creation time, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_simspaceweaver_simulation` table, you can use the `.inspect aws_simspaceweaver_simulation` command in Steampipe.

### Key columns:

- `simulation_id`: This is the unique identifier for each simulation. It can be used to join this table with other tables to gather more detailed information about each simulation.
- `status`: This column provides the current status of the simulation. It is useful for monitoring the progress and state of simulations.
- `tags`: This column stores metadata that is assigned to the simulation. It can be used to filter or categorize simulations based on user-defined criteria.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  status,
  execution_id,
  schema_error
from
  aws_simspaceweaver_simulation;
```

### List simulations older than 30 days

```sql
select
  name,
  arn,
  creation_time,
  status
from
  aws_simspaceweaver_simulation
where
  creation_time >= now() - interval '30' day;
```

### List failed simulations

```sql
select
  name,
  arn,
  creation_time,
  status
from
  aws_simspaceweaver_simulation
where
  status = 'FAILED';
```

### Get logging configurations of simulations

```sql
select
  name,
  arn,
  jsonb_pretty(d)
from
  aws_simspaceweaver_simulation,
  jsonb_array_elements(logging_configuration -> 'Destinations') as d;
```

### Get S3 bucket details of simulations

```sql
select
  s.name,
  s.arn,
  s.schema_s3_location ->> 'BucketName' as bucket_name,
  s.schema_s3_location ->> 'ObjectKey' as object_key,
  b.versioning_enabled,
  b.block_public_acls,
  b.acl
from
  aws_simspaceweaver_simulation as s,
  aws_s3_bucket as b
where
  s.schema_s3_location ->> 'BucketName' = b.name;
```
