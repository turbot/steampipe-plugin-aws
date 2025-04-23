---
title: "Steampipe Table: aws_simspaceweaver_simulation - Query AWS SimSpace Simulation using SQL"
description: "Allows users to query AWS SimSpace Simulation data. This table provides information about simulations within AWS SimSpace. Engineers can use it to query simulation-specific details, including simulation status, configuration, and associated metadata."
folder: "SimSpace Weaver"
---

# Table: aws_simspaceweaver_simulation - Query AWS SimSpace Simulation using SQL

The AWS SimSpace Simulation is a service within the Amazon Web Services infrastructure that enables users to create, run, and manage simulations at scale. It's part of the AWS SimSpace offering, which is designed to provide computational simulation as a service. This service eliminates the need for managing the underlying compute resources, allowing users to focus on analyzing results and optimizing designs.

## Table Usage Guide

The `aws_simspaceweaver_simulation` table in Steampipe provides you with information about simulations within AWS SimSpace. This table allows you to query simulation-specific details, including simulation status, configuration, and associated metadata. You can utilize this table to gather insights on simulations, such as simulation state, configuration details, and more. The schema outlines the various attributes of the simulation for you, including the simulation ID, creation time, status, and associated tags.

## Examples

### Basic info
Explore which AWS SimSpaceWeaver simulations have been created and their current status. This can help in identifying and addressing any potential issues with these simulations.

```sql+postgres
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

```sql+sqlite
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
Identify instances where simulations have been running for more than 30 days. This can be useful for managing resources and ensuring simulations are not unnecessarily consuming resources.

```sql+postgres
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

```sql+sqlite
select
  name,
  arn,
  creation_time,
  status
from
  aws_simspaceweaver_simulation
where
  creation_time >= datetime('now','-30 day');
```

### List failed simulations
Identify instances where simulations have been unsuccessful to assess potential issues or errors within the AWS SimSpaceWeaver system.

```sql+postgres
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

```sql+sqlite
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
Analyze the settings to understand the logging configurations of your AWS SimSpaceWeaver simulations. This can be useful to ensure that your logging configurations are correctly set up and to troubleshoot any issues that may arise during your simulations.

```sql+postgres
select
  name,
  arn,
  jsonb_pretty(d)
from
  aws_simspaceweaver_simulation,
  jsonb_array_elements(logging_configuration -> 'Destinations') as d;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(d.value, '$') as d
from
  aws_simspaceweaver_simulation,
  json_each(logging_configuration, '$.Destinations') as d;
```

### Get S3 bucket details of simulations
Determine the areas in which simulations and S3 bucket details intersect. This can be useful for gaining insights into simulation configurations and assessing how they align with your S3 bucket settings, particularly in terms of versioning, public access, and access control lists.

```sql+postgres
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

```sql+sqlite
select
  s.name,
  s.arn,
  json_extract(s.schema_s3_location, '$.BucketName') as bucket_name,
  json_extract(s.schema_s3_location, '$.ObjectKey') as object_key,
  b.versioning_enabled,
  b.block_public_acls,
  b.acl
from
  aws_simspaceweaver_simulation as s,
  aws_s3_bucket as b
where
  json_extract(s.schema_s3_location, '$.BucketName') = b.name;
```