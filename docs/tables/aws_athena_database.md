---
title: "Steampipe Table: aws_athena_database - Query AWS Athena Databases using SQL"
description: "Allows users to query AWS Athena Databases to retrieve information about databases within data catalogs, including their name and description."
folder: "Athena"
---

# Table: aws_athena_database - Query AWS Athena Databases using SQL

The AWS Athena Database is a logical grouping of tables within a data catalog. It provides a way to organize and manage your data in a structured manner, making it easier to query and analyze using Amazon Athena. Each database is associated with a specific data catalog and contains metadata about the tables and their structure.

## Table Usage Guide

The `aws_athena_database` table in Steampipe provides you with information about databases within AWS Athena. This table allows you, as a DevOps engineer or data analyst, to query database-specific details, including the database name, catalog name, and description. You can utilize this table to gather insights on databases, such as their configuration and associated metadata. The schema outlines the various attributes of the Athena database for you, including the database name, catalog name, description, and associated tags.

## Examples

### Basic info
Explore the basic information about your AWS Athena databases, including their names, catalog names, and descriptions. This can help you understand the structure and organization of your databases.

```sql+postgres
select
  name,
  catalog_name,
  description
from
  aws_athena_database;
```

```sql+sqlite
select
  name,
  catalog_name,
  description
from
  aws_athena_database;
```

### List databases in a specific catalog
Identify databases that are associated with a particular catalog, such as 'AwsDataCatalog'. This can be useful for managing and organizing your databases based on their catalog.

```sql+postgres
select
  name,
  catalog_name,
  description
from
  aws_athena_database
where
  catalog_name = 'AwsDataCatalog';
```

```sql+sqlite
select
  name,
  catalog_name,
  description
from
  aws_athena_database
where
  catalog_name = 'AwsDataCatalog';
``` 