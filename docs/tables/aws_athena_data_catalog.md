---
title: "Steampipe Table: aws_athena_data_catalog - Query AWS Athena Data Catalogs using SQL"
description: "Allows users to query AWS Athena Data Catalogs to retrieve information about data catalogs, including their type and associated parameters."
folder: "Athena"
---

# Table: aws_athena_data_catalog - Query AWS Athena Data Catalogs using SQL

The AWS Athena Data Catalog is a central repository that stores metadata about your data sources. It allows you to organize and manage your data in a structured way, making it easier to query and analyze using Amazon Athena. The data catalog can contain databases, tables, and other metadata that describe your data stored in various sources like Amazon S3.

## Table Usage Guide

The `aws_athena_data_catalog` table in Steampipe provides you with information about data catalogs within AWS Athena. This table allows you, as a DevOps engineer or data analyst, to query catalog-specific details, including the catalog name and type. You can utilize this table to gather insights on data catalogs, such as their configuration and type. The schema outlines the various attributes of the Athena data catalog for you, including the catalog name, type, and associated tags.

## Examples

### Basic info
Explore the basic information about your AWS Athena data catalogs, including their names and types. This can help you understand the structure and organization of your data catalogs.

```sql+postgres
select
  name,
  type
from
  aws_athena_data_catalog;
```

```sql+sqlite
select
  name,
  type
from
  aws_athena_data_catalog;
```

### List catalogs of a specific type
Identify data catalogs that are of a particular type, such as GLUE or LAMBDA. This can be useful for managing and organizing your data catalogs based on their type.

```sql+postgres
select
  name,
  type
from
  aws_athena_data_catalog
where
  type = 'GLUE';
```

```sql+sqlite
select
  name,
  type
from
  aws_athena_data_catalog
where
  type = 'GLUE';
``` 