---
title: "Steampipe Table: aws_glue_ml_transform - Query AWS Glue ML Transforms using SQL"
description: "Allows users to query AWS Glue ML Transforms, including transform details, parameters, and status information."
folder: "Glue"
---

# Table: aws_glue_ml_transform - Query AWS Glue ML Transforms using SQL

AWS Glue ML Transforms are machine learning transforms that use machine learning to learn the details of the transformation to be performed by learning from examples provided by humans. These transformations are then saved by Glue and can be used to process data.

## Table Usage Guide

The `aws_glue_ml_transform` table in Steampipe provides you with information about ML Transforms within AWS Glue. This table allows you, as a DevOps engineer, data engineer, or ML engineer, to query transform details including configuration parameters, status, performance metrics, and resource allocation. You can utilize this table to gather insights on ML transforms, such as transform status, worker configuration, input/output schemas, evaluation metrics, and associated tags.

**Important notes:**
- This table supports optional quals. Queries with optional quals are optimised to use AWS Glue filters. Optional quals are supported for the following columns:
  - `name`
  - `status`
  - `transform_type`
  - `glue_version`
  - `created_on`
  - `last_modified_on`

## Examples

### Basic info
Analyze the ML transforms to understand their configuration, status, and resource allocation. This is useful for monitoring transform health and performance.

```sql+postgres
select
  name,
  transform_id,
  status,
  glue_version,
  max_capacity,
  number_of_workers,
  worker_type,
  timeout,
  created_on,
  last_modified_on
from
  aws_glue_ml_transform;
```

```sql+sqlite
select
  name,
  transform_id,
  status,
  glue_version,
  max_capacity,
  number_of_workers,
  worker_type,
  timeout,
  created_on,
  last_modified_on
from
  aws_glue_ml_transform;
```

### List ML transforms that are in ready state
Identify ML transforms that are ready to be used for data processing. This helps in understanding which transforms are available for immediate use.

```sql+postgres
select
  name,
  transform_id,
  status,
  description,
  label_count,
  created_on
from
  aws_glue_ml_transform
where
  status = 'READY';
```

```sql+sqlite
select
  name,
  transform_id,
  status,
  description,
  label_count,
  created_on
from
  aws_glue_ml_transform
where
  status = 'READY';
```

### Check the resource allocation of ML transforms
Review the resource allocation for ML transforms to understand compute capacity and worker configuration.

```sql+postgres
select
  name,
  transform_id,
  max_capacity,
  number_of_workers,
  worker_type,
  max_retries,
  timeout,
  case
    when max_capacity is not null then 'MaxCapacity mode'
    when number_of_workers is not null then 'Worker mode'
    else 'Default mode'
  end as allocation_mode
from
  aws_glue_ml_transform;
```

```sql+sqlite
select
  name,
  transform_id,
  max_capacity,
  number_of_workers,
  worker_type,
  max_retries,
  timeout,
  case
    when max_capacity is not null then 'MaxCapacity mode'
    when number_of_workers is not null then 'Worker mode'
    else 'Default mode'
  end as allocation_mode
from
  aws_glue_ml_transform;
```

### Find transforms with evaluation metrics
Identify ML transforms that have evaluation metrics available, which indicates they have been trained and evaluated.

```sql+postgres
select
  name,
  transform_id,
  status,
  evaluation_metrics,
  label_count
from
  aws_glue_ml_transform
where
  evaluation_metrics is not null;
```

```sql+sqlite
select
  name,
  transform_id,
  status,
  evaluation_metrics,
  label_count
from
  aws_glue_ml_transform
where
  evaluation_metrics is not null;
```

### Review the input and output schemas for ML transforms
Review the input and output schemas for ML transforms to understand data structure requirements.

```sql+postgres
select
  name,
  transform_id,
  input_record_tables,
  schema,
  parameters
from
  aws_glue_ml_transform;
```

```sql+sqlite
select
  name,
  transform_id,
  input_record_tables,
  schema,
  parameters
from
  aws_glue_ml_transform;
```

### List transforms by worker type
Group ML transforms by their worker type to understand resource allocation patterns.

```sql+postgres
select
  worker_type,
  count(*) as transform_count,
  avg(max_capacity) as avg_max_capacity,
  avg(number_of_workers) as avg_workers
from
  aws_glue_ml_transform
group by
  worker_type
order by
  transform_count desc;
```

```sql+sqlite
select
  worker_type,
  count(*) as transform_count,
  avg(max_capacity) as avg_max_capacity,
  avg(number_of_workers) as avg_workers
from
  aws_glue_ml_transform
group by
  worker_type
order by
  transform_count desc;
```

### Check transform encryption settings
Review the encryption settings for ML transforms to ensure data security compliance.

```sql+postgres
select
  name,
  transform_id,
  transform_encryption,
  role
from
  aws_glue_ml_transform
where
  transform_encryption is not null;
```

```sql+sqlite
select
  name,
  transform_id,
  transform_encryption,
  role
from
  aws_glue_ml_transform
where
  transform_encryption is not null;
```


### Filter transforms by status
Find transforms with a specific status using the optional qualifier.

```sql+postgres
select
  name,
  transform_id,
  status,
  created_on
from
  aws_glue_ml_transform
where
  status = 'READY';
```

```sql+sqlite
select
  name,
  transform_id,
  status,
  created_on
from
  aws_glue_ml_transform
where
  status = 'READY';
```

### Filter transforms by creation date
Find transforms created after a specific date using the optional qualifier with operators.

```sql+postgres
select
  name,
  transform_id,
  status,
  created_on
from
  aws_glue_ml_transform
where
  created_on >= '2024-01-01';
```

```sql+sqlite
select
  name,
  transform_id,
  status,
  created_on
from
  aws_glue_ml_transform
where
  created_on >= '2024-01-01';
```

### Filter transforms by last modification date
Find transforms modified within a specific time range using the optional qualifier with operators.

```sql+postgres
select
  name,
  transform_id,
  status,
  last_modified_on
from
  aws_glue_ml_transform
where
  last_modified_on >= '2024-01-01'
  and last_modified_on <= '2024-12-31';
```

```sql+sqlite
select
  name,
  transform_id,
  status,
  last_modified_on
from
  aws_glue_ml_transform
where
  last_modified_on >= '2024-01-01'
  and last_modified_on <= '2024-12-31';
```

### Filter transforms by Glue version
Find transforms compatible with a specific Glue version using the optional qualifier.

```sql+postgres
select
  name,
  transform_id,
  glue_version,
  status
from
  aws_glue_ml_transform
where
  glue_version = '3.0';
```

```sql+sqlite
select
  name,
  transform_id,
  glue_version,
  status
from
  aws_glue_ml_transform
where
  glue_version = '3.0';
```
