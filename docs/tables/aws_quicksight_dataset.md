---
title: "Steampipe Table: aws_quicksight_dataset - Query AWS QuickSight Datasets using SQL"
description: "Allows users to query AWS QuickSight Datasets, providing details about dataset configurations, data sources, and logical data structures within QuickSight."
---

# Table: aws_quicksight_dataset - Query AWS QuickSight Datasets using SQL

AWS QuickSight Dataset is a core component that defines how source data should be used in QuickSight. It includes configurations for data import, transformations, calculated fields, and relationships between different data sources, enabling users to create meaningful visualizations and analyses.

## Table Usage Guide

The `aws_quicksight_dataset` table in Steampipe provides you with information about datasets within AWS QuickSight. This table allows you, as a data analyst or administrator, to query dataset-specific details, including import modes, permissions, and logical data structures. You can utilize this table to gather insights on datasets, such as their creation time, last update time, and various configuration settings.

## Examples

### Basic info

Obtain essential information about your QuickSight datasets to understand their configuration and structure.

```sql+postgres
select
  name,
  dataset_id,
  arn,
  created_time,
  last_updated_time,
  import_mode
from
  aws_quicksight_dataset;
```

```sql+sqlite
select
  name,
  dataset_id,
  arn,
  created_time,
  last_updated_time,
  import_mode
from
  aws_quicksight_dataset;
```

### List datasets using SPICE import mode

Identify datasets that are using SPICE (Super-fast, Parallel, In-memory Calculation Engine) for data import.

```sql+postgres
select
  name,
  dataset_id,
  created_time,
  last_updated_time
from
  aws_quicksight_dataset
where
  import_mode = 'SPICE';
```

```sql+sqlite
select
  name,
  dataset_id,
  created_time,
  last_updated_time
from
  aws_quicksight_dataset
where
  import_mode = 'SPICE';
```

### Find datasets with row-level permissions enabled

Examine datasets that implement row-level security to ensure proper data access controls.

```sql+postgres
select
  name,
  dataset_id,
  row_level_permission_data_set
from
  aws_quicksight_dataset
where
  row_level_permission_data_set is not null;
```

```sql+sqlite
select
  name,
  dataset_id,
  row_level_permission_data_set
from
  aws_quicksight_dataset
where
  row_level_permission_data_set is not null;
```

### List datasets with column-level permissions

Identify datasets that have column-level permission rules defined to understand access control at a granular level.

```sql+postgres
select
  name,
  dataset_id,
  jsonb_array_length(column_level_permission_rules) as permission_rule_count
from
  aws_quicksight_dataset
where
  column_level_permission_rules is not null and
  jsonb_array_length(column_level_permission_rules) > 0;
```

```sql+sqlite
select
  name,
  dataset_id,
  json_array_length(column_level_permission_rules) as permission_rule_count
from
  aws_quicksight_dataset
where
  column_level_permission_rules is not null;
```

### Get dataset output columns information

Examine the structure of the output columns in a specific dataset to understand the available fields for analysis.

```sql+postgres
select
  name,
  dataset_id,
  o ->> 'Name' as column_name,
  o ->> 'Type' as data_type
from
  aws_quicksight_dataset,
  jsonb_array_elements(output_columns) as o
where
  dataset_id = 'example-dataset-id';
```

```sql+sqlite
select
  name,
  dataset_id,
  json_extract(o.value, '$.Name') as column_name,
  json_extract(o.value, '$.Type') as data_type
from
  aws_quicksight_dataset,
  json_each(output_columns) as o
where
  dataset_id = 'example-dataset-id';
```

### Analyze datasets by creation date

Identify recently created datasets to track new analytical resources being developed.

```sql+postgres
select
  name,
  dataset_id,
  created_time,
  import_mode
from
  aws_quicksight_dataset
order by
  created_time desc
limit 10;
```

```sql+sqlite
select
  name,
  dataset_id,
  created_time,
  import_mode
from
  aws_quicksight_dataset
order by
  created_time desc
limit 10;
```
