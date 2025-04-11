---
title: "Steampipe Table: aws_quicksight_data_set - Query AWS QuickSight Datasets using SQL"
description: "Allows users to query AWS QuickSight Datasets and retrieve information about datasets within AWS QuickSight."
folder: "QuickSight"
---

# Table: aws_quicksight_data_set - Query AWS QuickSight Datasets using SQL

AWS QuickSight is a fully managed, serverless business intelligence service that makes it easy to extend the power of machine learning generated insights to everyone in your organization. AWS QuickSight Dataset is a collection of data, either from a single source or from multiple sources. You create a dataset before you can analyze your data.

## Table Usage Guide

The `aws_quicksight_data_set` table in Steampipe provides you with information about datasets within AWS QuickSight. This table allows you, as a data analyst or administrator, to query dataset-specific details, including the data source, import mode, and associated configurations. You can utilize this table to gather insights on datasets, such as dataset structure, associated data sources, import configurations, and more. The schema outlines the various attributes of the QuickSight dataset for you to query, including the dataset ID, ARN, creation time, and associated tags.

## Examples

### Basic info
Explore which datasets are available in your AWS QuickSight environment and when they were created. This can help you manage and organize your data resources effectively.

```sql+postgres
select
  name,
  dataset_id,
  created_time,
  last_updated_time
from
  aws_quicksight_data_set;
```

```sql+sqlite
select
  name,
  dataset_id,
  created_time,
  last_updated_time
from
  aws_quicksight_data_set;
```

### List datasets with SPICE import mode
Identify datasets that use SPICE as their import mode. This is useful for understanding which datasets leverage QuickSight's in-memory calculation engine for faster analysis.

```sql+postgres
select
  name,
  dataset_id,
  import_mode
from
  aws_quicksight_data_set
where
  import_mode = 'SPICE';
```

```sql+sqlite
select
  name,
  dataset_id,
  import_mode
from
  aws_quicksight_data_set
where
  import_mode = 'SPICE';
```

### Find datasets with row-level permissions
Determine which datasets have row-level permission controls in place. This helps identify datasets with enhanced security controls that restrict user access to specific rows of data.

```sql+postgres
select
  name,
  dataset_id,
  row_level_permission_data_set
from
  aws_quicksight_data_set
where
  row_level_permission_data_set is not null;
```

```sql+sqlite
select
  name,
  dataset_id,
  row_level_permission_data_set
from
  aws_quicksight_data_set
where
  row_level_permission_data_set is not null;
```

### Get datasets with column-level permissions
Discover datasets that implement column-level permission rules to control access to specific columns of data. This is helpful for assessing security controls across your datasets.

```sql+postgres
select
  name,
  dataset_id,
  column_level_permission_rules
from
  aws_quicksight_data_set
where
  column_level_permission_rules is not null;
```

```sql+sqlite
select
  name,
  dataset_id,
  column_level_permission_rules
from
  aws_quicksight_data_set
where
  column_level_permission_rules is not null;
```

### List datasets created in the last 30 days
Identify recently created datasets to monitor new data resources or track recent changes to your QuickSight environment.

```sql+postgres
select
  name,
  dataset_id,
  created_time
from
  aws_quicksight_data_set
where
  created_time > now() - interval '30 days'
order by
  created_time desc;
```

```sql+sqlite
select
  name,
  dataset_id,
  created_time
from
  aws_quicksight_data_set
where
  created_time > datetime('now', '-30 days')
order by
  created_time desc;
```
