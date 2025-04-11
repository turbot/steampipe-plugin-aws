---
title: "Steampipe Table: aws_quicksight_namespace - Query AWS QuickSight Namespaces using SQL"
description: "Allows users to query AWS QuickSight Namespaces, providing details about namespace configurations, identity stores, and status information."
folder: "QuickSight"
---

# Table: aws_quicksight_namespace - Query AWS QuickSight Namespaces using SQL

AWS QuickSight Namespace is a logical container that helps organize and manage QuickSight resources. Each namespace can have its own identity store and configuration settings, allowing for better organization and access control of QuickSight resources.

## Table Usage Guide

The `aws_quicksight_namespace` table in Steampipe provides you with information about namespaces within AWS QuickSight. This table allows you, as an administrator, to query namespace-specific details, including identity store configurations and creation status. You can utilize this table to gather insights on namespace management, such as creation status, capacity regions, and associated identity stores.

## Examples

### Basic info
Explore the fundamental details of your QuickSight namespaces to understand their organization and structure.

```sql+postgres
select
  name,
  arn,
  capacity_region,
  creation_status,
  identity_store
from
  aws_quicksight_namespace;
```

```sql+sqlite
select
  name,
  arn,
  capacity_region,
  creation_status,
  identity_store
from
  aws_quicksight_namespace;
```

### List namespaces with creation errors
Identify namespaces that encountered errors during creation to troubleshoot issues.

```sql+postgres
select
  name,
  creation_status,
  namespace_error
from
  aws_quicksight_namespace
where
  namespace_error is not null;
```

```sql+sqlite
select
  name,
  creation_status,
  namespace_error
from
  aws_quicksight_namespace
where
  namespace_error is not null;
```

### Get namespaces using IAM Identity Center
Determine which namespaces are configured to use IAM Identity Center for authentication.

```sql+postgres
select
  name,
  arn,
  identity_store,
  creation_status
from
  aws_quicksight_namespace
where
  identity_store = 'IAM_IDENTITY_CENTER';
```

```sql+sqlite
select
  name,
  arn,
  identity_store,
  creation_status
from
  aws_quicksight_namespace
where
  identity_store = 'IAM_IDENTITY_CENTER';
```

### List namespaces by capacity region
Analyze the distribution of namespaces across different capacity regions.

```sql+postgres
select
  capacity_region,
  count(*) as namespace_count,
  array_agg(name) as namespaces
from
  aws_quicksight_namespace
group by
  capacity_region;
```

```sql+sqlite
select
  capacity_region,
  count(*) as namespace_count,
  group_concat(name) as namespaces
from
  aws_quicksight_namespace
group by
  capacity_region;
```

### Get successfully created namespaces
List all namespaces that have been successfully created and are ready to use.

```sql+postgres
select
  name,
  arn,
  identity_store,
  capacity_region
from
  aws_quicksight_namespace
where
  creation_status = 'CREATED';
```

```sql+sqlite
select
  name,
  arn,
  identity_store,
  capacity_region
from
  aws_quicksight_namespace
where
  creation_status = 'CREATED';
```
